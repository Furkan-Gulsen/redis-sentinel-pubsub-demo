#!/bin/sh
sed -i "s/{SENTINEL_DOWN_AFTER}/$SENTINEL_DOWN_AFTER/" /etc/redis/sentinel.conf
sed -i "s/{SENTINEL_FAILOVER}/$SENTINEL_FAILOVER/" /etc/redis/sentinel.conf
exec "$@"
