FROM golang AS build
WORKDIR /
COPY countme.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o countme countme.go

# We need a shell so use alpine.
FROM alpine
COPY --from=build /countme .
ENV COUNTME_PORT 80
ENV GOGC off
EXPOSE 80
ENTRYPOINT ["/countme"]
