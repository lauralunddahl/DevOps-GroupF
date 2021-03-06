FROM golang:latest

#RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR ./src

COPY ./src .

COPY /src/go.mod .
COPY /src/go.sum .
RUN go mod download

RUN GOOS=linux go build -o ./out/main

EXPOSE 8080
EXPOSE 9090

CMD ["./out/main"]