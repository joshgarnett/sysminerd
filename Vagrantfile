# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.define "ubuntu" do |ubuntu|
    # Every Vagrant virtual environment requires a box to build off of.
    ubuntu.vm.box = "ubuntu/trusty64"

    # Create a private network, which allows host-only access to the machine
    # using a specific IP.
    ubuntu.vm.network "private_network", ip: "192.168.33.50"

    ubuntu.vm.hostname = "sysminerd-ubuntu.local"

    # If true, then any SSH connections made will enable agent forwarding.
    # Default value: false
    ubuntu.ssh.forward_agent = true

    # Share an additional folder to the guest VM.
    ubuntu.vm.synced_folder "./", "/home/vagrant/go/src/github.com/joshgarnett/sysminerd"

    ubuntu.vm.provider "virtualbox" do |v|
      v.memory = 768
      v.cpus = 2
    end

    # Setup all go dependencies
    ubuntu.vm.provision "shell" do |s|
      s.path = "bin/ubuntu/provision.sh"
      s.privileged = true
    end

    # Install go
    ubuntu.vm.provision "shell" do |s|
      s.path = "bin/install_go.sh"
      s.privileged = false
    end

    # Install graphite for testing
    ubuntu.vm.provision "shell" do |s|
      s.path = "bin/ubuntu/install_graphite.sh"
      s.privileged = true
    end

    # Install graphite for testing
    ubuntu.vm.provision "shell" do |s|
      s.path = "bin/ubuntu/install_grafana.sh"
      s.privileged = true
    end
  end

  config.vm.define "centos" do |centos|
    # Every Vagrant virtual environment requires a box to build off of.
    centos.vm.box = "chef/centos-6.5"

    # Create a private network, which allows host-only access to the machine
    # using a specific IP.
    centos.vm.network "private_network", ip: "192.168.33.51"

    centos.vm.hostname = "sysminerd-centos.local"

    # If true, then any SSH connections made will enable agent forwarding.
    # Default value: false
    centos.ssh.forward_agent = true

    # Share an additional folder to the guest VM.
    centos.vm.synced_folder "./", "/home/vagrant/sysminerd"

    centos.vm.provider "virtualbox" do |v|
      v.memory = 256
      v.cpus = 1
    end

  end

end
