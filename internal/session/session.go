package session

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"

	_ "github.com/joho/godotenv/autoload"
)

var SessionStore sessions.FilesystemStore

func InitSessions() {
	SessionStore = sessionStore()

}

// sessionStore() sets up the session filesystem store
func sessionStore() sessions.FilesystemStore {
	var key []byte
	if viper.GetString("web.sessionkey") != "" {
		key = []byte(viper.GetString("web.sessionkey"))
	} else {
    log.Println("Generating new session key")
		key = securecookie.GenerateRandomKey(32)
	}
  if _, err := os.Stat(viper.GetString("web.sessionpath")); err != nil {
    err := os.Mkdir(viper.GetString("web.sessionpath"), 700)
    if err != nil {
      log.Println("Error creating directory for sessions:", err)
      panic(1)
    }
  }
  log.Println("Sessions stored in ", viper.GetString("web.sessionpath"))
	return *sessions.NewFilesystemStore(viper.GetString("web.sessionpath"), key)
}

// AppendContainer() adds a container to a session
func AppendContainer(s *sessions.Session, ctName string) error {
	val := s.Values["containers"]
	if sliceString, ok := val.(string); ok {
		slice, err := deserializeFromJSON(sliceString)
		if err != nil {
			log.Println("Error deserializing Container list", err)
			return err
		}
		slice = append(slice, ctName)
		sliceString, err = serializeToJSON(slice)
		if err != nil {
			log.Println("Error serializing Container list", err)
			return err
		}
		s.Values["containers"] = sliceString
	} else {
		var err error
		s.Values["containers"], err = serializeToJSON([]string{ctName})
		if err != nil {
			log.Println("Error serializing Container list", err)
			return err
		}

	}
	return nil
}

// RemoveContainer() deletes a container from the session
func RemoveContainer(s *sessions.Session, ctName string) error {
	val := s.Values["containers"]
	if sliceString, ok := val.(string); ok {
		slice, err := deserializeFromJSON(sliceString)
		var returnSlice []string
		if err != nil {
			log.Println("Error deserializing Container list", err)
			return err
		}
		for _, ct := range slice {
			if ct != ctName {
				returnSlice = append(returnSlice, ct)
			}
		}
		sliceString, err = serializeToJSON(returnSlice)
		if err != nil {
			log.Println("Error serializing Container list", err)
			return err
		}
		s.Values["containers"] = sliceString
	} else {
		s.Values["containers"] = ""

	}
	return nil
}

// GetContainers returns a list of all containers owned by this session
func GetContainers(s *sessions.Session) ([]string, error) {
	val := s.Values["containers"]
	if sliceString, ok := val.(string); ok {
		slice, err := deserializeFromJSON(sliceString)
		if err != nil {
			log.Println("Error deserializing Container list", err)
			return nil, err
		}
		return slice, nil
	}
	return nil, nil
}

// OwnsContainer checks if a session owns a container (Prevents users terminating other users containers)
func OwnsContainer(s *sessions.Session, ctName string) (bool, error) {
	containerList, err := GetContainers(s)
	if err != nil {
		return false, err
	}
	for _, ct := range containerList {
		if ct == ctName {
			return true, nil
		}
	}

	return false, nil
}


// IsAdmin() checks if the session is logged in as Admin
func IsAdmin(s *sessions.Session) (bool, error) {
	if isAdmin, ok := s.Values["admin"].(bool); ok && isAdmin {
		return true, nil
	}

	return false, nil
}

// serializeToJSON() turns a slice to a JSON string
func serializeToJSON(slice []string) (string, error) {
	jsonData, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// deserializeFromJSON() turns a JSON string to a slice
func deserializeFromJSON(data string) ([]string, error) {
	var slice []string
	err := json.Unmarshal([]byte(data), &slice)
	return slice, err
}
