package main

import (
  "fmt"
  "labix.org/v2/mgo"
  //"labix.org/v2/mgo/bson"
)

func login(db *mgo.Database, user, pass string) bool {
  fmt.Println(db, user, pass)
  err := db.Login(user, pass)
  if err != nil {
    return false;
  } else {
    return true;
  }
}

func passwords(filename string) []string {
  return []string{"fuck", "off", "asshole", "123"};
}

func main() {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    panic(err)
  }

  username := "abc"
  passwords := passwords("pass.pass")

  defer session.Close()
  session.SetMode(mgo.Monotonic, true)

  db := session.DB("admin");
  
  for _, password := range passwords {
    if (login(db, username, password)) {
      fmt.Println("WE DID IT", username, password)
    } else {
      fmt.Println("WE FAILED :(")
    }
  }
}


