FROM golang:1.17-bullseye AS build-env

WORKDIR /go/src/app
# copy module files first so that they don't need to be downloaded again if no change
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/server

FROM gcr.io/distroless/static

COPY --from=build-env /go/bin/app /
USER nonroot:nonroot
EXPOSE 8080
CMD ["/app"]
