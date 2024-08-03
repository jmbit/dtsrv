package session

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	_ "github.com/joho/godotenv/autoload"
)

var SessionStore sessions.FilesystemStore

func init() {
  SessionStore = sessionStore()

}

func sessionStore() sessions.FilesystemStore {
  var key []byte
  if keystring, ok := os.LookupEnv("SESSION_KEY"); ok == true {
    key = []byte(keystring)
  } else {
    key = securecookie.GenerateRandomKey(32)
  }
  return *sessions.NewFilesystemStore(os.Getenv("SESSION_PATH"), key)
}


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

//GetContainers returns a list of all containers owned by this session
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

//OwnsContainer checks if a session owns a container
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


func serializeToJSON(slice []string) (string, error) {
	jsonData, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func deserializeFromJSON(data string) ([]string, error) {
	var slice []string
	err := json.Unmarshal([]byte(data), &slice)
	return slice, err
}
