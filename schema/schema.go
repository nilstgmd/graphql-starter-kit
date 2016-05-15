package schema

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/nilstgmd/graphql-starter-kit/cassandra"
	"github.com/nilstgmd/graphql-starter-kit/mongo"
	"gopkg.in/mgo.v2/bson"
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
	Views int    `json:"views"`
}

// Author representation.
type Author struct {
	Name      string `json:"lastName"`
	FirstName string `json:"firstName"`
}

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"views": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var authorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"posts": &graphql.Field{
			Type: graphql.NewList(postType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var results []Post

				session, collection := mongo.Get()
				defer session.Close()
				err := collection.Find(bson.M{}).All(&results)
				if err != nil {
					log.Fatal(err)
				}

				return results, nil
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
				session := cassandra.Get()
				defer session.Close()

				var name string
				var firstName string

				iter := session.Query(
					`SELECT name, firstName FROM graphql.author WHERE name = ? AND firstName = ? ALLOW FILTERING;`,
					p.Args["lastName"].(string),
					p.Args["firstName"].(string)).Iter()
				for iter.Scan(&name, &firstName) {
					return Author{FirstName: firstName, Name: name}, nil
				}

				if err := iter.Close(); err != nil {
					log.Fatal(err)
				}

				return Author{}, nil
			},
		},
		"getFortuneCookie": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				response, err := http.Get("http://fortunecookieapi.com/v1/cookie")
				if err != nil {
					log.Fatal(err)
				}
				defer response.Body.Close()

				decoder := json.NewDecoder(response.Body)

				var cookies []Cookie
				err = decoder.Decode(&cookies)
				if err != nil {
					log.Fatal(err)
				}

				return cookies[0].Fortune.Message, nil
			},
		},
	},
})

// Schema is the GraphQL schema served by the server.
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})
