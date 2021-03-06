package validate

import (
	"context"

	admissionregistration "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/api/apps/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	v1Options "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ignoreNotFound(err error) error {
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func Deployment(ctx context.Context, c client.Client, name, namespace string) (bool, error) {
	dp := v1.Deployment{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	err := c.Get(ctx, key, &dp)
	if err != nil {
		return false, ignoreNotFound(err)
	}

	expected := *dp.Spec.Replicas
	actual := dp.Status.AvailableReplicas
	return actual == expected, nil
}

func Webhook(ctx context.Context, c client.Client, name string) (bool, error) {
	webhook := admissionregistration.MutatingWebhookConfiguration{}
	key := client.ObjectKey{Name: name}
	err := c.Get(ctx, key, &webhook)
	return err == nil, ignoreNotFound(err)
}

func CRD(config *rest.Config, crdName string) (bool, error) {
	apiextensionsclientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return false, err
	}
	_, err = apiextensionsclientset.ApiextensionsV1beta1().
		CustomResourceDefinitions().Get(crdName, v1Options.GetOptions{})
	return err == nil, ignoreNotFound(err)
}
