package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Cookie containing fortunes.
type Cookie struct {
	Fortune Fortune `json:"fortune"`
}

// Fortune contains the message.
type Fortune struct {
	Message string `json:"message"`
}

// Post representation.
type Post struct {
	Title string `json:"title"`
	Views string `json:"views"`
}

// PostList represents multiple Posts.
var PostList []Post

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"views": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var authorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"firstName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "Peter", nil
			},
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "Lustig", nil
			},
		},
		"posts": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "List of todos",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return PostList, nil
			},
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"author": &graphql.Field{
			Type: authorType,
			Args: graphql.FieldConfigArgument{
				"firstName": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lastName": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				firstName := p.Args["firstName"].(string)
				lastName := p.Args["lastName"].(string)
				return firstName + " " + lastName, nil
			},
		},
	},
})

// Schema is the GraphQL schema served by the server.
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

func init() {
	// populate "database" with posts.
	post1 := Post{Title: "a", Views: "b"}
	post2 := Post{Title: "b", Views: "c"}
	PostList = append(PostList, post1, post2)
}

func main() {
	// Creates a GraphQL-go HTTP handler with the previously defined schema and we also set it to return pretty JSON output
	handler := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	// serve a GraphQL endpoint at `/graphql`
	http.Handle("/graphql", handler)

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
			cleanupDone <- true
		}
	}()

	<-cleanupDone
}
