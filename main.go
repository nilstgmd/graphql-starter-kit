package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/graphql-go/handler"
	"github.com/nilstgmd/graphql-starter-kit/schema"
)

// Host is the connection string for MongoDB.
const Host string = "localhost"

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	log.Println("Seeding mock data to MongoDB")
	session, err := mgo.Dial(Host)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	collection := session.DB("graphql").C("post")
	err = collection.Insert(&schema.Post{Title: "Learn Golang + GraphQL + Relay", Views: r1.Intn(100)},
		&schema.Post{Title: "Tutorial: How to build a GraphQL server", Views: r1.Intn(100)},
		&schema.Post{Title: "The Go Programming Language", Views: r1.Intn(100)},
		&schema.Post{Title: "Microservices in Go", Views: r1.Intn(100)},
		&schema.Post{Title: "Programming in Go: Creating Applications for the 21st Century", Views: r1.Intn(100)})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Creates a GraphQL-go HTTP handler with the previously defined schema and we also set it to return pretty JSON output
	handler := handler.New(&handler.Config{
		Schema: &schema.Schema,
		Pretty: true,
	})

	// serve a GraphQL endpoint at `/graphql`
	http.Handle("/graphql", handler)

	// serve a graphiql IDE
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	go func() {
		log.Println("Starting GraphQL Server on localhost:8080")
		http.ListenAndServe(":8080", nil)
	}()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			log.Println("Received an interrupt, stopping GraphQL Server...")
			cleanup()
			cleanupDone <- true
		}
	}()

	<-cleanupDone
}

func cleanup() {
	log.Println("Cleaning up mock data...")
	session, err := mgo.Dial(Host)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	collection := session.DB("graphql").C("post")
	_, err = collection.RemoveAll(bson.M{})
	if err != nil {
		log.Fatal(err)
	}
}
