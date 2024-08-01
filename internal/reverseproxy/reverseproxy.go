package reverseproxy

import (
	"dtsrv/internal/containers"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

var containerProxies sync.Map

func newProxy(rawUrl string) (*httputil.ReverseProxy, error) {
    url, err := url.Parse(rawUrl)
    if err != nil {
        return nil, err
    }
  	proxy := &httputil.ReverseProxy{
  		Rewrite: func(r *httputil.ProxyRequest) {
  			r.SetURL(url)
  			r.Out.Host = r.In.Host
  		},
  	}
    log.Println("created Proxy for ", rawUrl)

    return proxy, nil
}

func NewContainerProxy(ctName string, url string) (error) {
  _ , ok := containerProxies.Load(ctName)
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

func HandleReverseProxy(w http.ResponseWriter, r *http.Request) {
  ctName := r.PathValue("ctName")
  if (ctName == "" || !(strings.HasPrefix(ctName, "dtsrv"))) {
    log.Println("Container", ctName, "not found")
    http.Error(w, "Container not found", http.StatusNotFound)
    return
  }
  proxy, err := loadProxy(ctName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
  proxy.ServeHTTP(w, r)
}


func loadProxy(ctName string) (*httputil.ReverseProxy, error) {
var proxyObject any
 proxyObject, ok := containerProxies.Load(ctName)
  if ok == false {
    url, err := containers.GetContainerUrl(ctName)
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
  proxy, ok := proxyObject.(*httputil.ReverseProxy)
  if ok == false {
        return nil, fmt.Errorf("This is not a reverse proxy")
  }
  return proxy, nil
}
