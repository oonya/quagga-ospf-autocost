# host
sudo modprobe mac80211_hwsim radios=2 dyndbg=+p
sudo apt install -y hostapd
sudo apt install -y wpasupplicant

sudo iw phy phy1 set netns name /run/netns/ue
sudo iw phy phy0 set netns name /run/netns/proxy

sudo ip netns exec proxy ip link set dev wlan0 up
sudo ip netns exec proxy hostapd -B -f /home/oonya/wlan/hostapd.log -i wlan0 /home/oonya/quagga-ospf-autocost/etc/wlan/hostapd.conf
sudo ip netns exec proxy ip addr add 192.168.120.3/24 dev wlan0

sudo ip netns exec ue wpa_supplicant -B -c /home/oonya/quagga-ospf-autocost/etc/wlan/wpa_supplicant.conf -f /home/oonya/wlan/wpa_supplicant.log -i wlan1
sudo ip netns exec ue ip addr add 192.168.120.1/24 dev wlan1
