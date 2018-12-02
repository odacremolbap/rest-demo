
#!/bin/bash

HOST=${HOST:-localhost}
PORT=${PORT:-9101}
TASK_ID=${TASK_ID:-1}
curl -X DELETE \
    http://${HOST}:${PORT}/v1/tasks/${TASK_ID} \
    | jq

