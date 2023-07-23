package restclientexercise11

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func RestclientExercise() {
	// 使用client-go
	// 1. config
	// clientcmd.RecommendedHomeFile：实际上是kube config的路径
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	// 2. client
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	config.APIPath = "/api"
	resetClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	// 3. get data
	pod := v1.Pod{}
	err = resetClient.Get().Resource("pods").Namespace("gt-ns").Name("minio-deploy-dc59c9b9d-m6p8l").Do(context.TODO()).Into(&pod)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pod.Name)
		fmt.Println(pod)
	}

}
