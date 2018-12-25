FROM circleci/golang:1.11

WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download

COPY . /app
