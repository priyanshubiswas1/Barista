FROM golang:alpine AS build
RUN apk add --no-cache git gcc g++ make ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go/bin/app ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates python3 py3-pip
WORKDIR /app
COPY --from=build /go/bin/app /usr/bin/app
COPY services/photogrammetry /app/services/photogrammetry
COPY requirements.txt /app/requirements.txt
RUN pip3 install --no-cache-dir -r /app/requirements.txt
EXPOSE 8080
CMD ["/usr/bin/app"]