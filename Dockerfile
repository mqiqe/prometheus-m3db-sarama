#FROM golang:1.11-alpine as build
#
#RUN mkdir /prometheus-m3db-sarama
#WORKDIR /prometheus-m3db-sarama
#ENV GO111MODULE=on
#ENV GOPROXY=https://goproxy.io
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
#COPY . .
#RUN make

#FROM alpine:3.8
#COPY --from=build /prometheus-m3db-sarama/saramam3db /
#LABEL maintainer mqiqe@163.com
#ENTRYPOINT ["/saramam3db"]

FROM alpine:3.8
MAINTAINER mqiqe@163.com
#COPY ./conf /conf/
COPY ./saramam3db /
ENTRYPOINT ["/saramam3db"]

