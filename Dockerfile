FROM golang:1.19-alpine3.16 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -o app .

FROM alpine:3.16
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/app .

RUN addgroup -S myusr && adduser -S myusr -G myusr
USER myusr

ENTRYPOINT [ "./app" ]