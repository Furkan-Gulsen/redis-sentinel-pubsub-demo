FROM redis:3

# Sentinel configuration
ENV SENTINEL_QUORUM 1
ENV SENTINEL_DOWN_AFTER 5000
ENV SENTINEL_FAILOVER 10000
ENV SENTINEL_PORT 26000

# Copy custom scripts and configuration
COPY sentinel-entrypoint.sh /usr/local/bin/
COPY sentinel.conf /etc/redis/sentinel.conf

# Set execute permission for entrypoint script
RUN chmod +x /usr/local/bin/sentinel-entrypoint.sh

# Set the entrypoint and default command
ENTRYPOINT ["sentinel-entrypoint.sh"]
CMD ["redis-server", "/etc/redis/sentinel.conf", "--sentinel"]
