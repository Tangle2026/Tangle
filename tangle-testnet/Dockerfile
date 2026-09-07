FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags='-s -w' -o /out/tangled ./cmd/tangled
FROM alpine:3.22
WORKDIR /app
COPY --from=build /out/tangled /usr/local/bin/tangled
COPY frontend ./frontend
VOLUME ["/app/data"]
EXPOSE 8080
ENTRYPOINT ["tangled","--data","/app/data/tangle.json","--http",":8080"]
