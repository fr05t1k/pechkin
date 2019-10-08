FROM golang:1.12-alpine as builder

RUN apk --no-cache add ca-certificates git

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

ENV PECHKIN_BOT_TOKEN=""

CMD ["./main"]


