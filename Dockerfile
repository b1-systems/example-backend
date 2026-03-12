FROM golang:1.26 AS build-stage
WORKDIR /app
COPY example-backend.go go.mod go.sum ./
COPY ini /usr/local/go/src/example-backend/ini
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /example-backend

FROM scratch AS release-stage
COPY --from=build-stage /example-backend /example-backend
COPY example-backend.ini.sample /example-backend.ini
ENTRYPOINT ["/example-backend"]
ENV CLIENT_ID=example-backend
ENV PROVIDER_URL=
ENV LISTEN_ADDRESS=0.0.0.0:8080
