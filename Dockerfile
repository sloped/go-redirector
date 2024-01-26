FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/main

RUN touch redirects

FROM busybox
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/redirects .
EXPOSE 8080/tcp
ENTRYPOINT ["./main"]
