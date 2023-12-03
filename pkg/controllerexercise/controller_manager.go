package controllerexercise

import (
	"log"

	"k8s.io/client-go/informers"

	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// 管理控制器
func Manager() {
	// 获取kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inclusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("can not get kube config")
		}
		config = inclusterConfig
	}

	// 创建client
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("can not create client")
	}

	// 创建informer
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()
	// 注册事件处理方法
	controller := NewController(clientSet, serviceInformer, ingressInformer)
	// informer.start
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	// 这里的waitforCacheSync作用是什么
	factory.WaitForCacheSync(stopCh)

	// todo: 为什么在这里Run
	controller.Run(stopCh)

}
