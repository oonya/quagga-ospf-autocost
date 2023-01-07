# ue
sudo modprobe ip_gre
# 5GSAなら10.45.0.2
sudo ip netns exec ue ip tunnel add rangre mode gre remote 192.168.200.3 local 10.60.0.1  ttl 255
sudo ip netns exec ue ip addr add 10.10.10.1/24 dev rangre
sudo ip netns exec ue ip link set rangre up

sudo ip netns exec ue ip tunnel add wlangre mode gre remote 192.168.120.3 local 192.168.120.1  ttl 255
sudo ip netns exec ue ip addr add 10.10.20.1/24 dev wlangre
sudo ip netns exec ue ip link set wlangre up

# proxy 
# 5GSAなら10.45.0.2
sudo ip netns exec proxy ip tunnel add rangre mode gre remote 10.60.0.1  local 192.168.200.3  ttl 255
sudo ip netns exec proxy ip addr add 10.10.10.3/24 dev rangre
sudo ip netns exec proxy ip link set rangre up

sudo ip netns exec proxy ip tunnel add wlangre mode gre remote 192.168.120.1 local 192.168.120.3  ttl 255
sudo ip netns exec proxy ip addr add 10.10.20.3/24 dev wlangre
sudo ip netns exec proxy ip link set wlangre up

