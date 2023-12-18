# Dockerfile
FROM alpine:latest

WORKDIR /app

# Install Go and Python
RUN apk add --no-cache go python3

CMD ["sh"]