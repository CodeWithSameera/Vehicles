FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/github.com/CodeWithSameera/Vehicles/
WORKDIR /go/src/github.com/CodeWithSameera/Vehicles
RUN go mod download
COPY . /go/src/github.com/CodeWithSameera/Vehicles
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/Vehicles github.com/CodeWithSameera/Vehicles
FROM golang:1.14.6-alpine3.12
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/CodeWithSameera/Vehicles/build/Vehicles /usr/bin/Vehicles
COPY .env /usr/bin
EXPOSE 8080 8080
ENTRYPOINT ["/bin/bash"]