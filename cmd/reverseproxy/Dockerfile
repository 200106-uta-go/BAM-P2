FROM golang:alpine AS builder
# Create the directory where the application will reside
WORKDIR /app
ADD . .
RUN apk add git
# Install any dependencies
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -o rproxy *.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/. /app/.
EXPOSE 8080
ENTRYPOINT [ "./rproxy" ]

