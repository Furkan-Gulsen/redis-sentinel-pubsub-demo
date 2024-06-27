REDIS_HOST="127.0.0.1"
REDIS_PORT="6383"

for i in {1..1000}; do
  key=$(uuidgen)
  value="value_$i"
  redis-cli -h $REDIS_HOST -p $REDIS_PORT SET $key $value
  echo "Set $key to $value"

  sleep_time=$(( (RANDOM % 491) + 10 ))
  sleep_time_sec=$(echo "scale=3; $sleep_time / 1000" | bc)
  echo "Sleeping for $sleep_time_sec seconds..."
  sleep $sleep_time_sec
done