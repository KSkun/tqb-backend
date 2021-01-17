# build
FROM golang:1.15.5-alpine3.12 AS build
COPY ./src /src
WORKDIR /src
ENV "GOPROXY" "https://goproxy.io"
RUN go build -o /build/app

# image
FROM alpine:latest
COPY --from=build /build/app /bin/app
RUN mkdir /env
WORKDIR /env
ENTRYPOINT [ "/bin/app" ]