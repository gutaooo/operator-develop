package clientsetexercise

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	log = logf.Log.WithName("ClientSet-Exercise")
)

func ClientSetFun() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	pod, err := clientset.CoreV1().Pods("gt-ns").Get(context.TODO(), "minio-deploy-dc59c9b9d-m6p8l", v1.GetOptions{})
	if err != nil {
		log.Error(err, "clientSet get pod failed")
	} else {
		fmt.Println(pod.Name, "-----", pod.APIVersion, "----", pod.Namespace)
	}
}
