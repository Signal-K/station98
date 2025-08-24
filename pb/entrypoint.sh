#!/bin/sh

/pb/pocketbase serve --http=0.0.0.0:8080 --dir /pb/data &
PB_PID=$!

# Wait for PocketBase to be ready
until curl -s http://localhost:8080/api/collections/events/records > /dev/null; do
  sleep 1
done

# Wait for PocketBase process to exit
wait $PB_PID
