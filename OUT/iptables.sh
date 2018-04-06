#!/bin/sh
export PATH="/bin:/sbin:/usr/sbin:/usr/bin"

DEV="eth1"
GWPrefix="192.168"
SubNet="192.168.12.0/24"

TUNGATEWAY=$(ip addr show $DEV| grep $GWPrefix | sed 's/.*inet *\([0-9.]*\).*/\1/')

if [ -z "$TUNGATEWAY" ]; then
    exit 1
fi

iptables -t nat --flush
iptables -t nat -A POSTROUTING -s $SubNet -o $DEV -j MASQUERADE
#iptables -t nat -A POSTROUTING -s $SubNet -o $DEV -j SNAT --to-source $TUNGATEWAY

. ./iptables.rule

iptables -t nat -D POSTROUTING -s $SubNet -o $DEV -j MASQUERADE
iptables -t nat -A POSTROUTING -s $SubNet -o $DEV -j MASQUERADE
#iptables -t nat -D POSTROUTING -s $SubNet -o $DEV -j SNAT --to-source $TUNGATEWAY
#iptables -t nat -A POSTROUTING -s $SubNet -o $DEV -j SNAT --to-source $TUNGATEWAY

# iptables -t nat -nvL --line-number
