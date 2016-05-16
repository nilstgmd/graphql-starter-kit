FROM alpine:3.3
MAINTAINER Nils Meder <nilstgmd@gmx.de>

ADD ./graphql-server /

RUN mkdir /static
COPY ./static /static

WORKDIR /

EXPOSE 8080

ENTRYPOINT /graphql-server
