package mongo

import (
	"log"
	"math/rand"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const host string = "mongo"

// Init will seed mock data into the database.
func Init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	log.Println("Seeding mock data to MongoDB")
	session, collection := Get()
	defer session.Close()

	_, err := collection.RemoveAll(bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Example from https://blog.golang.org
	err = collection.Insert(bson.M{"title": "Go 1.6 is released", "views": r1.Intn(100), "author": bson.M{"firstName": "Andrew", "lastName": "Gerrand"}},
		bson.M{"title": "Six years of Go", "views": r1.Intn(100), "author": bson.M{"firstName": "Andrew", "lastName": "Gerrand"}},
		bson.M{"title": "Testable Examples in Go", "views": r1.Intn(100), "author": bson.M{"firstName": "Andrew", "lastName": "Gerrand"}},
		bson.M{"title": "Go, Open Source, Community", "views": r1.Intn(100), "author": bson.M{"firstName": "Russ", "lastName": "Cox"}},
		bson.M{"title": "Generating code", "views": r1.Intn(100), "author": bson.M{"firstName": "Rob", "lastName": "Pike"}},
		bson.M{"title": "Arrays, slices (and strings): The mechanics of 'append'", "views": r1.Intn(100), "author": bson.M{"firstName": "Rob", "lastName": "Pike"}},
		bson.M{"title": "Errors are values", "views": r1.Intn(100), "author": bson.M{"firstName": "Rob", "lastName": "Pike"}})
	if err != nil {
		log.Fatal(err)
	}
}

// Cleanup will remove all mock data from the database.
func Cleanup() {
	log.Println("Cleaning up MongoDB...")
	session, collection := Get()
	defer session.Close()

	_, err := collection.RemoveAll(bson.M{})
	if err != nil {
		log.Fatal(err)
	}
}

//Get returns the session and a reference to the post collection.
func Get() (*mgo.Session, *mgo.Collection) {
	session, err := mgo.Dial(host)
	if err != nil {
		log.Fatal(err)
	}

	collection := session.DB("graphql").C("post")
	return session, collection
}
