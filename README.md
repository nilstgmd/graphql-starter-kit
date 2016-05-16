# GraphQL Starter Kit

Golang implementation of [Tutorial: How to build a GraphQL server](https://medium.com/apollo-stack/tutorial-building-a-graphql-server-cddaa023c035#.wy5h1htxs).

## Quick Start

Get the code:
```sh
go get github.com/nilstgmd/graphql-starter-kit
```

Start the Docker container running the GraphQL server, a container with Cassandra and a container with MongoDB:
```sh
cd  $GOPATH/src/github.com/nilstgmd/graphql-starter-kit/ && make all
```
Use the GraphiQL IDE in your browser `http://localhost:8080/` or the API `http://localhost:8080/graphql?query=...` to execute queries against the server.

## Example

Using GraphiQL:
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
  getFortuneCookie
}
```

Using cURL:
```sh
curl -XPOST http://localhost:8080/graphql \
-H 'Content-Type: application/graphql' \
-d 'query Root{ author(firstName:"Chuck",lastName:"Norris"){firstName,lastName,posts{title,views}},getFortuneCookie }'
```
