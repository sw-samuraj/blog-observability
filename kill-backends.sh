#!/bin/sh -eu

LOG_FILE="_logs/observability.log"

kill_app () {
  APP_PID=$(pgrep "${APP}" || echo "" )
  if [ -z "${APP_PID}" ]
  then
    echo "${APP} is not running."
  else
    echo "Killing ${APP} running with PID ${APP_PID}."
    kill -9 "${APP_PID}"
  fi
}

APP="loki"
kill_app
APP="promtail"
kill_app
APP="prometheus"
kill_app

if [ -f "${LOG_FILE}" ]
then
  echo "Deleting the old my-app log file..."
  rm "${LOG_FILE}"
fi
