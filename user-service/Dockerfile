# Base off Golang Alpine image
# I prefer to version lock so updates don't cause failures
FROM golang:1.16.3-alpine3.12

# Create our Code directory and alias it to the working directory
RUN mkdir -p /go/src/user-service
WORKDIR /go/src/user-service

# Copy the src code into the container
ADD src /go/src/user-service

# Set the PORT env variable and expose it
ENV PORT 8080
EXPOSE 8080

# Build the service
RUN go mod download
RUN go build -o bin/user-service
RUN chmod -R +x /go

ENTRYPOINT ["/go/src/user-service/bin/user-service"]
