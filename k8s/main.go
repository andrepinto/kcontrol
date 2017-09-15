package k8s

import (

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/informers"
	"time"
	informercorev1 "k8s.io/client-go/informers/core/v1"
)

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

	kc.Watch(sharedInformers.Core().V1().Secrets())

	sharedInformers.Start(nil)

	return  kc, nil

}

func(kc *KubeClient) Watch(secretInformer informercorev1.SecretInformer){

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