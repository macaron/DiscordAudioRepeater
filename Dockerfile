FROM golang:alpine3.14 as builder

WORKDIR /go/src/app/
COPY . .
RUN apk add --update-cache alpine-sdk

RUN go build -o main main.go

FROM alpine:3.14

WORKDIR /root/
COPY --from=builder /go/src/app/main .

RUN apk add --update-cache ffmpeg

CMD [ "./main" ]
