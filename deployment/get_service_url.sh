#!/bin/bash

set -e

eval "$(jq -r '@sh "PROJECT=\(.project) REGION=\(.region)"')"

while true
do
	URL=$(gcloud beta run services describe find-new-rent-service \
		--platform managed \
		--project $PROJECT \
		--region $REGION \
		--format json | jq --raw-output '.status.url // empty')

	if [ ! -z "$URL" ]
	then
		break
	fi
	sleep 5
done

echo "{\"url\": \"$URL\"}"
