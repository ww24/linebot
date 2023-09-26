FROM golang:1.21 AS build

ARG VERSION
ARG REVISION
ARG REPOSITORY
ARG BUILD_INFO="github.com/ww24/linebot/internal/buildinfo"

WORKDIR $GOPATH/src/github.com/ww24/screenshot
COPY . .

ENV CGO_ENABLED=0
RUN go build -buildmode pie \
  -ldflags "-X ${BUILD_INFO}.version=${VERSION} -X ${BUILD_INFO}.revision=${REVISION} -X ${BUILD_INFO}.repository=${REPOSITORY} -s -w" \
  -trimpath \
  -o /usr/local/bin/screenshot ./cmd/screenshot

FROM chromedp/headless-shell:100.0.4896.56

RUN apt-get update \
 && apt-get install -y --no-install-recommends dumb-init fonts-ipafont ca-certificates \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*
ENTRYPOINT ["dumb-init", "--"]

ENV TZ=Asia/Tokyo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezones

COPY --from=build /usr/local/bin /usr/local/bin
CMD [ "screenshot" ]
