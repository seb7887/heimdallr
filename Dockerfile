# Two stage build to create one 
# 1. Start from the golang base image as the builder
FROM golang:alpine AS builder

# Set the current workdir inside the cointainer
WORKDIR /go/src/heimdallr

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached
RUN go mod download

# Mkdir src
RUN mkdir ./src

# Copy the source from the current directory to the Working Directory inside the container
COPY . ./src

# Build the Go app
RUN cd ./src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o heimdallr .

# Copy the built api-gateway to top level
RUN cp ./src/heimdallr .

# Remove source codes that no longer needed
RUN rm -rf go.* *.go src

# 2. Use scratch image
FROM scratch

# Set working directory
WORKDIR /root/

# Copy file from builder image
COPY --from=builder /go/src/heimdallr/heimdallr .

# Expose heimdallr ports
EXPOSE 7000
EXPOSE 7001

# Run
CMD ["./heimdallr"]