# Use an official Go runtime as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /build

# Copy the Go modules manifest and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

# Build the Go application inside the container
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ENV=local AWS_ACCESS_KEY_ID=AKIAXYRLKZCY7LBEYN7Z AWS_SECRET_ACCESS_KEY=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
RUN go build -o endava-coding-exercise

# Set the entry point command for the container
CMD ["./rie-kaneko/credit-cards-summary"]
