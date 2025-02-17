#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright 2021 Authors of KubeArmor

# update repo
sudo sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-Linux-*
sudo sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-Linux-*

# make a directory to build bcc
mkdir -p /tmp/build; cd /tmp/build

# download bcc
sudo dnf -y install git
git -C /tmp/build/ clone --branch v0.24.0 --depth 1 https://github.com/iovisor/bcc.git

# install dependencies for bcc
sudo dnf -y install gcc gcc-c++ make cmake bison flex ethtool iperf3 \
                    python3-netaddr python3-pip zlib-devel elfutils-libelf-devel libarchive \
                    wget zip clang clang-devel llvm llvm-devel llvm-static ncurses-devel

# install bcc
mkdir -p /tmp/build/bcc/build; cd /tmp/build/bcc/build
cmake .. -DENABLE_LLVM_SHARED=1 && make -j$(nproc) && sudo make install
if [ $? != 0 ]; then
    echo "Failed to install bcc"
    exit 1
fi

# install dependencies for selinux
sudo dnf -y install policycoreutils-devel setools-console

if [[ $(hostname) = kubearmor-dev* ]]; then
    echo >> /home/vagrant/.bashrc
    echo "alias lz='ls -lZ'" >> /home/vagrant/.bashrc
    echo >> /home/vagrant/.bashrc
    mkdir -p /home/vagrant/go; chown -R vagrant:vagrant /home/vagrant/go
elif [ -z "$GOPATH" ]; then
    echo >> ~/.bashrc
    echo "alias lz='ls -lZ'" >> ~/.bashrc
    echo >> ~/.bashrc
fi

# enable audit mode
sudo semanage dontaudit off

# install golang
echo "Installing golang binaries..."
goBinary=$(curl -s https://go.dev/dl/ | grep linux | head -n 1 | cut -d'"' -f4 | cut -d"/" -f3)
wget --quiet https://dl.google.com/go/$goBinary -O /tmp/build/$goBinary
sudo tar -C /usr/local -xzf /tmp/build/$goBinary

if [[ $(hostname) = kubearmor-dev* ]]; then
    echo >> /home/vagrant/.bashrc
    echo "export GOPATH=\$HOME/go" >> /home/vagrant/.bashrc
    echo "export GOROOT=/usr/local/go" >> /home/vagrant/.bashrc
    echo "export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin" >> /home/vagrant/.bashrc
    echo >> /home/vagrant/.bashrc
    mkdir -p /home/vagrant/go; chown -R vagrant:vagrant /home/vagrant/go
elif [ -z "$GOPATH" ]; then
    echo >> ~/.bashrc
    echo "export GOPATH=\$HOME/go" >> ~/.bashrc
    echo "export GOROOT=/usr/local/go" >> ~/.bashrc
    echo "export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin" >> ~/.bashrc
    echo >> ~/.bashrc
fi

# download protoc
mkdir -p /tmp/build/protoc; cd /tmp/build/protoc
wget --quiet https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_64.zip -O /tmp/build/protoc/protoc-3.19.4-linux-x86_64.zip

# install protoc
unzip protoc-3.19.4-linux-x86_64.zip
sudo mv bin/protoc /usr/local/bin/
sudo chmod 755 /usr/local/bin/protoc

# apply env
if [[ $(hostname) = kubearmor-dev* ]]; then
    export GOPATH=/home/vagrant/go
    export GOROOT=/usr/local/go
    export PATH=$PATH:/usr/local/go/bin:/home/vagrant/go/bin
elif [ -z "$GOPATH" ]; then
    export GOPATH=$HOME/go
    export GOROOT=/usr/local/go
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
fi

# download protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# install kubebuilder
wget --quiet https://github.com/kubernetes-sigs/kubebuilder/releases/download/v3.1.0/kubebuilder_linux_amd64 -O /tmp/build/kubebuilder
chmod +x /tmp/build/kubebuilder; sudo mv /tmp/build/kubebuilder /usr/local/bin

if [[ $(hostname) = kubearmor-dev* ]]; then
    echo >> /home/vagrant/.bashrc
    echo 'export PATH=$PATH:/usr/local/kubebuilder/bin' >> /home/vagrant/.bashrc
elif [ -z "$GOPATH" ]; then
    echo >> ~/.bashrc
    echo 'export PATH=$PATH:/usr/local/kubebuilder/bin' >> ~/.bashrc
fi

# install kustomize
cd /tmp/build/
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
sudo mv kustomize /usr/local/bin

# remove downloaded files
cd; sudo rm -rf /tmp/build
