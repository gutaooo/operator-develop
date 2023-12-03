package sharedinformerexercise

// import (
// 	"fmt"

// 	"k8s.io/client-go/informers"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/tools/cache"
// 	"k8s.io/client-go/tools/clientcmd"
// 	"k8s.io/client-go/util/workqueue"
// 	logf "sigs.k8s.io/controller-runtime/pkg/log"
// )

// var (
// 	log = logf.Log.WithName("SharedInformer-Exercise")
// )

// // Sharedinformer 创建
// func SharedinformerFun() {

// 	// 1. 创建 config
// 	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 2. 创建 client
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// workqueue的使用
// 	// 创建一个限速队列
// 	limitQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "LimitQueue")

// 	// 3. 获取 Informer
// 	// 通过工厂方法创建Informer
// 	factory := informers.NewSharedInformerFactory(clientset, 0)
// 	// 这里是直接创建client-go内建的Informer了，直接调用clientgo-informerfactory-pod的Infomer()方法，进而调用到factory的InformerFor()方法
// 	informer := factory.Core().V1().Pods().Informer()

// 	// 4. 添加事件处理方法， envet handler, informer监听到事件了
// 	informer.AddEventHandler(cache.ResourceEventHandlerDetailedFuncs{
// 		AddFunc: func(obj interface{}, isInInitialList bool) {
// 			fmt.Println("add obj...")
// 			// 事件产生就入队，不做事件处理
// 			key, err := cache.MetaNamespaceKeyFunc(obj)
// 			if err != nil {
// 				fmt.Println("get key failed in AddFunc")
// 			}
// 			limitQueue.AddRateLimited(key)
// 		},
// 		UpdateFunc: func(oldObj, newObj interface{}) {
// 			fmt.Println("update obj...")
// 		},
// 		DeleteFunc: func(oldObj interface{}) {
// 			fmt.Println("delete obj...")
// 		},
// 	})

// 	// 5. 启动 Informer
// 	stopCh := make(chan struct{})
// 	factory.Start(stopCh)
// 	factory.WaitForCacheSync(stopCh)
// 	<-stopCh
// }
