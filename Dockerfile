FROM ubuntu:22.04

WORKDIR /app

# Copy binary dari hasil orchestrion build di Jenkins
COPY services-auth .

# Expose port kalau perlu
EXPOSE 7070

CMD ["./services-auth"]
