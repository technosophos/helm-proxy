package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Masterminds/httputil"
	"github.com/technosophos/helm-proxy/transcode"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	addr := ":44133" // 44134-1
	paddr := "localhost:44134"
	proxy := transcode.New(paddr)
	http.HandleFunc("/", bootstrap(proxy))
	log.Printf("starting server on %s to %s", addr, paddr)
	http.ListenAndServe(addr, nil)
}

func bootstrap(proxy *transcode.Proxy) http.HandlerFunc {
	api := routes(proxy)
	rslv := httputil.NewResolver(routeNames(api))

	// The main http.HandlerFunc delegates to the right route handler.
	hf := func(w http.ResponseWriter, r *http.Request) {
		cfg, err := cleanKubeConfig()
		if err != nil {
			log.Printf("cannot create config: %s", err)
			http.Error(w, "authentication required", http.StatusUnauthorized)
		}

		// This merely ensures that the proxy can auth.
		//log.Printf("Connecting to %q", cli.Api)
		if cli, err := kubernetes.NewForConfig(cfg); err != nil {
			log.Printf("cannot get new Kube client: %s", err)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		} else if _, err := cli.Namespaces().List(v1.ListOptions{}); err != nil {
			log.Printf("cannot get namespaces: %s", err)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		path, err := rslv.Resolve(r)
		if err != nil {
			http.NotFound(w, r)
		}
		for _, rr := range api {
			if rr.path == path {
				if err := rr.handler(w, r); err != nil {
					log.Printf("error on path %q: %s", path, err)
					http.Error(w, "proxy operation failed", 500)
				}
			}
		}
	}
	return hf
}

type routeHandler func(w http.ResponseWriter, r *http.Request) error
type route struct {
	path    string
	handler routeHandler
}

func routes(proxy *transcode.Proxy) []route {
	return []route{
		// Status
		{"GET /", index},
		// List
		{"GET /v1/releases", proxy.List},
		// Get
		{"GET /v1/releases/*", proxy.Get},
		// Install
		{"POST /v1/releases", proxy.Install},
		// Upgrade
		{"POST /v1/releases/*", proxy.Upgrade},
		// Delete
		{"DELETE /v1/releases/*", proxy.Uninstall},
		// History
		{"GET /v1/releases/*/history", proxy.History},
		// Rollback
		{"POST /v1/releaes/*/history/*", proxy.Rollback},
	}
}

func routeNames(r []route) []string {
	rn := make([]string, len(r))
	for i, rr := range r {
		rn[i] = rr.path
	}
	return rn
}

func index(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte(`{status: "ok", versions:["v1"]}`))
	return err
}

func kubeConfig() (*rest.Config, error) {
	// Try in-cluster config:
	c, err := rest.InClusterConfig()
	if err == nil {
		return c, nil
	}
	c, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return c, err
	}
	return c, nil
}

const ProxyUserAgent = "helm-proxy"

func cleanKubeConfig() (*rest.Config, error) {
	c, err := kubeConfig()
	if err != nil {
		return c, err
	}

	// Scrub the data. One of the issues with the Config struct is that it
	// blends connection info, content preferences, TLS info, authn info, and
	// various other bits in a way that doesn't allow clean separation of
	// concerns. So we basically have to manually destroy certain data.
	c.Username = ""
	c.Password = ""
	c.BearerToken = ""
	c.Impersonate = ""
	c.UserAgent = ProxyUserAgent

	return c, nil
}

func setCredentials(r *http.Request, c *rest.Config) {
	// Bearer token:
	t := r.Headers().Get("x-auth-token")
	c.BearerToken = t

	// HTTP Basic auth
}
