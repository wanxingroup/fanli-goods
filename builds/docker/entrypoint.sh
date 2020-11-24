#!/usr/bin/env sh

configFile=${CONFIG_FILE:-/app/config.json}
logLevel=${LOG_LEVEL:-info}
serviceId=${SERVICE_ID:-0}

/app/service --serviceId="${serviceId}" --config="${configFile}"
