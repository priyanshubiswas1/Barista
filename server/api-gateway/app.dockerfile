FROM golang:alpine AS build
RUN apk add --no-cache git gcc g++ make ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go/bin/app ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/bin/app .
EXPOSE 8080
CMD ["/usr/bin/app"]