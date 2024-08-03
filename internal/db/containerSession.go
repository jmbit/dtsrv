package db

import "log"


func addContainer(ctName string, sessionID string) error {
  _, err := conn.Exec("INSERT INTO containersession (container, session) VALUES ?, ?", ctName, sessionID)
  if err != nil {
    log.Printf("Unable add Container %s to session %s: %v\n", ctName, sessionID, err)
  }
  return err
}

func delContainer(ctName string) error {
  _, err := conn.Exec("DELETE FROM containersession WHERE container LIKE ?", ctName)
  if err != nil {
    log.Printf("Unable add Container %s from sessions table: %v\n", ctName, err)
  }
  return err
}
