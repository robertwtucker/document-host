FROM golang:1.20.5-bullseye AS build-env

ARG BUILD_VERSION=development
ARG BUILD_REVISION=unknown
ARG PROJECT="github.com/robertwtucker/document-host"

WORKDIR /go/src/app
# copy module files first so that they don't need to be downloaded again if no change
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags \
  "-X ${PROJECT}/internal/config.appVersion=${BUILD_VERSION} -X ${PROJECT}/internal/config.revision=${BUILD_REVISION}" \
  -o /go/bin/app ./cmd/docuhost

FROM gcr.io/distroless/static

COPY --from=build-env /go/bin/app /
USER nonroot:nonroot
EXPOSE 8080
CMD ["/app","serve"]
