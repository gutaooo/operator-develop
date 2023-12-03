package clientgo_example

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func ClientSet_example() {
	// 1. 构建config
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/gutao/gutaodev/gocode/operator-develop/pkg/clientgo_example/minikube_kubeconfig")
	if err != nil {
		panic(err)
	}

	// 2. 创建 ClientSet 对象
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 3. 使用ClientSet，直接获取已经实现好的client
	pods, err := clientSet.
		CoreV1().
		Pods("kube-system").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, pod := range pods.Items {
		fmt.Printf("namespace: %s, name: %s\n", pod.Namespace, pod.Name)
	}
}
