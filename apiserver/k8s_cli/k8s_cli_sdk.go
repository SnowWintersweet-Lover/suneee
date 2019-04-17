package k8s_cli

import (
	"fmt"
	"github.com/zhaozf-zhiming/suneee/apiserver/common/types"
	"github.com/zhaozf-zhiming/suneee/apiserver/etc/apiconfig"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"strings"
)

var ConfigPath = "/k8s_cli/dmin.conf"
var iCount = 0 //用于记录分页索引,也作为返回total值

func InitClient() (*kubernetes.Clientset, error) {
	kubeconfig, err := ioutil.ReadFile(filepath.Join(apiconfig.GetServerDir(), ConfigPath))
	if err != nil {
		return nil, err
	}
	restConf, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func QueryK8sInfo(queryInfo types.QueryDeployment) (*types.QueryOut, error) {
	clientset, err := InitClient()
	if err != nil {
		return nil, err
	}
	iCount = 0
	return QueryNamespace(clientset, queryInfo)
}
func QueryNamespace(clientset *kubernetes.Clientset, queryInfo types.QueryDeployment) (*types.QueryOut, error) {
	queryOut := new(types.QueryOut)

	if queryInfo.Namespace == "" {
		k8sNamespacelist, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, valNamespace := range k8sNamespacelist.Items {
			queryInfo.Namespace = valNamespace.Name //对命名空间重新赋值
			rtVal, err := QueryName(clientset, queryInfo)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					continue
				}
				return nil, err
			}
			queryOut.List = append(queryOut.List, rtVal.DeployList...)
		}
	} else {
		rtVal, err := QueryName(clientset, queryInfo)
		if err != nil {
			return nil, err
		}
		queryOut.List = rtVal.DeployList
		//deloymentOut.Namespaces = append(deloymentOut.Namespaces, *rtVal)
	}
	queryOut.Total = iCount
	return queryOut, nil
}

func QueryName(clientset *kubernetes.Clientset, queryInfo types.QueryDeployment) (*types.NamespaceInfo, error) {
	var deloyinfo types.DeploymentInfo
	namespace := new(types.NamespaceInfo)
	if queryInfo.Name == "" {
		k8sDeploylist, err := clientset.AppsV1beta1().Deployments(queryInfo.Namespace).List(v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, val := range k8sDeploylist.Items {
			deloyinfo.Name = val.ObjectMeta.GetName()
			deloyinfo.Namespace = val.ObjectMeta.GetNamespace()
			valImage := val.Spec.Template.Spec.Containers[0]
			deloyinfo.ImageName = valImage.Image
			iAvailable := val.Status.AvailableReplicas
			iReady := val.Status.ReadyReplicas
			deloyinfo.InsCount = fmt.Sprintf("%d/%d", iAvailable, iReady)
			if iAvailable == iReady {
				deloyinfo.Status = "0"
			} else if iAvailable < iReady {
				deloyinfo.Status = "1"
			} else {
				deloyinfo.Status = "2"
			}
			if iCount >= queryInfo.Start {
				if queryInfo.Limit == 0 {
					namespace.DeployList = append(namespace.DeployList, deloyinfo)
				} else {
					if iCount-queryInfo.Start < queryInfo.Limit {
						namespace.DeployList = append(namespace.DeployList, deloyinfo)
					}
				}
			}
			iCount++
		}
	} else {
		k8sDeploy, err := clientset.AppsV1beta1().Deployments(queryInfo.Namespace).Get(queryInfo.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		deloyinfo.Name = k8sDeploy.ObjectMeta.GetName()
		deloyinfo.Namespace = k8sDeploy.ObjectMeta.GetNamespace()
		valImage := k8sDeploy.Spec.Template.Spec.Containers[0]
		deloyinfo.ImageName = valImage.Image
		iAvailable := k8sDeploy.Status.AvailableReplicas
		iReady := k8sDeploy.Status.ReadyReplicas
		deloyinfo.InsCount = fmt.Sprintf("%d/%d", iAvailable, iReady)
		if iAvailable == iReady {
			deloyinfo.Status = "0"
		} else if iAvailable < iReady {
			deloyinfo.Status = "1"
		} else {
			deloyinfo.Status = "2"
		}
		if iCount >= queryInfo.Start {
			if queryInfo.Limit == 0 {
				namespace.DeployList = append(namespace.DeployList, deloyinfo)
			} else {
				if iCount-queryInfo.Start < queryInfo.Limit {
					namespace.DeployList = append(namespace.DeployList, deloyinfo)
				}
			}
		}
		iCount++
	}
	return namespace, nil
}
