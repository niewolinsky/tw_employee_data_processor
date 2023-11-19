FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /tw-employee-data-processor ./cmd/api

EXPOSE 4001

CMD ["/tw-employee-data-processor"]
