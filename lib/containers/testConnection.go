package containers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// map keeping track of timeouts
var containerTimeouts sync.Map

// max retry count (should probably be an env variable)
var maxTimeout = 30

// testConnectionToContainer() takes the Id of a container and checks if it's reachable on Port 3000
// if port is nil, will try to automagically find port
func TestConnectionToContainer(ctName string, port int) (bool, error) {
	timeoutCounter := getTimeoutCount(ctName)
	cturl, err := GetContainerUrl(ctName, port)
	if err != nil {
    log.Println("Error getting URL for Container", ctName)
		return false, err
	}

	log.Printf("Trying to connect to %s (%s)\n", ctName, cturl)
	resp, err := http.Get(cturl)
	if err != nil {
    log.Println("Error getting response from", ctName)
		return false, err
	}
	log.Println("Code", resp.Status, err)
	if resp.StatusCode == 200 {
		return true, nil
	}
	log.Printf("Trying to connect to http://%s/view/%s/", cturl, ctName)
	resp, err = http.Get(fmt.Sprintf("%s/view/%s/", cturl, ctName))
	log.Println("Code", resp.Status, err)
	if resp.StatusCode == 200 {
		containerTimeouts.Delete(ctName)
		return true, nil
	}
	if timeoutCounter > maxTimeout {
		containerTimeouts.Delete(ctName)
		return false, fmt.Errorf("Could not connect to %s on %s", ctName, cturl)
	}
	return false, nil

}

// getTimeoutCount() retrieves the timeout from the map
func getTimeoutCount(ctName string) int {
	counter, ok := containerTimeouts.Load(ctName)
	if !ok {
		containerTimeouts.Store(ctName, 0)
		return 0
	} else {
		return counter.(int)
	}

}
