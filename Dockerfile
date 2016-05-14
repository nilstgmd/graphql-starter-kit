FROM alpine:3.3
MAINTAINER Nils Meder <nilstgmd@gmx.de>

ADD ./graphql-server /

EXPOSE 8080

ENTRYPOINT /graphql-server
