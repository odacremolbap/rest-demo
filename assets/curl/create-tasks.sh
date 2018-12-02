
#!/bin/bash

HOST=${HOST:-localhost}
PORT=${PORT:-9101}

curl -X POST \
    http://${HOST}:${PORT}/v1/tasks \
    -H "Content-Type: application/json" \
    -d '{
        "name": "first task",
        "due_date": "2018-12-31T23:59:59Z"
        }' \
    | jq

