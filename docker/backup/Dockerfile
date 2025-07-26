# DeshChain Backup Utility Container
FROM alpine:3.18

# Install required packages
RUN apk add --no-cache \
    bash \
    curl \
    postgresql-client \
    aws-cli \
    kubectl \
    tar \
    gzip \
    jq \
    ca-certificates \
    tzdata

# Create backup user
RUN addgroup -g 1000 backup && \
    adduser -u 1000 -G backup -s /bin/bash -D backup

# Create backup directories
RUN mkdir -p /backups && \
    chown -R backup:backup /backups

# Copy backup scripts
COPY backup-scripts/ /usr/local/bin/
RUN chmod +x /usr/local/bin/*.sh

# Switch to backup user
USER backup

# Set working directory
WORKDIR /backups

# Default command
CMD ["/usr/local/bin/backup.sh"]