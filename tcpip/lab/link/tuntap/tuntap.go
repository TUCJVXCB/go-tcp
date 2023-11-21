package tuntap

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	TUN = 1
	TAP = 2
)

var (
	ErrDeviceMode = errors.New("unsupport device mode")
)

type rawSockaddr struct {
	Family uint16
	Data   [14]byte
}

type Config struct {
	Name string
	Mode int // tap/tun
}

func NewNetDev(c *Config) (fd int, err error) {
	switch c.Mode {
	case TUN:
		fd, err = newTUN(c.Name)
	case TAP:
		fd, err = newTAP(c.Name)
	default:
		err = ErrDeviceMode
		return
	}
	if err != nil {
		return
	}
	return
}

func SetLinkUp(name string) (err error) {
	output, err := exec.Command("ip", "link", "set", name, "up").CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v%v", err, string(output))
		return
	}
	return
}

func AddIP(name string, ip string) (err error) {
	out, cmdErr := exec.Command("ip", "addr", "add", ip, "dev", name).CombinedOutput()
	if cmdErr != nil {
		err = fmt.Errorf("%v:%v", cmdErr, string(out))
		return
	}
	return
}

func newTUN(name string) (int, error) {
	return open(name, syscall.IFF_TUN|syscall.IFF_NO_PI)
}

func newTAP(name string) (int, error) {
	return open(name, syscall.IFF_TAP|syscall.IFF_NO_PI)
}

func open(name string, flags uint16) (int, error) {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR, 0)
	if err != nil {
		return -1, err
	}

	var ifr struct {
		name  [16]byte
		flags uint16
		_     [22]byte
	}
	copy(ifr.name[:], name)
	ifr.flags = flags
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		syscall.Close(fd)
		return -1, errno
	}
	return fd, nil
}

func SetRoute(name, cidr string) (err error) {
	// ip route add 192.168.1.0/24 dev tap0
	out, cmdErr := exec.Command("ip", "route", "add", cidr, "dev", name).CombinedOutput()
	if cmdErr != nil {
		err = fmt.Errorf("%v:%v", cmdErr, string(out))
		return
	}
	return
}
