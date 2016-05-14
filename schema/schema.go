package schema

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"gopkg.in/mgo.v2"
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
			Type: graphql.NewList(postType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var results []Post

				session, err := mgo.Dial("localhost")
				if err != nil {
					log.Fatal(err)
				}
				defer session.Close()

				err = session.DB("graphql").C("post").Find(bson.M{}).All(&results)
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
				firstName := p.Args["firstName"].(string)
				lastName := p.Args["lastName"].(string)
				return firstName + " " + lastName, nil
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
