package apilxd

import (
	"os"
	"fmt"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"github.com/lxc/lxd/shared"
)

type Config struct {
	//Mainhost holds the name of the mainhost we are accesing.
	MainHost string `yaml:"main-host"`
	//Hosts holds the access configuration for all hosts.
	Hosts map[string]HostConfig `yaml:"hosts"`
	//This holds where the configuration is going to be stored.
	ConfigDir string `yaml:"config-dir"`
}

type HostConfig struct {
	//This holds the name/alias of the host.
	Name string `yaml:"name"`
	IPAddress string `yaml:"address"`
	//This is the socket address the host is using.
	SocketAddress string `yaml:"socket-address"`

}

var Main = HostConfig{
	Name: "local",
	IPAddress: "localhost",
	SocketAddress: "unix:///var/lib/lxd/unix.socket"}

var Host2 = HostConfig{
	Name: "lxdpm02",
	IPAddress: "10.0.0.17",
	SocketAddress: "unix:///var/lib/lxd/unix.socket"}

var Host3 = HostConfig{
	Name: "lxdpm03",
	IPAddress: "10.0.0.18",
	SocketAddress: "unix:///var/lib/lxd/unix.socket"}

var Host4 = HostConfig{
	Name: "lxdpm04",
	IPAddress: "10.0.0.19",
	SocketAddress: "unix:///var/lib/lxd/unix.socket"}

var DefaultHosts = map[string]HostConfig{
	"local": 	Main,
	"lxdpm02":	Host2,
	"lxdpm03":	Host3}

var DefaultConfig = Config{
	MainHost:	"local",
	Hosts:		DefaultHosts,
	ConfigDir:	"./config"}

func saveConfig(c *Config,filename string) error {
	os.Remove(filename + ".new")
	os.Mkdir(filepath.Dir(filename), 0700)
	f, err := os.Create(filename + ".new")
	if err != nil {
		return fmt.Errorf("cannot create config file: %v", err)
	}
	defer f.Close()
	defer os.Remove(filename + ".new")

	data, err := yaml.Marshal(c)
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("cannot write configuration: %v", err)
	}
	f.Close()
	err = shared.FileMove(filename+".new", filename)
	if err != nil {
		return fmt.Errorf("cannot rename temporary config file: %v", err)
	}
	return nil
}
