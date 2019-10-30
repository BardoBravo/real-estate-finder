#!/bin/bash

set -e

TRIGGER_ID=$(gcloud alpha builds triggers list \
	--filter='triggerTemplate.repoName:github_bardobravo_finder-imob' \
	--format json | jq --raw-output '.[0].id')

gcloud alpha builds triggers run $TRIGGER_ID \
	--branch master > /dev/null
