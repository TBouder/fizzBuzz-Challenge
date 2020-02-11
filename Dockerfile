FROM golang:1.13.3

WORKDIR /go/src/github.com/TBouder/fizzBuzz-Challenge

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /go/src/github.com/TBouder/fizzBuzz-Challenge

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o fizzBuzz

ENTRYPOINT ["./fizzBuzz"]
EXPOSE 8000