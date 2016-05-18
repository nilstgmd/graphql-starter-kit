# GraphQL Starter Kit for GoLang

Golang implementation of a GraphQL server using Cassandra, MongoDB and REST backends inspired by [Tutorial: How to build a GraphQL server](https://medium.com/apollo-stack/tutorial-building-a-graphql-server-cddaa023c035#.wy5h1htxs).

## Quick Start

Get the code:
```sh
go get github.com/nilstgmd/graphql-starter-kit
```

Start the Docker container running the GraphQL server, a container with Cassandra and a container with MongoDB ([showterm.io](http://showterm.io/df13d2ece08deb5cc5564)):
```sh
cd  $GOPATH/src/github.com/nilstgmd/graphql-starter-kit/ && make all
```
Use the GraphiQL IDE in your browser `http://localhost:8080/` or the API `http://localhost:8080/graphql?query=...` to execute queries against the server.

<p align="center">
  <img width=600 src="https://raw.githubusercontent.com/nilstgmd/graphql-starter-kit/master/images/graphiql.png" alt="GraphiQL IDE"/>
</p>

## Example

Using GraphiQL:
```javascript
{
  author(firstName:"Rob", lastName: "Pike"){
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
-d 'query Root{ author(firstName:"Andrew",lastName:"Gerrand"){firstName,lastName,posts{title,views}},getFortuneCookie }'
```

## Extra

Just because it is cool: Flamegraphs :fire:

<p align="center">
  <img width=600 src="https://raw.githubusercontent.com/nilstgmd/graphql-starter-kit/master/images/flame.png" alt="Flamegraph"/>
</p>

The GraphQL server offers a [pprof](https://golang.org/pkg/net/http/pprof/) endpoint which can be consumed to create flamegraphs with [go-torch](https://github.com/uber/go-torch).
