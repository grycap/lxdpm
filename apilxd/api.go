package apilxd

import (
	"fmt"
	"net/http"
	"runtime"
	"flag"

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
	mux 	*mux.Router
	Cli 	*lxd.Client
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
//	containerCmd,
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

	if err := http.ListenAndServe(*port,lx.mux); err != nil {
		log.Fatal("ListenAndServe:",err)
	}

}
