sudo apt install -y quagga
sudo systemctl disable ospfd.service
sudo systemctl disable zebra.service

sudo mkdir /etc/quagga/ue
sudo mkdir /etc/quagga/proxy

sudo cp /usr/share/doc/quagga-core/examples/zebra.conf.sample /etc/quagga/ue/zebra.conf
sudo cp /usr/share/doc/quagga-core/examples/zebra.conf.sample /etc/quagga/proxy/zebra.conf
sudo cp /usr/share/doc/quagga-core/examples/vtysh.conf.sample /etc/quagga/ue/vtysh.conf
sudo cp /usr/share/doc/quagga-core/examples/vtysh.conf.sample /etc/quagga/proxy/vtysh.conf
sudo cp /usr/share/doc/quagga-core/examples/ospfd.conf.sample /etc/quagga/ue/ospfd.conf
sudo cp /usr/share/doc/quagga-core/examples/ospfd.conf.sample /etc/quagga/proxy/ospfd.conf

sudo mkdir /var/run/quagga
sudo chmod g+w /var/run/quagga/
sudo chown root:quagga /var/run/quagga/

sudo chown root:root /etc/quagga/ue/*
sudo chmod +r /etc/quagga/ue/zebra.conf
sudo chown root:root /etc/quagga/proxy/*
sudo chmod +r /etc/quagga/proxy/zebra.conf

sudo mkdir /etc/quagga/server
sudo cp /usr/share/doc/quagga-core/examples/zebra.conf.sample /etc/quagga/server/zebra.conf
sudo cp /usr/share/doc/quagga-core/examples/vtysh.conf.sample /etc/quagga/server/vtysh.conf
sudo cp /usr/share/doc/quagga-core/examples/ospfd.conf.sample /etc/quagga/server/ospfd.conf

sudo chown root:root /etc/quagga/server/*
sudo chmod +r /etc/quagga/server/zebra.conf

