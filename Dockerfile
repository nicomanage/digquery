FROM golang:1.24.1-bookworm

WORKDIR /app

# Install dig package
RUN apt-get update && apt-get install -y bind9-dnsutils && rm -rf /var/lib/apt/lists/*

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]