FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_COMMAND=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/main ./cmd/main.go
    
# FROM alpine/ca-certificates:latest AS runner

# # Set the working directory (e.g., /app for consistency or / for root)
# WORKDIR /app

# # Copy only the compiled binary from the builder stage
# COPY --from=builder /app/main .

# # Create a non-root user for security
# # Add a group and then a user, specifying a home directory and no password.
# RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser

# # Set the user to run the application
# USER appuser

# # Expose the port your application listens on
# EXPOSE 8080

# # Define a health check for Kubernetes/orchestrators
# # Adjust the path and port as per your application's health endpoint
# HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
#     CMD wget --quiet --tries=1 --timeout=5 http://localhost:8080/healthz || exit 1

# Command to run the application
CMD ["./main"]