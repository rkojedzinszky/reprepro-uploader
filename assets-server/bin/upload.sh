#!/bin/sh

set -e

[ $# -gt 0 ] && REPREPRO_REPOS="$@"

cd reprepro

tmpdir=$(mktemp -d)
trap "rm -rf $tmpdir" EXIT

tar xf - -C "$tmpdir"
for chg in $tmpdir/*.changes ; do
	test -f "$chg" || continue

	for repo in ${REPREPRO_REPOS}; do
		reprepro --waitforlock=${REPREPRO_WAITFORLOCK:-30} --export=never --keepunreferencedfiles include $repo "$chg"
	done
done

reprepro --waitforlock=${REPREPRO_WAITFORLOCK:-30} export
