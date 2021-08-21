FROM golang:1.17 AS build

WORKDIR $GOPATH/src/github.com/ww24/linebot
COPY . .

ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/linebot ./cmd/linebot

FROM gcr.io/distroless/static

COPY --from=build /usr/local/bin/linebot /usr/local/bin/linebot
ENTRYPOINT [ "linebot" ]
