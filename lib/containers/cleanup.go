package containers

import (
	"log"
	"time"
)

// startCleanup() starts the goroutines responsible for cleaning up unused containers
// maxAge int in seconds
// interval int in seconds
func StartCleanup(maxAge int64, interval int64) {
  go deleteOldContainers(maxAge, interval)
}

// deleteOldContainers() continually loops and deletes old containers
func deleteOldContainers(maxAge int64, interval int64) {
  for {
    cts, err := ListContainers()
    if  err != nil {
      log.Println("Error in cleanup task:", err)
    }
    for _, ct := range cts {
      created := ct.Created
      ageDate := time.Now().Unix() - maxAge
      if created > ageDate {
        if ct.State == "running" {
          err := StopContainer(ct.Names[0])
          if err != nil {
            log.Printf("Error stopping Container %s: %v\n", ct.Names[0], err)
          }
        }
        err := DeleteContainer(ct.Names[0]) 
        if err != nil {
          log.Printf("Error deleting Container %s: %v\n", ct.Names[0], err)
        }
      }
      
    }
    time.Sleep(time.Second*time.Duration(interval))
  }


}
