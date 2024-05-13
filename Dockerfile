FROM golang:1.22.0

# Set destination for COPY
WORKDIR /app

# Copy Go app
COPY go.mod src/*.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ../manager

# Run
ENTRYPOINT [ "/manager" ]