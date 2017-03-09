#!/bin/bash

DEV="eth1"
GRE="gre1"
GRE_TABLE="vpn"
DEF_TABLE="main"
SubNet="192.168.3.0/24"
NODE_A_PUB="104.207.150.187"
NODE_A_PRI="192.168.0.1/30"
NODE_B_PUB="60.205.140.248"
NODE_B_PRI="192.168.0.2/30"
MODE="A"

# chdir
cd /data/sh

# gre tunnel
# echo "200 vpn" >> /etc/iprouter/rt_table
ip tunnel add $GRE mode gre remote $NODE_B_PUB local $NODE_A_PUB ttl 255
ip link set $GRE up
ip addr add $NODE_A_PRI peer $NODE_B_PRI dev $GRE
if [ "$MODE" = "A" ]; then
    ip route add 8.8.0.0/16 dev gre1
fi

# vpn
/usr/local/sbin/ipsec start

# vpn route
ip route add default dev $GRE table $GRE_TABLE

. ./ip.rule

# iptables -A PREROUTING -t mangle -s $SubNet -j MARK --set-mark 3
# ip rule add fwmark 3 table $GRE_TABLE
ip rule add from $SubNet table $GRE_TABLE pref 10001
iptables -t nat -A POSTROUTING -s $SubNet -j MASQUERADE
