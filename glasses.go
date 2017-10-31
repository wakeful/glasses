package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	k8sHostname   string
	matchPattern  = "map[kubernetes.io/ingress.class:traefik]"
	hostFile      = flag.String("host-file", "/etc/hosts", "host file location")
)

const (
	sectionStart = "# generated using glasses start #"
	sectionEnd   = "# generated using glasses end #"
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
	flag.Parse()

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

	block := []byte(fmt.Sprintf("%s\n%s\n%s\n", sectionStart, hostEntries, sectionEnd))

	fileContent, err := ioutil.ReadFile(*hostFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	re := regexp.MustCompile(fmt.Sprintf("(?ms)%s(.*)%s", sectionStart, sectionEnd))
	if re.Match(fileContent) {
		fileContent = re.ReplaceAll(fileContent, block)
	} else {
		fileContent = append(fileContent, block...)
	}

	if err := ioutil.WriteFile(*hostFile, fileContent, 0644); err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(hostEntries)

}
