FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/merch_store
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:latest
EXPOSE 8080
WORKDIR /root
COPY --from=build /app/cmd/merch_store/main .
CMD ["./main"]
