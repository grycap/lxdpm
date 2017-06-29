package apilxd

import (
	"encoding/pem"
	//"fmt"
	"net/http"
	"os"
	//"reflect"
	"syscall"
	"io/ioutil"

	"gopkg.in/lxc/go-lxc.v2"

	//"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	//"github.com/lxc/lxd/shared/osarch"
	"github.com/lxc/lxd/shared/version"
)


func api10Get(lx *LxdpmApi, r *http.Request) Response {
	srv := api.ServerUntrusted{
		/* List of API extensions in the order they were added.
		 *
		 * The following kind of changes require an addition to api_extensions:
		 *  - New configuration key
		 *  - New valid values for a configuration key
		 *  - New REST API endpoint
		 *  - New argument inside an existing REST API call
		 *  - New HTTPs authentication mechanisms or protocols
		 */
		APIExtensions: []string{
			"storage_zfs_remove_snapshots",
			"container_host_shutdown_timeout",
			"container_syscall_filtering",
			"auth_pki",
			"container_last_used_at",
			"etag",
			"patch",
			"usb_devices",
			"https_allowed_credentials",
			"image_compression_algorithm",
			"directory_manipulation",
			"container_cpu_time",
			"storage_zfs_use_refquota",
			"storage_lvm_mount_options",
			"network",
			"profile_usedby",
			"container_push",
			"container_exec_recording",
			"certificate_update",
			"container_exec_signal_handling",
			"gpu_devices",
			"container_image_properties",
			"migration_progress",
			"id_map",
			"network_firewall_filtering",
			"network_routes",
			"storage",
			"file_delete",
			"file_append",
			"network_dhcp_expiry",
		},
		APIStatus:  "stable",
		APIVersion: version.APIVersion,
		Public:     false,
		Auth:       "untrusted",
	}

	// If untrusted, return now
	/*
	if !lx.isTrustedClient(r) {
		return SyncResponseETag(true, srv, nil)
	}*/

	srv.Auth = "trusted"

	/*
	 * Based on: https://groups.google.com/forum/#!topic/golang-nuts/Jel8Bb-YwX8
	 * there is really no better way to do this, which is
	 * unfortunate. Also, we ditch the more accepted CharsToString
	 * version in that thread, since it doesn't seem as portable,
	 * viz. github issue #206.
	 */
	uname := syscall.Utsname{}
	if err := syscall.Uname(&uname); err != nil {
		return InternalError(err)
	}

	kernel := ""
	for _, c := range uname.Sysname {
		if c == 0 {
			break
		}
		kernel += string(byte(c))
	}

	kernelVersion := ""
	for _, c := range uname.Release {
		if c == 0 {
			break
		}
		kernelVersion += string(byte(c))
	}

	kernelArchitecture := ""
	for _, c := range uname.Machine {
		if c == 0 {
			break
		}
		kernelArchitecture += string(byte(c))
	}

	addresses := []string{}
	/*addresses, err := lx.ListenAddresses()
	if err != nil {
		return InternalError(err)
	}*/
	//Modify this
	var certificate string
	content,_ := ioutil.ReadFile("../clientlxd.crt")
	certificate = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: content}))

	architectures := []string{}
	//architectures = append(architectures, "test")
	/*
	for _, architecture := range lx.architectures {
		architectureName, err := osarch.ArchitectureName(architecture)
		if err != nil {
			return InternalError(err)
		}
		architectures = append(architectures, architectureName)
	}
	*/
	env := api.ServerEnvironment{
		Addresses:          addresses,
		Architectures:      architectures,
		Certificate:        certificate,
		Driver:             "lxc",
		DriverVersion:      lxc.Version(),
		Kernel:             kernel,
		KernelArchitecture: kernelArchitecture,
		KernelVersion:      kernelVersion,
		Server:             "lxd",
		ServerPid:          os.Getpid(),
		ServerVersion:      version.Version}

	//drivers := readStoragePoolDriversCache()
	/*for _, driver := range drivers {
		// Initialize a core storage interface for the given driver.
		sCore, err := storageCoreInit(driver)
		if err != nil {
			continue
		}

		if env.Storage != "" {
			env.Storage = env.Storage + " | " + driver
		} else {
			env.Storage = driver
		}

		// Get the version of the storage drivers in use.
		sVersion := sCore.GetStorageTypeVersion()
		if env.StorageVersion != "" {
			env.StorageVersion = env.StorageVersion + " | " + sVersion
		} else {
			env.StorageVersion = sVersion
		}
	}*/

	fullSrv := api.Server{ServerUntrusted: srv}
	fullSrv.Environment = env
	//fullSrv.Config = daemonConfigRender()

	return SyncResponseETag(true, fullSrv, fullSrv.Config)
}

var api10Cmd = Command{name: "", untrustedGet: true, get: api10Get}