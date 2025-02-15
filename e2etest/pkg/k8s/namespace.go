// SPDX-License-Identifier:Apache-2.0

package k8s

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
)

func CreateNamespace(cs clientset.Interface, name string, labels map[string]string) error {
	_, err := cs.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}, metav1.CreateOptions{})
	return err
}

func DeleteNamespace(cs clientset.Interface, name string) error {
	backoff := wait.Backoff{
		Steps:    5,
		Duration: 10 * time.Second,
		Factor:   2,
	}
	err := cs.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = wait.ExponentialBackoff(backoff, func() (bool, error) {
		_, err := cs.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return true, nil
		}
		if err != nil {
			return false, err
		}
		return false, nil
	})
	return err
}

func ApplyLabelsToNamespace(cs clientset.Interface, name string, labels map[string]string) error {
	ns, err := cs.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if ns.Labels == nil {
		ns.Labels = make(map[string]string)
	}
	for k, v := range labels {
		ns.Labels[k] = v
	}
	_, err = cs.CoreV1().Namespaces().Update(context.Background(), ns, metav1.UpdateOptions{})
	return err
}
