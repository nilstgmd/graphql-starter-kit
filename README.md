# graphql-starter-kit

Golang implementation of https://github.com/apollostack/apollo-starter-kit.

## Schema

```go
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
```

## Example Query

```javascript
{
  author(firstName:"Chuck", lastName: "Norris"){
    firstName
    lastName
    posts{
      title
      views
    }
  }
}
```

```
http://localhost:8080/graphql?query={author(firstName:%22Chuck%22,%20lastName:%22Norris%22){firstName,lastName,posts{title,%20views}}}
```
