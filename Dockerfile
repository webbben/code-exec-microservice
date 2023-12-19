# Dockerfile
FROM alpine:latest

WORKDIR /app

# Install Go and Python
RUN apk add --no-cache go python3

# install bash
RUN apk update
RUN apk add --no-cache bash

CMD ["sh"]