#!/bin/bash

# Install vim
echo "Installing apt packages"
apt-get update -qq
apt-get install -qq -y vim curl git mercurial

echo "Remove unneeded packages"
apt-get remove -qq -y puppet puppet-common chef chef-zero
