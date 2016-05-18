package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/graphql-go/handler"
	"github.com/nilstgmd/graphql-starter-kit/cassandra"
	"github.com/nilstgmd/graphql-starter-kit/mongo"
	"github.com/nilstgmd/graphql-starter-kit/schema"

	_ "net/http/pprof"
)

func main() {
	// Creates a GraphQL-go HTTP handler with the defined schema
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
		log.Println("Starting GraphQL Server on http://localhost:8080/")
		http.ListenAndServe(":8080", nil)
	}()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
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
