# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  # Every Vagrant virtual environment requires a box to build off of.
  config.vm.box = "ubuntu/trusty64"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  config.vm.network "private_network", ip: "192.168.33.50"

  config.vm.hostname = "sysminerd.local"

  # If true, then any SSH connections made will enable agent forwarding.
  # Default value: false
  config.ssh.forward_agent = true

  # Share an additional folder to the guest VM.
  config.vm.synced_folder "./", "/home/vagrant/go/src/github.com/joshgarnett/sysminerd"

  config.vm.provider "virtualbox" do |v|
    v.memory = 768
    v.cpus = 2
  end

  # Setup all go dependencies
  config.vm.provision "shell" do |s|
    s.path = "bin/provision.sh"
    s.privileged = false
  end

  # Install graphite for testing
  config.vm.provision "shell" do |s|
    s.path = "bin/install_graphite.sh"
    s.privileged = true
  end

  # Install graphite for testing
  config.vm.provision "shell" do |s|
    s.path = "bin/install_grafana.sh"
    s.privileged = true
  end
end
