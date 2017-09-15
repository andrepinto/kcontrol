package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"os"
	"net"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/andrepinto/kcontrol/k8s"

)

type HttpResponse struct{
	Data 		interface{} 	`json:"data"`
	Ip 		[]string 	`json:"ip"`
	Environment 	string  	`json:"environment"`
}



func main() {

	client, err := k8s.NewKubeClient(os.Getenv("KUBECONFIG"))

	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	routerControler := Router{
		client: client,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", routerControler.Index)
	//router.HandleFunc("/info/{data}", routerControler.Info)
	router.HandleFunc("/k8s/nodes", routerControler.ListNode)
	router.HandleFunc("/k8s/config-maps/{namespace}", routerControler.ListConfigMaps)

	log.Println("server: localhost:3000")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Fatal(http.ListenAndServe(":3000", loggedRouter))
}

type Router struct {
	client *k8s.KubeClient
}

func (rt *Router) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Simple Server!\n")
}

func (rt *Router)  Info(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := vars["data"]

	ips := []string{}

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ips = append(ips, ipv4.String())
			//fmt.Println("IPv4: ", ipv4)
		}
	}



	response := &HttpResponse{
		Data: data,
		Ip: ips,
		Environment: os.Getenv("K8S_NAMESPACE"),
	}

	json.NewEncoder(w).Encode(response)
}

func (rt *Router)  ListNode(w http.ResponseWriter, r *http.Request){
	response, err := rt.client.ListNode()

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(response)
}

func (rt *Router)  ListConfigMaps(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	response, err := rt.client.ConfigMaps(namespace)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(response)
}