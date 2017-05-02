package main

import (
	"fmt"
	"io/ioutil"
	"flag"
	"lxdpm/apilxd"
	"github.com/lxc/lxd"
)

var cert = flag.String("cert","./certs/client.crt","Path to cert to use for client auth.")
var key = flag.String("key","./certs/client.key","Path to key to use for client auth.")
var servercrt = flag.String("servercrt","./certs/servercerts/server.crt","Path to the server cert we are accesing.")

func main() {
	saveConfig(&DefaultConfig,DefaultConfig.ConfigDir+"/config.yaml")
	var apilx = apilxd.LxdpmApi{}
	//Client initialization
	fmt.Println("Creating client")
	var cli,err = lxd.NewClientFromInfo(initClientSocket())
	if err != nil {
		fmt.Println(err)
	} else {
		apilx.Cli = cli
	}

	apilx.Init()

	
}

func initClientHttp() (lxd.ConnectInfo){
	var CasaRemote = lxd.RemoteConfig{
		Addr:   "https://localhost:8443",
		Protocol: "https://",
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

func initClientSocket() (lxd.ConnectInfo){
	var CasaRemote = lxd.RemoteConfig{
		Addr:   "unix:///var/lib/lxd/unix.socket",
		Protocol: "unix://",
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