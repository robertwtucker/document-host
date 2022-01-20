FROM golang:1.17-buster AS build
WORKDIR /go/src/app
# copy module files first so that they don't need to be downloaded again if no change
COPY go.* .
RUN go mod download
RUN go mod vendor
RUN go mod verify
COPY . .
# RUN make build
RUN CGO_ENABLED=0 go build -o /go/bin/server ./cmd/server
RUN chmod +x ./cmd/server/entrypoint.sh

FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /go/bin/server ./
COPY --from=build /go/src/app/cmd/server/entrypoint.sh /
COPY --from=build /go/src/app/config/*.yaml ./config/
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
