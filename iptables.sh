#!/bin/sh
export PATH="/bin:/sbin:/usr/sbin:/usr/bin"

DEV="eth1"
GWPrefix="192.168"
SubNet="192.168.42.0/24"

IPSECGATEWAY=$(ip addr show $DEV| grep $GWPrefix | sed 's/.*inet *\([0-9.]*\).*/\1/')

if [ -z "$IPSECGATEWAY" ]; then
    exit 1
fi

iptables -t nat --flush

. ./iptables.rule

iptables -t nat -A POSTROUTING -s $SubNet -o $DEV -j SNAT --to-source $IPSECGATEWAY

# iptables -t nat -nvL --line-number
