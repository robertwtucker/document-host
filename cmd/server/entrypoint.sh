#!/bin/sh
APP_ENV=${APP_ENV:-local}

echo "Running entrypoint script in the '${APP_ENV}' environment..."
CONFIG_FILE=/app/config/${APP_ENV}.yaml

/app/server -config ${CONFIG_FILE}
