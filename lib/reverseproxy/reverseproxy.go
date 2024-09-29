package reverseproxy

import (
	"github.com/jmbit/dtsrv/lib/containers"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

// map of pointers to proxies
var containerProxies sync.Map

// map of timestamps for proxy/	container access
var lastAccessed sync.Map

// newProxy generates a new reverse proxy for a given URL and returns a pointer to it
func newProxy(rawUrl string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	proxy := &httputil.ReverseProxy{
		//This is technically default behaviour, but explicitly written out for clarity
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(url)
			r.Out.Host = r.In.Host
		},
	}
	log.Println("created Proxy for ", rawUrl)

	return proxy, nil
}

// NewContainerProxy is a wrapper around newProxy that fetches the required info
// about the container, creates a proxy and stores the pointer
func NewContainerProxy(ctName string, url string) error {
	_, ok := containerProxies.Load(ctName)
	if ok == true {
		log.Println("Proxy for Container", ctName, "already exists")
		return nil
	}
	proxy, err := newProxy(url)
	if err != nil {
		return err
	}
	containerProxies.Store(ctName, proxy)
	log.Println("created proxy for ", ctName)
	return nil
}

// DeleteContainerProxy() removes proxy from map. Should be called when destroying the container
func DeleteContainerProxy(ctName string) {
	containerProxies.Delete(ctName)
}

// HandleUnauthorized() returns an unauthorized HTTP error (Used to e.g. block file explorer)
func HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

// HandleReverseProxy() gets the proxy from the map and serves it
func HandleReverseProxy(w http.ResponseWriter, r *http.Request) {
	ctName := r.PathValue("ctName")
	//Check if the parameter is not empty or doesn't have the full prefix
	if ctName == "" || !(strings.HasPrefix(ctName, "dtsrv")) {
		log.Println("Container", ctName, "not found")
		http.Error(w, "Invalid container requested", http.StatusBadRequest)
		return
	}

	proxy, err := loadProxy(ctName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Hand over to proxy
	proxy.ServeHTTP(w, r)
}

// loadProxy() retrieves the pointer to the proxy from the map
// If no proxy exists, tries to create one with the port guessed by NewContainerProxy()
func loadProxy(ctName string) (*httputil.ReverseProxy, error) {
	// sync.Map doesn't support typed objects
	var proxyObject any
	proxyObject, ok := containerProxies.Load(ctName)
	if ok == false {
		url, err := containers.GetContainerUrl(ctName, 0)
		if err != nil {
			return nil, err
		}
		err = NewContainerProxy(ctName, url)
		if err != nil {
			return nil, err
		}
		proxyObject, err = loadProxy(ctName)
		if err != nil {
			return nil, err
		}

	}
	// type checking
	proxy, ok := proxyObject.(*httputil.ReverseProxy)
	if ok == false {
		return nil, fmt.Errorf("This is not a reverse proxy")
	}
	return proxy, nil
}
