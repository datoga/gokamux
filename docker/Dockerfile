FROM golang:1.17-alpine AS builder

ARG APP_VERSION
ARG APP_BUILD

RUN apk --no-cache add make gcc g++ dev86

WORKDIR /src

# Cache libraries
COPY go.mod go.sum /src/
RUN go mod download

# Compile the code
COPY . ./
RUN make

# Final OS with the app
FROM alpine:3.14.2
RUN apk --no-cache add ca-certificates

WORKDIR /app

ENV PATH=/app:$PATH

COPY --from=builder /src/bin/gokamux .

VOLUME ["/etc/gokamux"]

CMD ["/app/gokamux"]