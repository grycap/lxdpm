package main

import (
	"fmt"
	//"html"
	"log"
	"net/http"
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"github.com/lxc/lxd"

	"github.com/gorilla/mux"
)
//var dummyConfig Config = {}
//var host string = "192.168.1.135" 
var host = flag.String("host", "158.42.104.141", "The port of the application.")
var port = flag.String("port", ":8080", "The port of the application.")
var cert = flag.String("cert","/certs/client.crt","Path to cert to use for client auth.")
var key = flag.String("key","/certs/client.key","Path to key to use for client auth.")
var servercrt = flag.String("servercrt","/certs/server.crt","Path to the server cert we are accesing.")

func main() {
	flag.Parse() // parse the flags

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/status", statusCmdClient)
	router.HandleFunc("/trusted", trustCmdClient)
	router.HandleFunc("/info", infoCmdClient)

	log.Println("Starting LXD platform manager server on ", *port)
	if err := http.ListenAndServe(*port,router); err != nil {
		log.Fatal("ListenAndServe:",err)
	}
	
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is a work in progress for the new LXD platform manager. Stay tuned!")
}

func initClient() (lxd.ConnectInfo){
	var CasaRemote = lxd.RemoteConfig{
		Addr:   fmt.Sprintf("%s:8443",*host),
		//Static: false,
		//Public: false,
	}/*
	var Remotes = map[string]lxd.RemoteConfig{
		"casa": 	CasaRemote,
	}
	/*var ConfigCasa = lxd.Config{
		DefaultRemote: 	"casa",
		Remotes: 		CasaRemotes,
	}*/

	var certPem,errcert = ioutil.ReadFile(*cert)
	if errcert != nil {
		fmt.Println(errcert)
	}
	var keyPem,errkey = ioutil.ReadFile(*key)
	if errkey != nil {
		fmt.Println(errkey)
	}
	var servercertPem,errcertserv = ioutil.ReadFile(*servercrt)
	if errcertserv != nil {
		fmt.Println(errcertserv)
	}

	//fmt.Println(string([]byte(certPem)))
	//fmt.Println(string([]byte(keyPem)))

	var ConnectInfoCasa = lxd.ConnectInfo{
		Name:	"casa",
		RemoteConfig: CasaRemote,
		ClientPEMCert: string([]byte(certPem)),
		ClientPEMKey: string([]byte(keyPem)),
		ServerPEMCert: string([]byte(servercertPem)),

	}
	return ConnectInfoCasa
}

func infoCmdClient(w http.ResponseWriter, r *http.Request) {

	//Init the client
	var info = initClient()
	//var cli,err = lxd.NewClient(&ConfigCasa,"casa")
	var cli,err = lxd.NewClientFromInfo(info)
	if err != nil {
		fmt.Println(err)
	}
	//Get server status through client
	var serverStatus, serverr = cli.ServerStatus()
	if serverr != nil {
		fmt.Println(serverr)
	}

	data, d_err := yaml.Marshal(&serverStatus)
	if d_err != nil {
			fmt.Println(d_err)
	}

	fmt.Printf("\n%s",data)
}

func trustCmdClient(w http.ResponseWriter, r *http.Request) {

	
	var info = initClient()
	//var cli,err = lxd.NewClient(&ConfigCasa,"casa")
	var cli,err = lxd.NewClientFromInfo(info)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fmt.Sprintf("%+v",cli))
	var response = cli.AmTrusted()
	fmt.Println("after request")
	fmt.Println(response)
}

func statusCmdClient(w http.ResponseWriter, r *http.Request) {

	
	var info = initClient()
	//var cli,err = lxd.NewClient(&ConfigCasa,"casa")
	var cli,err = lxd.NewClientFromInfo(info)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("%+v",cli))

	var response,resp_err = cli.GetServerConfigString()
	if resp_err != nil {
		fmt.Println(resp_err)
	}
	fmt.Println(response)
}