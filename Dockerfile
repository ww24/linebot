# syntax=docker/dockerfile:1.3-labs
FROM golang:1.17 AS build

ARG VERSION

WORKDIR $GOPATH/src/github.com/ww24/linebot
COPY . .

ENV CGO_ENABLED=0
RUN <<EOF
go build \
  -ldflags "-X main.version=${VERSION} -w" \
  -buildmode pie \
  -o /usr/local/bin/linebot ./cmd/linebot
EOF

FROM gcr.io/distroless/base:nonroot

COPY --from=build /usr/local/bin/linebot /usr/local/bin/linebot
ENTRYPOINT [ "linebot" ]
