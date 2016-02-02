#!/bin/bash
set +e
set -o noglob

#
# Set Colors
#

reset=$(tput sgr0)
red=$(tput setaf 1)
white=$(tput setaf 7)

#
# Headers and Logging
#

info() { printf "${white}➜ %s${reset}\n" "$@"
}
error() { printf "${red}✖ %s${reset}\n" "$@"
}

main() {

  info 'starting ecs deploy'

  # Check variables
  if [ -z "$WERCKER_AWS_ECS_KEY" ]; then
    error "Please set the 'key' variable"
    exit 1
  fi

  if [ -z "$WERCKER_AWS_ECS_SECRET" ]; then
    error "Please set the 'secret' variable"
    exit 1
  fi

  if [ -z "$WERCKER_AWS_ECS_CLUSTER_NAME" ]; then
    error "Please set the 'cluster-name' variable"
    exit 1
  fi

  if [ -z "$WERCKER_AWS_ECS_SERVICE_NAME" ]; then
    error "Please set the 'service-name' variable"
    exit 1
  fi

  if [ -z "$WERCKER_AWS_ECS_TASK_DEFINITION_NAME" ]; then
    error "Please set the 'task-definition-name' variable"
    exit 1
  fi

  if [ -z "$WERCKER_AWS_ECS_TASK_DEFINITION_FILE" ]; then
    error "Please set the 'task-definition-file' variable"
    exit 1
  fi

  ./ecs \
  -key "$WERCKER_AWS_ECS_KEY" \
  -secret "$WERCKER_AWS_ECS_SECRET" \
  -region "${WERCKER_AWS_ECS_REGION:-us-east-1}" \
  -cluster-name "$WERCKER_AWS_ECS_CLUSTER_NAME" \
  -service-name "$WERCKER_AWS_ECS_SERVICE_NAME" \
  -task-definition-name "$WERCKER_AWS_ECS_TASK_DEFINITION_NAME" \
  -task-definition-file "$WERCKER_AWS_ECS_TASK_DEFINITION_FILE" \
  -minimum-running-tasks "${WERCKER_AWS_ECS_MINIMUM_RUNNING_TASKS:-1}"
}

main