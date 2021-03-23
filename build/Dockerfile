FROM golang:alpine AS build
RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
            ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.* ./

# copy source files and build the binary
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /app/
COPY --from=build /app/bin/crawler .
RUN ls -la
ENTRYPOINT ["/app/crawler"]