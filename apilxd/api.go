package apilxd

import (
	"fmt"
	"net/http"
	"runtime"
	"flag"
	"os"
	"crypto/x509"
	"bytes"
	"database/sql"

	"github.com/lxc/lxd"
	"github.com/lxc/lxd/shared"
	"github.com/gorilla/mux"
	"log"
	"gopkg.in/inconshreveable/log15.v2"
)
var host = flag.String("host", "158.42.104.141", "The port of the application.")
var port = flag.String("port", ":8080", "The port of the application.")
var debug = true

type LxdpmApi struct {
	mux 		*mux.Router
	Cli 		*lxd.Client
	db 			*sql.DB
	clientCerts	[]x509.Certificate
	DefConfig 	*Config
}

type Command struct {
	name          string
	untrustedGet  bool
	untrustedPost bool
	get           func(lx *LxdpmApi, r *http.Request) Response
	put           func(lx *LxdpmApi, r *http.Request) Response
	post          func(lx *LxdpmApi, r *http.Request) Response
	delete        func(lx *LxdpmApi, r *http.Request) Response
	patch         func(lx *LxdpmApi, r *http.Request) Response
}

var api10 = []Command{
	containersCmd,
	certificatesCmd,
	api10Cmd,
	containerCmd,
}

func (lx *LxdpmApi) createCmd(version string, c Command) {
	var uri string
	if c.name == "" {
		uri = fmt.Sprintf("/%s", version)
	} else {
		uri = fmt.Sprintf("/%s/%s", version, c.name)
	}

	lx.mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var resp Response
		resp = NotImplemented

		switch r.Method {
		case "GET":
			if c.get != nil {
				resp = c.get(lx, r)
			}
		case "PUT":
			if c.put != nil {
				resp = c.put(lx, r)
			}
		case "POST":
			if c.post != nil {
				resp = c.post(lx, r)
			}
		case "DELETE":
			if c.delete != nil {
				resp = c.delete(lx, r)
			}
		case "PATCH":
			if c.patch != nil {
				resp = c.patch(lx, r)
			}
		default:
			resp = NotFound
		}

		if err := resp.Render(w); err != nil {
			err := InternalError(err).Render(w)
			if err != nil {
				shared.LogErrorf("Failed writing error for error, giving up")
			}
		}
		runtime.GC()
	})
}

func (lx *LxdpmApi) Init() {
	/* Initialize the database */
	err := initializeDbObject(lx, "./lxdpm.db")
	if err != nil {
		fmt.Println(err)
	}

	var certpath string = ""
	var keypath string = ""
	if _, err := os.Stat("../serverlxd.crt"); err == nil {
		certpath = "../serverlxd.crt"
		keypath = "../serverlxd.key"
	}
	//APIServer initialization
	lx.mux = mux.NewRouter()
	lx.mux.StrictSlash(false)

	lx.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		SyncResponse(true, []string{"/1.0"}).Render(w)
	})

	for _, c := range api10 {
		lx.createCmd("1.0", c)
	}

	lx.mux.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shared.LogInfo("Sending top level 404", log15.Ctx{"url": r.URL})
		w.Header().Set("Content-Type", "application/json")
		NotFound.Render(w)
	})

	log.Println("Starting LXD platform manager server on ", *port)
	if certpath != "" {
		if err := http.ListenAndServeTLS(*port,certpath,keypath,lx.mux); err != nil {
			log.Fatal("ListenAndServe:",err)
		}
	} else {
		if err := http.ListenAndServe(*port,lx.mux); err != nil {
			log.Fatal("ListenAndServe:",err)
		}
	}

}

func (lx *LxdpmApi) isTrustedClient(r *http.Request) bool {
	if r.RemoteAddr == "@" {
		// Unix socket
		return true
	}

	if r.TLS == nil {
		return false
	}

	for i := range r.TLS.PeerCertificates {
		if lx.CheckTrustState(*r.TLS.PeerCertificates[i]) {
			return true
		}
	}

	return false
}

func (lx *LxdpmApi) CheckTrustState(cert x509.Certificate) bool {
	for k, v := range lx.clientCerts {
		if bytes.Compare(cert.Raw, v.Raw) == 0 {
			//shared.LogDebug("Found cert", log.Ctx{"k": k})
			fmt.Printf("Found cert %s",k)
			return true
		}
	}

	return false
}
