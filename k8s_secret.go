package k8sutils

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyCorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metaCorev1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type Secret map[string][]byte

func GetSecrets(cs *kubernetes.Clientset, namespace string) (*corev1.SecretList, error) {
	secrets, err := cs.CoreV1().Secrets(namespace).
		List(context.TODO(), metav1.ListOptions{})

	return secrets, err
}

func GetSecret(cs *kubernetes.Clientset, name string, namespace string) (*corev1.Secret, error) {
	secret, err := cs.CoreV1().Secrets(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})

	return secret, err
}

func ApplySecret(cs *kubernetes.Clientset, name string, namespace string, secret Secret) (*corev1.Secret, error) {
	kind := "Secret"
	apiVersion := "v1"

	s, err := cs.CoreV1().Secrets(namespace).
		Apply(
			context.TODO(),
			&applyCorev1.SecretApplyConfiguration{
				TypeMetaApplyConfiguration: metaCorev1.TypeMetaApplyConfiguration{
					Kind:       &kind,
					APIVersion: &apiVersion,
				},
				ObjectMetaApplyConfiguration: &metaCorev1.ObjectMetaApplyConfiguration{
					Name: &name,
				},
				Data: secret,
			},
			metav1.ApplyOptions{
				FieldManager: "golang-utils",
			},
		)

	return s, err
}

func UpdateSecret(cs *kubernetes.Clientset, name string, namespace string, secret Secret) (*corev1.Secret, error) {
	secretClient := cs.CoreV1().Secrets(namespace)

	var apiSecret *v1.Secret

	// Retrieve the latest version of Secret before attempting update
	// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var getErr, updateErr error
		apiSecret, getErr = secretClient.Get(context.TODO(), name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Secret: %v", getErr))
		}

		if apiSecret.Data == nil {
			// var emptyData
			apiSecret.Data = make(map[string][]byte)
		}

		for k, v := range secret {
			apiSecret.Data[k] = v
		}

		apiSecret, updateErr = secretClient.Update(context.TODO(), apiSecret, metav1.UpdateOptions{})
		return updateErr
	})

	return apiSecret, retryErr
}

func DeleteSecret(cs *kubernetes.Clientset, name string, namespace string) error {
	err := cs.CoreV1().Secrets(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})

	return err
}
