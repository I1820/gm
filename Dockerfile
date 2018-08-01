# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github.com/aiotrc/gm
RUN apk update && apk add git
RUN go get -u github.com/golang/dep/cmd/dep
RUN cd $GOPATH/src/github.com/aiotrc/gm/ && dep ensure && go build -v -o /gm

# Final stage
FROM alpine
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build-env /gm /app/
ENTRYPOINT ["./gm"]
