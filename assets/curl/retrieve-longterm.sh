
#!/bin/bash

HOST=${HOST:-localhost}
PORT=${PORT:-9101}

curl http://${HOST}:${PORT}/v1/tasks?category=longterm | jq