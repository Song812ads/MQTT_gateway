# syntax=docker/dockerfile:1

FROM golang:1.23.3

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . ./

# COPY ../docker-compose.yml ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 12345

# Run
CMD ["/gateway"]