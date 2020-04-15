FROM golang AS build
WORKDIR /
RUN pwd
COPY prime.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o prime prime.go

# We need a shell so use alpine.
FROM alpine
COPY --from=build /prime .
ENTRYPOINT ["/prime"]
