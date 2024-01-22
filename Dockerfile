FROM alpine:latest

WORKDIR /app

# Install Go and Python
RUN apk add --no-cache go python3

# install bash
RUN apk update
RUN apk add --no-cache bash

# Copy the go code and build the executable
COPY . .
RUN go build -o execService

CMD ["./execService"]