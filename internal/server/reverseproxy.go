package server

import (
	"dtsrv/internal/containers"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var containerProxies sync.Map

func newProxy(rawUrl string) (*httputil.ReverseProxy, error) {
    url, err := url.Parse(rawUrl)
    if err != nil {
        return nil, err
    }
    proxy := httputil.NewSingleHostReverseProxy(url)
    log.Println("created Proxy for ", rawUrl)

    return proxy, nil
}

func NewContainerProxy(ctName string) ( error) {

  _ , ok := containerProxies.Load(ctName)
  if ok == true {
    log.Println("Proxy for Container", ctName, "already exists")
    return nil
  }
  url, err := containers.GetContainerUrl(ctName)
  if err != nil {
    return err
  }
  proxy, err := newProxy(url)
  if err != nil {
    return err
  }
  containerProxies.Store(ctName, proxy)
    return nil

}

func handleReverseProxy(w http.ResponseWriter, r *http.Request) {
  ctName := r.PathValue("ctid")
  if ctName == "" {
    http.Error(w, "Container not found", http.StatusNotFound)
    return
  }
  proxyObject, ok := containerProxies.Load(ctName)
  if ok == false {
    http.Error(w, "ReverseProxy not found", http.StatusNotFound)
    return
  }
  proxy, ok := proxyObject.(*httputil.ReverseProxy)
  if ok == false {
    http.Error(w, "supplied ReverseProxy is not a Proxy", http.StatusNotFound)
    return
  }
  proxy.ServeHTTP(w, r)
}
