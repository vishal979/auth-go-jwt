FROM golang:alpine as builder

LABEL maintainer="Vishal Sharma <vishal.sharma09890@gmail.com>"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
RUN mkdir logs

EXPOSE 7070

CMD ["./main"]

