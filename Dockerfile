FROM golang:1.23.2-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/radionica cmd/radionica/main.go
EXPOSE 8080
CMD ["/app/radionica"]