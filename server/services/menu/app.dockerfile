FROM golang:alpine AS build
RUN apk add --no-cache git gcc g++ make ca-certificates
WORKDIR /src
# Build context is repo root. Copy only what we need.
COPY contracts ./contracts
COPY server/services/menu ./server/services/menu
WORKDIR /src/server/services/menu
RUN go mod download
RUN go build -o /go/bin/app ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates python3 py3-pip
WORKDIR /app
COPY --from=build /go/bin/app /usr/bin/app
COPY --from=build /src/server/services/menu/services/photogrammetry /app/services/photogrammetry
COPY --from=build /src/server/services/menu/requirements.txt /app/requirements.txt
RUN pip3 install --no-cache-dir -r /app/requirements.txt
EXPOSE 8080
CMD ["/usr/bin/app"]