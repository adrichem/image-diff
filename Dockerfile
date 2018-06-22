# build stage
FROM golang:alpine AS build-env
ADD . /go/src/github.com/adrichem/image-diff
RUN cd /go/src/github.com/adrichem/image-diff && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /src/image-diff

# final stage
FROM scratch
WORKDIR /app
COPY --from=build-env /src/image-diff /app/
EXPOSE 80
ENTRYPOINT ["/app/image-diff"]