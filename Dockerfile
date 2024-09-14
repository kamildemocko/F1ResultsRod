FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o F1ResultsRod ./cmd/app

RUN chmod +x F1ResultsRod

FROM alpine:latest

WORKDIR /app

COPY .env .

COPY --from=builder /app/F1ResultsRod .

CMD [ "/app/F1ResultsRod" ]
