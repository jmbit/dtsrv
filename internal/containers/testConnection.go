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
func TestConnectionToContainer(ctName string) (bool, error) {
	timeoutCounter := getTimeoutCount(ctName)
	cturl, err := GetContainerUrl(ctName)
	if err != nil {
		return false, err
	}

	log.Printf("Trying to connect to %s", cturl)
	resp, err := http.Get(cturl)
	log.Println("Code", resp.Status, err)
	if resp.StatusCode == 200 {
		return true, nil
	}
	log.Printf("Trying to connect to http://%s/view/%s/", cturl, ctName)
	resp, err = http.Get(fmt.Sprintf("http://%s/view/%s/", cturl, ctName))
	log.Println("Code", resp.Status, err)
	if resp.StatusCode == 200 {
		containerTimeouts.Delete(ctName)
		return true, nil
	}
	if timeoutCounter > maxTimeout {
		containerTimeouts.Delete(ctName)
		return false, fmt.Errorf("Could not connect to %s on ", ctName, cturl)
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
