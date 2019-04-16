package k8s_cli

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"k8s.io/api/apps/v1beta1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitClient() (clientset *kubernetes.Clientset, err error) {
	var restConf *rest.Config

	if restConf, err = GetRestConf(); err != nil {
		return nil, err
	}

	if clientset, err = kubernetes.NewForConfig(restConf); err != nil {
		return nil, err
	}
	return clientset, nil
}

func GetRestConf() (restConf *rest.Config, err error) {
	var kubeconfig []byte
	if kubeconfig, err = ioutil.ReadFile("./admin.conf"); err != nil {
		return nil, err
	}

	if restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		return nil, err
	}
	return restConf, nil
}

func QueryK8sInfo() {
	var (
		clientset *kubernetes.Clientset
		podsList  *core_v1.PodList
		k8sDeploy *v1beta1.Deployment
		err       error
	)
	if clientset, err = InitClient(); err != nil {
		return
	}

	if podsList, err = clientset.CoreV1().Pods("bms-pre").List(meta_v1.ListOptions{}); err != nil {
		return
	}
	fmt.Printf("%v", podsList)

	k8sDeploy, err = clientset.AppsV1beta1().Deployments("bms-pre").Get("biz-rest", v1.GetOptions{})
	fmt.Println(k8sDeploy)

	if b, er := json.Marshal(k8sDeploy); er != nil {
		log.Fatal(er)
	} else {
		fmt.Println(string(b))
		fmt.Println("-===============================================")
	}

	return
}
