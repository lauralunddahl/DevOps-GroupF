FROM golang:latest

#RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR ./src

COPY ./src .

COPY /src/go.mod .
COPY /src/go.sum .
RUN go mod download

RUN GOOS=linux go build -o ./out/minitwit

EXPOSE 8080

CMD ["./out/minitwit"]