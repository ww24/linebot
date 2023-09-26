FROM golang:1.21 AS build

ARG VERSION
ARG REVISION
ARG REPOSITORY

WORKDIR $GOPATH/src/github.com/ww24/linebot
COPY . .

ENV CGO_ENABLED=0
RUN go build -buildmode pie \
  -ldflags "-X buildinfo.version=${VERSION} -X buildinfo.revision=${REVISION} -X buildinfo.repository=${REPOSITORY} -s -w" \
  -trimpath \
  -o /usr/local/bin/linebot ./cmd/linebot

FROM gcr.io/distroless/base:nonroot

COPY --from=build /usr/local/bin/linebot /usr/local/bin/linebot
ENTRYPOINT [ "linebot" ]
