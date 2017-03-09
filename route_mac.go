package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	ipfile = "delegated-apnic-latest.txt"
	ipurl  = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"
	upBase = `#!/bin/sh
export PATH="/bin:/sbin:/usr/sbin:/usr/bin"

OLDGW=$(netstat -nr | grep '^default' | grep -v 'ppp' | grep -v '::' | sed 's/default *\([0-9.]*\).*/\1/')

if [ ! -e /tmp/pptp_oldgw ]; then
    echo "${OLDGW}" > /tmp/pptp_oldgw
fi

dscacheutil -flushcache

`
	downBase = `#!/bin/sh
export PATH="/bin:/sbin:/usr/sbin:/usr/bin"

if [ ! -e /tmp/pptp_oldgw ]; then
        exit 0
fi

OLDGW=$(cat /tmp/pptp_oldgw)

`
)

func main() {
	err := download(ipfile, ipurl)
	if err != nil {
		log.Fatalf("ip route file download error: %v\n", err)
	}

	f, err := os.Open(ipfile)
	if err != nil {
		log.Fatalf("ip route file open error: %v\n", err)
	}
	defer f.Close()

	router := make([]string, 0)
	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("file readLine error: %v\n", err)
		}
		if strings.Contains(string(line), "CN|ipv4") == false {
			continue
		}
		columns := strings.Split(string(line), "|")
		netaddr := columns[3]
		ipcount, _ := strconv.Atoi(columns[4])

		i := 0
		n := 256
		for ipcount > n {
			n *= 2
			i++
		}
		router = append(router, fmt.Sprintf("%s/%d", netaddr, 24-i))
	}

	err = writeRouter(router)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("router parse success.")
	fmt.Println("copy ip-up and ip-down to /etc/ppp/ and don't forget to make them executable.")
}

func writeRouter(router []string) error {
	var err error
	// ip-up
	upData := make([]string, 0)
	for i := range router {
		upData = append(upData, fmt.Sprintf("route add %s \"${OLDGW}\"", router[i]))
	}
	err = writeFile("ip-up", []byte(upBase+strings.Join(upData, "\n")))
	if err != nil {
		return err
	}

	// ip-down
	downData := make([]string, 0)
	for i := range router {
		downData = append(downData, fmt.Sprintf("route delete %s \"${OLDGW}\"", router[i]))
	}
	downData = append(downData, "\n\nrm /tmp/pptp_oldgw\n")
	err = writeFile("ip-down", []byte(downBase+strings.Join(downData, "\n")))
	if err != nil {
		return err
	}

	return nil
}

func download(file, url string) error {
	if _, err := os.Stat(file); err == nil {
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	return writeFile(file, body)
}

func writeFile(file string, body []byte) error {
	f, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		return err
	}
	if _, err = f.Write(body); err != nil {
		return err
	}
	return nil
}
