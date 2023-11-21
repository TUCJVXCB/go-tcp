package main

import (
	"go-tcp/tcpip/lab/link/raw"
	"go-tcp/tcpip/lab/link/tuntap"
	"log"
)

func main() {
	tapName := "tap0"
	c := &tuntap.Config{tapName, tuntap.TAP}
	fd, err := tuntap.NewNetDev(c)

	if err != nil {
		panic(err)
	}

	_ = tuntap.SetLinkUp(tapName)
	_ = tuntap.SetRoute(tapName, "192.168.1.0/24")

	buf := make([]byte, 1<<16)
	for {
		rn, err := raw.BlockingRead(fd, buf)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("read %d bytes", rn)
	}
}
