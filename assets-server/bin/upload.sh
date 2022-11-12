#!/bin/sh

set -e

[ $# -gt 0 ] && REPREPRO_REPOS="$@"

cd reprepro

tmpdir=$(mktemp -d)
trap "rm -rf $tmpdir" EXIT

tar xf - -C "$tmpdir"
for chg in $tmpdir/*.changes ; do
	test -f "$chg" || continue

	echo "Processing $(basename "$chg")"
	for repo in ${REPREPRO_REPOS}; do
		echo " Including in $repo"
		reprepro --waitforlock=${REPREPRO_WAITFORLOCK:-30} --export=never --keepunreferencedfiles include $repo "$chg"
	done
	echo " Done"
	echo ""
done

reprepro --waitforlock=${REPREPRO_WAITFORLOCK:-30} export
