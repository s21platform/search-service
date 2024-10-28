FROM golang:1.22 as builder

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o build/kafka cmd/workers/notification/main.go
RUN go build -o build/main cmd/service/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /usr/src/service/build/main /app
RUN apk add --no-cache gcompat
# RUN chmod +x main

# RUN ls -l /app/
COPY --from=builder /usr/src/service ./scripts
COPY --from=builder /usr/src/service/build/kafka .

#CMD ["/app/main","/app/kafka"]


CMD ./main & ./kafka