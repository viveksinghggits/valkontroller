package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	klusterv1alpha1 "github.com/viveksinghggits/kluster/pkg/apis/viveksingh.dev/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ValidateKlusterVersion(spek klusterv1alpha1.KlusterSpec) (bool, error) {
	client := initClient(spek.TokenSecret)
	if client == nil {
		return false, errors.New("failed to create DO client")
	}
	options, _, err := client.Kubernetes.GetOptions(context.Background())
	if err != nil {
		return false, err
	}

	for _, version := range options.Versions {
		if spek.Version == version.Slug {
			return true, nil
		}
	}
	return false, errors.New("The version is not supported")
}

func initClient(tokenSecret string) *godo.Client {
	token, err := getToken(tokenSecret)
	if err != nil {
		fmt.Printf("Error %s getttign token", err.Error())
		return nil
	}

	client := godo.NewFromToken(token)
	return client
}

func getToken(sec string) (string, error) {
	// this must go inside `init`
	// konfig := flag.String("konfig", "/home/vivek/.kube/config", "location to your kubeconfig file")
	// config, err := clientcmd.BuildConfigFromFlags("", *konfig)
	// if err != nil {
	// handle error
	// fmt.Printf("erorr %s building config from flags\n", err.Error())
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("error %s, getting inclusterconfig", err.Error())
	}
	// }
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		// handle error
		fmt.Printf("error %s, creating clientset\n", err.Error())
	}

	namespace := strings.Split(sec, "/")[0]
	name := strings.Split(sec, "/")[1]
	s, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(s.Data["token"]), nil
}
