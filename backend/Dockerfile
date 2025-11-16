
# Build stage
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o pgvms .

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/pgvms /app/pgvms
ENV PORT=8080
EXPOSE 8080
CMD ["/app/pgvms"]
