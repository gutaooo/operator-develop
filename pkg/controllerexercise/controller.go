package controllerexercise

import (
	"context"
	"reflect"
	"time"

	"log"

	ingressv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apisv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	serviceinformer "k8s.io/client-go/informers/core/v1"
	ingressinformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	servicelister "k8s.io/client-go/listers/core/v1"
	ingresslister "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	client        kubernetes.Interface
	serviceLister servicelister.ServiceLister
	ingressLister ingresslister.IngressLister
	// EventHandler 和 Worker 之间处理速度不一致
	// todo: EventHandler只是将informer产生的时间记录下来？不做实际的处理？
	// todo: 在workqueue中，通过key可以去informer中拿到对应的事件
	queue workqueue.RateLimitingInterface
}

func (c *Controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c *Controller) updateService(oldObj, newObj interface{}) {
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}

func (c *Controller) deleteIngress(obj interface{}) {
	ingress := obj.(*ingressv1.Ingress)
	owerreference := apisv1.GetControllerOf(ingress)
	if owerreference == nil {
		log.Println("有ingress, 没有service")
	}
	if owerreference.Kind != "Service" {
		log.Println("ingress的refernece不是service")
		return
	}
	// 自己生产key
	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.Add(key)
}

// 开始work
func (c *Controller) Run(stopCh <-chan struct{}) {
	workNum := 5
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopCh)
	}

	<-stopCh
}

// work要从queue取key
func (c *Controller) worker() {
	for c.processNetxItem() {

	}
}

func (c *Controller) processNetxItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	// todo: queue处理后要done
	defer c.queue.Done(item)

	key := item.(string)

	err := c.syncService(key)
	if err != nil {
		c.handlerError(key, err)
		return false
	}
	return true
}

// 调谐逻辑
// todo: 无法通过key知道是删除、新增还是更新事件
func (c *Controller) syncService(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// 可能是service删除事件触发的
	service, err := c.serviceLister.Services(namespace).Get(name)
	if err != nil {
		log.Println("获取servic失败, service可能不存在")
		return nil
	}

	// 新增和更新
	_, haveAnno := service.Annotations["Controller/svc-igs"]
	// 对比ingress
	ingress, err := c.ingressLister.Ingresses(namespace).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return nil
	}

	// 根据service上annotation状态去创建或删除ingress
	if haveAnno && errors.IsNotFound(err) {
		_, err := c.client.NetworkingV1().Ingresses(namespace).Create(context.TODO(), c.constructIngress(namespace, name), apisv1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	if !haveAnno && ingress != nil {
		err := c.client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, apisv1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) constructIngress(namespace, name string) *ingressv1.Ingress {
	ingress := &ingressv1.Ingress{}
	ingress.Namespace = namespace
	ingress.Name = name
	pathType := ingressv1.PathTypePrefix
	ingress.Spec = ingressv1.IngressSpec{
		Rules: []ingressv1.IngressRule{
			{
				Host: "gutaooohost",
				IngressRuleValue: ingressv1.IngressRuleValue{
					HTTP: &ingressv1.HTTPIngressRuleValue{
						Paths: []ingressv1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: ingressv1.IngressBackend{
									Service: &ingressv1.IngressServiceBackend{
										Name: name,
										Port: ingressv1.ServiceBackendPort{
											Name:   name,
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return ingress
}

// 发生错误，重新入队
func (c *Controller) handlerError(key string, err error) {
	retryNum := 10
	if c.queue.NumRequeues(key) <= retryNum {
		c.queue.AddRateLimited(key)
	}
	log.Println("[controller handle error]", err)
	// todo: forget的作用
	c.queue.Forget(key)
}

func NewController(client kubernetes.Interface, serviceInformer serviceinformer.ServiceInformer, ingressInformer ingressinformer.IngressInformer) Controller {
	controller := Controller{
		client:        client,
		serviceLister: serviceInformer.Lister(),
		ingressLister: ingressInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ControllerExercise"),
	}
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addService,
		UpdateFunc: controller.updateService,
	})
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: controller.deleteIngress,
	})
	return controller
}
