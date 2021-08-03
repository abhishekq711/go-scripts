FROM golang:1.16-alpine
WORKDIR /go/src/app

COPY . .
RUN go mod download

RUN go build -o /go-code-server
EXPOSE 8080
CMD [ "/go-code-server" ]
