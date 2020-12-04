FROM golang:1.15 as build
WORKDIR /src
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 go build -o hilink-exporter

FROM scratch
COPY --from=build /src/hilink-exporter /bin/hilink-exporter

EXPOSE      9770
ENTRYPOINT  [ "/bin/hilink-exporter" ]
