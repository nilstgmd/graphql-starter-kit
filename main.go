package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/graphql-go/handler"
	"github.com/nilstgmd/graphql-starter-kit/cassandra"
	"github.com/nilstgmd/graphql-starter-kit/mongo"
	"github.com/nilstgmd/graphql-starter-kit/schema"
)

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

func init() {
	mongo.Init()
	cassandra.Init()
}

func cleanup() {
	mongo.Cleanup()
	cassandra.Cleanup()
}
