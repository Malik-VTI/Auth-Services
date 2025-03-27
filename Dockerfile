FROM alpine:latest

WORKDIR /app

# Copy binary dari hasil orchestrion build di Jenkins
COPY auth-service .

# Expose port kalau perlu
EXPOSE 7070

CMD ["./auth-service"]
