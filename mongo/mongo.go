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

	err := collection.Insert(bson.M{"title": "Learn Golang + GraphQL + Relay", "views": r1.Intn(100)},
		bson.M{"title": "Tutorial: How to build a GraphQL server", "views": r1.Intn(100)},
		bson.M{"title": "The Go Programming Language", "views": r1.Intn(100)},
		bson.M{"title": "Microservices in Go", "views": r1.Intn(100)},
		bson.M{"title": "Programming in Go: Creating Applications for the 21st Century", "views": r1.Intn(100)})
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
