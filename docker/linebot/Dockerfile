FROM golang:1.21 AS build

ARG VERSION

WORKDIR $GOPATH/src/github.com/ww24/linebot
COPY . .

ENV CGO_ENABLED=0
RUN go build -buildmode pie \
  -ldflags "-X main.version=${VERSION} -s -w" -trimpath \
  -o /usr/local/bin/linebot ./cmd/linebot

FROM gcr.io/distroless/base:nonroot

COPY --from=build /usr/local/bin/linebot /usr/local/bin/linebot
ENTRYPOINT [ "linebot" ]
