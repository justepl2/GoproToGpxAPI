FROM golang:latest
WORKDIR /app

RUN apt-get update && apt-get install -y libimage-exiftool-perl

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download
COPY ./src ./
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]