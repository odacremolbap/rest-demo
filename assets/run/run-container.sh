#!/bin/bash

# This will relies on environemnt variable loaded from $(REPO_ROOT)/assets/run/environment
# to work on targeted development platforms

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
source $REPO_ROOT/assets/run/environment

OS=`go env GOOS`
OUTPUT_DIR=_output
BIN=todolist

EXEC=$REPO_ROOT/$OUTPUT_DIR/$OS/$ARCH/$BIN

# All flag values need to be defined at the environment file

echo "TODO - create golang container + volume map + host network"
exit 1

$EXEC \
  server \
  --port $SERVER_PORT \
  --db-host $DB_HOST \
  --db-port $DB_PORT \
  --db-user $DB_USER \
  --db-password $DB_PASSWORD \
  --db-name $DB_NAME \
  --log-formatter $LOG_FORMATTER \
  --v $VERBOSITY  
