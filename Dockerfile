FROM golang:1.24.3-alpine3.21

# Copy scripts.
COPY --chmod=755 /scripts/ /scripts/

# Add dependencies
RUN apk add --no-cache bash make

# Set default working directory.
WORKDIR "/go/src/github.com/Le-BlitzZz/streaming-etl-app"

# Copy source to image.
COPY . .

# Expose HTTP port.
EXPOSE 8080

# Set the default command.
CMD ["/scripts/cmd.sh", "tail", "-f", "/dev/null"]
