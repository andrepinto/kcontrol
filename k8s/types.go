package k8s

type Node struct {
	Name string
	Namespace string
	Tags map[string]string
}

type ConfigMaps struct {
	Name string
	Namespace string
	Data map[string]string
}