package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	k8sHostname   string
	matchPattern  = "map[kubernetes.io/ingress.class:traefik]"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE")
}

func k8sHost(config *rest.Config) string {
	u, err := url.Parse(config.Host)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return u.Hostname()
}

func main() {

	fmt.Println("# reading k8s config...")
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir(), ".kube", "config"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	k8sHostname = k8sHost(config)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	ingress, err := client.ExtensionsV1beta1().Ingresses("").List(metaV1.ListOptions{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	var hostEntries string

	for _, elem := range ingress.Items {
		for _, annotation := range elem.Annotations {
			if annotation == matchPattern {
				for _, rule := range elem.Spec.Rules {
					hostEntries = hostEntries + fmt.Sprintf("%s %s # %s\n", k8sHostname, rule.Host, elem.Name)
				}
			}
		}
	}

	fmt.Println(hostEntries)

}
