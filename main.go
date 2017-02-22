package main

import (
	"fmt"
	"log"
	"net/http"
	"flag"
	"io/ioutil"
	"encoding/json"
	"reflect"
	
	"lxdpm/lxdcli"

	"github.com/gorilla/mux"
)
//var dummyConfig Config = {}
//var host string = "192.168.1.135" 
var host = flag.String("host", "158.42.104.141", "The port of the application.")
var port = flag.String("port", ":8080", "The port of the application.")
var cert = flag.String("cert","./certs/client.crt","Path to cert to use for client auth.")
var key = flag.String("key","./certs/client.key","Path to key to use for client auth.")
var servercrt = flag.String("servercrt","./certs/server.crt","Path to the server cert we are accesing.")

func main() {
	flag.Parse() // parse the flags

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/info", infoCmdClient)
	router.HandleFunc("/delete/{cname}", deleteCmdClient)
	router.HandleFunc("/list", listCmdClient)
	router.HandleFunc("/start/{cname}", startCmdClient)
	router.HandleFunc("/stop/{cname}", stopCmdClient)
	router.HandleFunc("/launch",launchCmdClient).Methods("POST")


	log.Println("Starting LXD platform manager server on ", *port)
	if err := http.ListenAndServe(*port,router); err != nil {
		log.Fatal("ListenAndServe:",err)
	}
	
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is a work in progress for the new LXD platform manager. Stay tuned!")
}

func launchCmdClient(w http.ResponseWriter, r *http.Request) {
	body,err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var t lxdcli.LaunchPostJSONModel
	var args_buffer []string  
	
	json.Unmarshal(body,&t)
	

	reflected_t := reflect.ValueOf(&t).Elem()
	for i := 0;i < reflected_t.NumField(); i++ {
		sfield := reflected_t.Field(i)

		if sfield.String() != "" {
			args_buffer = append(args_buffer,sfield.String())
		}
	}
	fmt.Println(args_buffer)
	launch := lxdcli.LaunchCommand(args_buffer...)
	//fmt.Println("ARGS:",launch.Cmd.Args)
	errlaunch := launch.Do()
	if errlaunch != nil {
		fmt.Println(errlaunch)
	}
}

func stopCmdClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stop := lxdcli.StopCommand(vars["cname"])
	err := stop.Do()
	if err != nil {
		fmt.Println(err)
	}
}

func startCmdClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	start := lxdcli.StartCommand(vars["cname"])
	err := start.Do()
	if err != nil {
		fmt.Println(err)
	}
}


func listCmdClient(w http.ResponseWriter, r *http.Request) {
	list := lxdcli.ListCommand()
	err := list.Do()
	if err != nil {
		fmt.Println(err)
	}
}

func deleteCmdClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	delete := lxdcli.DeleteCommand(vars["cname"])
	err := delete.Do()
	if err != nil {
		fmt.Println(err)
	}
}

func infoCmdClient(w http.ResponseWriter, r *http.Request) {
	info := lxdcli.InfoCommand()
	err := info.Do()
	if err != nil {
		fmt.Println(err)
	}
}
