package k8s

import (

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/informers"
	"time"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
	"fmt"
	apicorev1 "k8s.io/client-go/pkg/api/v1"
	"log"
)

type ServiceData struct {
	Code string `yaml:"code"`
	Name string `yaml:"name"`
	Schema string `yaml:"schema"`
	Endpoint string `yaml:"endpoint"`
	Port int `yaml:"port"`
	BasePath string `yaml:"basePath"`
	Cluster  string `yaml:"cluster"`
}

type Service struct {
	Services ServicesData
}

type ServicesData []* ServiceData

type KubeClient struct {
	Client *kubernetes.Clientset
}

func NewKubeClient(kubeconfig string) (*KubeClient, error){

	var (
		config *rest.Config
		err    error
	)

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	client := kubernetes.NewForConfigOrDie(config)

	kc := &KubeClient{
		Client: client,
	}


	sharedInformers := informers.NewSharedInformerFactory(client, 10*time.Minute)

	kc.Watch(sharedInformers.Core().V1().ConfigMaps())

	sharedInformers.Start(nil)

	return  kc, nil

}

func(kc *KubeClient) Watch(secretInformer informercorev1.ConfigMapInformer){



	secretInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				//fmt.Println("add:",obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				data := newObj.(*apicorev1.ConfigMap)
				log.Println("upd:",data.Name, " ", data.Namespace)
				key, _ := cache.MetaNamespaceKeyFunc(newObj)
				log.Println(key)

			},
			DeleteFunc: func(obj interface{}) {
				fmt.Println("del:",obj)
			},
		},
	)
}

func(kc *KubeClient) ListNode() ([]*Node, error){
	list, err := kc.Client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	listNodes := []*Node{}

	for _, node := range list.Items {
		listNodes = append(listNodes,&Node{
			Name:node.Name,
			Namespace:node.Namespace,
			Tags:node.Labels,
		})
	}

	return listNodes, nil

}

func(kc *KubeClient) ConfigMaps(namespace string) ([]*ConfigMaps, error){
	list, err := kc.Client.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	listNodes := []*ConfigMaps{}

	for _, node := range list.Items {
		listNodes = append(listNodes,&ConfigMaps{
			Name:node.Name,
			Namespace:node.Namespace,
			Data:node.Data,
		})
	}

	return listNodes, nil

}