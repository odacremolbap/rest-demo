#!/bin/bash

# This relies on environemnt variable loaded from $(REPO_ROOT)/assets/run/environment
# to work on targeted development platforms

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
source $REPO_ROOT/assets/run/environment

BIN=todolist
MAIN=$REPO_ROOT/cmd/$BIN/main.go

# All flag values need to be defined at the environment file
go run $MAIN \
  server \
  --port $SERVER_PORT \
  --db-host $DB_HOST \
  --db-port $DB_PORT \
  --db-user $DB_USER \
  --db-password $DB_PASSWORD \
  --db-name $DB_NAME \
  --log-formatter $LOG_FORMATTER \
  --v $VERBOSITY \  
