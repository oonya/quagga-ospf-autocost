sudo ip netns add ue
sudo ip netns add proxy

sudo ip netns exec ue ip link set lo up
sudo ip netns exec proxy ip link set lo up

# 100.0
sudo ip link add name ranEth1 type veth peer name ranEth2
sudo ip link set ranEth1 netns ue

# 200.0
sudo ip link add name veth2 type veth peer name veth3
sudo ip link set veth3 netns proxy

# for wlan
sudo ip link add name wlanEth1 type veth peer name wlanEth3
sudo ip link set wlanEth1 netns ue
sudo ip link set wlanEth3 netns proxy

# add addr
sudo ip addr add 192.168.100.2/24 dev ranEth2
sudo ip link set dev ranEth2 up
sudo ip netns exec ue ip addr add 192.168.100.1/24 dev ranEth1
sudo ip netns exec ue sudo ip link set dev ranEth1 up

sudo ip addr add 192.168.200.2/24 dev veth2
sudo ip link set dev veth2 up
sudo ip netns exec proxy ip addr add 192.168.200.3/24 dev veth3
sudo ip netns exec proxy sudo ip link set dev veth3 up

sudo ip netns exec ue ip addr add 192.168.110.1/24 dev wlanEth1
sudo ip netns exec ue ip link set dev wlanEth1 up
sudo ip netns exec proxy ip addr add 192.168.110.3/24 dev wlanEth3
sudo ip netns exec proxy ip link set dev wlanEth3 up


sudo ip netns add server
sudo ip netns add waste

sudo ip netns exec server ip link set lo up

sudo ip link add name vnic type veth peer name  hide
sudo ip link set hide netns waste
sudo ip link set vnic netns ue

sudo ip link add name vserver type veth peer name  hide
sudo ip link set vserver netns proxy
sudo ip link set hide netns server

sudo ip netns exec ue ip addr add 192.168.1.1/24 dev vnic
sudo ip netns exec ue ip link set vnic up
sudo ip netns exec waste ip link set hide up

sudo ip netns exec proxy ip addr add 192.168.3.3/24 dev vserver
sudo ip netns exec server ip addr add 192.168.3.4/24 dev hide
sudo ip netns exec proxy ip link set vserver up
sudo ip netns exec server ip link set hide up

