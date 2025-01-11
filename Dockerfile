FROM golang:1.23 AS build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN apt-get update && apt-get install -y gcc libc-dev unzip
RUN curl -fsSL https://deno.land/install.sh | sh
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN cd frontend && /root/.deno/bin/deno install && /root/.deno/bin/deno task build && cd ..
RUN go build main.go
RUN chmod +x /build/thecommongame

FROM ubuntu:latest AS dev
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=build /build/thecommongame /usr/local/bin/nowplaying
COPY --from=build /build/.env /usr/local/bin/.env
RUN chmod +x /usr/local/bin/thecommongame
EXPOSE 8080
ENTRYPOINT [ "/usr/local/bin/thecommongame" ]

FROM alpine:latest AS production
WORKDIR /app
COPY --from=build /build/thecommongame /app
EXPOSE 8080
ENTRYPOINT [ "/app/thecommongame" ]

