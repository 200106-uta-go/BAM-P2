#build stage
FROM golang:alpine AS builder 
RUN mkdir /app 
ADD . /app
WORKDIR /app
RUN apk add --no-cache git
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -o loggerd *.go

#production image
FROM scratch
COPY --from=builder /app .
EXPOSE 9090
CMD ["./loggerd"]