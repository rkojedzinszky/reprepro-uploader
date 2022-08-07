#!/bin/sh

set -e

prepare_dot_ssh()
{
	echo "$DOT_SSH" | base64 -d | tar xzf - -C $HOME
	cat > $HOME/.ssh/config <<-EOF
	Host *
	 BatchMode=yes
	 CheckHostIP=no
	 StrictHostKeyChecking=yes
	EOF
}

reprepro_upload()
{
	cd ${REPREPRO_OUTPUT_PATH:-/output}
	echo ""
	echo "$(date) Uploading"
	echo ""
	ls -l | sed -e "s/^/  /"
	echo ""
	tar cf - * | ssh -p ${REPREPRO_PORT:-2222} -l ${REPREPRO_USER:-reprepro} ${REPREPRO_SERVER:-reprepro-server.ci-cd} /bin/upload.sh ${DISTRIBUTION}
}

prepare_dot_ssh
reprepro_upload
