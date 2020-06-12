#!/bin/bash
rm -rf /etc/pki/etcd-ca

MY_ADDRESS=`hostname -I`
MY_ADDRESS=`echo $MY_ADDRESS | sed -e 's/^ *//g' -e 's/ *$//g'`
ORG_DIR=`pwd`

# wget Insert
sudo apt -y install wget

# Etcd Install
export RELEASE="3.3.13"
wget https://github.com/etcd-io/etcd/releases/download/v${RELEASE}/etcd-v${RELEASE}-linux-amd64.tar.gz
tar xvf etcd-v${RELEASE}-linux-amd64.tar.gz

cd etcd-v${RELEASE}-linux-amd64
mv etcd etcdctl /usr/local/bin 

mkdir -p /var/lib/etcd/
mkdir /etc/etcd

# Change Directory
cd $ORG_DIR

# etcd Create CA 
cd /etc/pki
mkdir etcd-ca
cd etcd-ca

mkdir private certs newcerts crl
wget https://raw.githubusercontent.com/kelseyhightower/etcd-production-setup/master/openssl.cnf
touch index.txt
echo '01' > serial

apt-get install -y expect
expect <<-EOF
        set timeout 60
        spawn openssl req -config openssl.cnf -new -x509 -extensions v3_ca \
                          -keyout private/ca.key -out certs/ca.crt 
        expect {
                "Enter PEM pass phrase" {send "ketilinux\r";exp_continue}
                "Verifying - Enter PEM pass phrase" {send "ketilinux\r";exp_continue}
		"Country Name (2 letter code)" {send "\r"; exp_continue}
                "Common Name (FQDN)" {send "ca.etcd.example.com\r";exp_continue}
                "Organization Name (eg, company)" {send "\r";exp_continue}       
	        eof
	}
EOF

export SAN="IP:127.0.0.1, IP:${MY_ADDRESS}"
expect <<-EOF
        set timeout 60
        spawn openssl req -config openssl.cnf -new -nodes \
                          -keyout private/etcd0.example.com.key -out etcd0.example.com.csr
        expect {
                "Country Name (2 letter code)" {send "\r"; exp_continue}
                "Common Name (FQDN)" {send "etcd0.example.com\r";exp_continue}
                "Organization Name (eg, company)" {send "\r";exp_continue}
        	eof
	}
EOF
expect <<-EOF
        set timeout 60
        spawn openssl ca -config openssl.cnf -extensions etcd_server \
			  -keyfile private/ca.key \
			  -cert certs/ca.crt \
			  -out certs/etcd0.example.com.crt -infiles etcd0.example.com.csr
        expect {
                "Enter pass phrase for private/ca.key" {send "ketilinux\r"; exp_continue}
                "Sign the certificate" {send "y\r";exp_continue}
                "1 out of 1 certificate requests certified, commit?" {send "y\r";exp_continue}
	        eof
	}
EOF
openssl x509 -in certs/etcd0.example.com.crt -noout -text

unset SAN

expect <<-EOF
        set timeout 60
        spawn openssl req -config openssl.cnf -new -nodes \
			  -keyout private/etcd-client.key -out etcd-client.csr
        expect {
                "Country Name (2 letter code)" {send "\r"; exp_continue}
                "Common Name (FQDN)" {send "etcd-client\r";exp_continue}
                "Organization Name (eg, company)" {send "\r";exp_continue}
	        eof
	}
EOF

expect <<-EOF
	set timeout 60
	spawn openssl ca -config openssl.cnf -extensions etcd_client \
			 -keyfile private/ca.key \
			 -cert certs/ca.crt \
			 -out certs/etcd-client.crt -infiles etcd-client.csr
	expect {
		"Enter pass phrase for private/ca.key" {send "ketilinux\r"; exp_continue}
		"Sign the certificate?" {send "y\r"; exp_continue}
		"1 out of 1 certificate requests certified, commit?" {send "y\r"; exp_continue}
		eof
	}
EOF

# CD
cd $ORG_DIR

# Etcd Server Service Start
cp etcd.service /etc/systemd/system/etcd.service
sed -i "s/<IPADDRESS>/${MY_ADDRESS}/" /etc/systemd/system/etcd.service
systemctl daemon-reload
systemctl start etcd.service

# bashrc update
echo "export ETCDCTL_API=3" >> ~/.bashrc
echo "alias e='etcdctl --endpoints 127.0.0.1:2379   --cert /etc/pki/etcd-ca/certs/etcd-client.crt   --key /etc/pki/etcd-ca/private/etcd-client.key   --cacert /etc/pki/etcd-ca/certs/ca.crt'" >> ~/.bashrc

# Delete File
rm -rf *.gz











