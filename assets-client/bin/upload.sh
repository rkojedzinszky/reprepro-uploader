#!/bin/sh

set -e

cd ${REPREPRO_OUTPUT_PATH:-/output}
echo ""
echo "$(date) Uploading"
echo ""
ls -l | sed -e "s/^/  /"
echo ""
tar czf - * | curl -s --upload-file - -H "Authorization: bearer $REPREPRO_TOKEN" -H "Content-encoding: gzip" "$REPREPRO_SERVER?dists=$DISTRIBUTION"
