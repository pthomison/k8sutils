package k8sutils

import (
	"context"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetClientSet() (*kubernetes.Clientset, error) {
	cs, err := internalClientSet()
	if err == nil {
		return cs, nil
	} else {
		// fmt.Println("Unable To Load Internal ClientSet, attempting External ClientSet")
		return externalClientSet()
	}
}

func externalClientSet() (*kubernetes.Clientset, error) {
	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		home := homedir.HomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func internalClientSet() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GetPods(cs *kubernetes.Clientset, namespace string) (*corev1.PodList, error) {
	listOptions := metav1.ListOptions{}
	pods, err := cs.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	return pods, err
}

func GetDeployments(cs *kubernetes.Clientset, namespace string) (*appsv1.DeploymentList, error) {
	listOptions := metav1.ListOptions{}
	deployments, err := cs.AppsV1().Deployments(namespace).List(context.TODO(), listOptions)

	return deployments, err
}
