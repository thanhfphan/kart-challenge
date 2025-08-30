#!/usr/bin/env bash

set -e 
OS=$(uname | tr '[:upper:]' '[:lower:]')
if [[ "$OS" == *"mingw"* ]] || [[ "$OS" == *"windows"* ]]; then
   ROOT_DIR=$(pwd -W) 
else
   ROOT_DIR=$(pwd)
fi
PROGNAME="$(basename $0)"

function help() {
  echo 1>&2 "Usage: ${PROGNAME} <command>"
  echo 1>&2 ""
  echo 1>&2 "Commands:"
  echo 1>&2 "  start        start the service"
  echo 1>&2 "  up           pull and start infrastructure images"
  echo 1>&2 "  down         stop all infrastructure images"
  echo 1>&2 "  tests        run unit test"
  echo 1>&2 "  lint         check lint"
  echo 1>&2 "  migrate		  run migration"
}

function setup_env() {
    set -a
    export $(grep -v '^#' "$ROOT_DIR/builders/.base.env" | xargs) >/dev/null 2>&1
    set +a
}

function start() {
    setup_env
    go run ./cmd/api/main.go
}

function up() {
	docker compose -f ./builders/docker-compose.dev.yml up -d
}

function down() {
	docker compose -f ./builders/docker-compose.dev.yml down
}

function run_test() {
  echo 'Run unit testing'
  go test ./... -short || {
    echo 'Unit testing failed'
    exit 1
  }
}

function migrate() {
	setup_env
	go run ./cmd/migrate/main.go
}

function lint() {
    setup_env
    echo "Running linter..."
    # Add linter logic here
}

SUBCOMMAND="${1:-}"
case "${SUBCOMMAND}" in
  "" | "help" | "-h" | "--help" )
    help
    ;;

  "start" )
    shift
    start "$@"
    ;;

  "up" )
    shift
    up "$@"
    ;;

  "down" )
    shift
    down "$@"
    ;;

  "tests" )
    shift
    run_test"$@"
    ;;

  "lint" )
    shift
    lint "$@"
    ;;

  "migrate" )
    shift
    migrate "$@"
    ;;
  
  *)
    help
    exit 1
    ;;
esac
