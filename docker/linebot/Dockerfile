FROM golang:1.22 AS build

ARG VERSION
ARG REVISION
ARG REPOSITORY
ARG BUILD_INFO="github.com/ww24/linebot/internal/buildinfo"

WORKDIR $GOPATH/src/github.com/ww24/linebot
COPY . .

ENV CGO_ENABLED=0
RUN go build -buildmode pie \
  -ldflags "-X ${BUILD_INFO}.version=${VERSION} -X ${BUILD_INFO}.revision=${REVISION} -X ${BUILD_INFO}.repository=${REPOSITORY} -s -w" \
  -trimpath \
  -o /usr/local/bin/linebot ./cmd/linebot

FROM gcr.io/distroless/base:nonroot

COPY --from=build /usr/local/bin/linebot /usr/local/bin/linebot
ENTRYPOINT [ "linebot" ]
