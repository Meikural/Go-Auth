FROM alpine:latest

WORKDIR /app

# Copy the pre-built binary
COPY auth-service .

# Copy .env file (optional, can also be passed via environment variables)
COPY .env .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./auth-service"]