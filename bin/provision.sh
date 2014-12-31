#!/bin/bash

# Install vim
echo "Installing apt packages"
sudo apt-get update -qq
sudo apt-get install -qq -y vim curl git mercurial

# Install go
GO_VERSION="1.4"
GO_PKG="go${GO_VERSION}.linux-amd64.tar.gz"
GO_SHA="cd82abcb0734f82f7cf2d576c9528cebdafac4c6"
VAGRANT_SU="sudo su - vagrant -c 'echo foo'"

if [ -d "/usr/local/go" ]; then
	echo "Go is already installed"
else
	echo "Installing go "
	curl https://storage.googleapis.com/golang/${GO_PKG} > /tmp/${GO_PKG}
	if (( $? != 0 )); then
	  echo "fetching go package failed!"
	  exit 1
	fi

	file_sha=`sha1sum /tmp/${GO_PKG} | cut -f 1 -d " "`

	if (( $file_sha != $GO_SHA )); then
	  echo "sha1sum of ${GO_PKG} doesn't match!"
	  exit 1
	fi

	sudo tar -C /usr/local -xzf /tmp/${GO_PKG}
	if (( $? != 0 )); then
	  echo "sudo tar -C /usr/local -xzf /tmp/${GO_PKG} failed!"
	  exit 1
	fi

	cat <<'EOT' >> /home/vagrant/.bashrc
export GOROOT=/usr/local/go
export GOPATH=/home/vagrant/go
export PATH=${GOROOT}/bin:${GOPATH}/bin:$PATH
EOT

fi

export GOROOT=/usr/local/go
export GOPATH=/home/vagrant/go
export PATH=${GOROOT}/bin:${GOPATH}/bin:$PATH

sudo chown -R vagrant:vagrant /home/vagrant/go

# Install go packages
which godep > /dev/null 2>&1
if (( $? != 0 )); then
	echo "Installing godep"
	go get github.com/tools/godep
fi

which golint > /dev/null 2>&1
if (( $? != 0 )); then
	echo "Installing golint"
	go get github.com/golang/lint/golint
fi

echo "Restoring godep packages"
cd /home/vagrant/go/src/github.com/joshgarnett/sysminerd
godep restore
if (( $? != 0 )); then
  echo "godep restore failed!"
  exit 1
fi
