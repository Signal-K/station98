FROM golang:1.24.5-alpine

# Install curl for API checking
RUN apk add --no-cache curl

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o backend main.go
RUN go build -o utils cmd/utils-main.go

# Create startup script
RUN echo '#!/bin/sh' > /app/startup.sh && \
    echo 'echo "🚀 Starting backend service with initial sync and cleanup..."' >> /app/startup.sh && \
    echo '# Start backend in background to perform initial sync' >> /app/startup.sh && \
    echo './backend &' >> /app/startup.sh && \
    echo 'BACKEND_PID=$!' >> /app/startup.sh && \
    echo '' >> /app/startup.sh && \
    echo '# Wait for initial sync to complete (check for events in PocketBase)' >> /app/startup.sh && \
    echo 'echo "⏳ Waiting for initial sync to complete..."' >> /app/startup.sh && \
    echo 'COUNTER=0' >> /app/startup.sh && \
    echo 'while [ $COUNTER -lt 10 ]; do' >> /app/startup.sh && \
    echo '  # Check if we have events in the database' >> /app/startup.sh && \
    echo '  EVENT_COUNT=$(curl -s "${PB_URL}/api/collections/events/records?page=1&perPage=1" | grep -o "totalItems":[0-9]* | cut -d: -f2 | head -n1)' >> /app/startup.sh && \
    echo '  if [ "$EVENT_COUNT" != "" ] && [ "$EVENT_COUNT" -gt 0 ]; then' >> /app/startup.sh && \
    echo '    echo "✅ Found $EVENT_COUNT events, initial sync completed!"' >> /app/startup.sh && \
    echo '    break' >> /app/startup.sh && \
    echo '  fi' >> /app/startup.sh && \
    echo '  echo "🔄 Waiting for events to be synced... ($COUNTER/10)"' >> /app/startup.sh && \
    echo '  sleep 2' >> /app/startup.sh && \
    echo '  COUNTER=$((COUNTER + 1))' >> /app/startup.sh && \
    echo 'done' >> /app/startup.sh && \
    echo '' >> /app/startup.sh && \
    echo '# Run cleanup after sync is complete' >> /app/startup.sh && \
    echo 'echo "🧹 Running event cleanup..."' >> /app/startup.sh && \
    echo './utils -cleanup-events' >> /app/startup.sh && \
    echo '' >> /app/startup.sh && \
    echo '# Wait for backend to continue running' >> /app/startup.sh && \
    echo 'echo "✅ Cleanup complete! Backend service running..."' >> /app/startup.sh && \
    echo 'wait $BACKEND_PID' >> /app/startup.sh && \
    chmod +x /app/startup.sh

CMD ["/app/startup.sh"]
