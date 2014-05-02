package ss

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// ListenTarget describes an address and an associated file descriptor
type ListenTarget struct {
	Name string // host:port | port. Currently unix sockets are not supported
	Fd   uintptr
}

/*

Ports parses SERVER_STARTER_PORT environment variable, and returns a list of
ListenTarget structs that can be passed to NewListenerOn()

*/
func Ports() ([]ListenTarget, error) {
	ssport := os.Getenv("SERVER_STARTER_PORT")
	if ssport == "" {
		return nil, fmt.Errorf("error: No environment variable SERVER_STARTER_PORT available")
	}

	return ParsePorts(ssport)
}

/*

ParsePorts parses the given string and returns a list of
ListenTarget structs that can be passed to NewListenerOn()

*/
func ParsePorts(ssport string) ([]ListenTarget, error) {
	ret := []ListenTarget{}
	for _, pairstring := range strings.Split(ssport, ";") {
		pair := strings.Split(pairstring, "=")
		port, err := strconv.ParseUint(pair[1], 10, 0)
		if err != nil {
			return nil, fmt.Errorf("error: Failed to parse '%s'", pairstring)
		}
		ret = append(ret, ListenTarget{pair[0], uintptr(port)})
	}
	return ret, nil
}

/*

NewListenerOrDefault creates a new listener from SERVER_STARTER_PORT, or
if that fails, binds to the "default" bind address

*/
func NewListenerOrDefault(proto, defaultBindAddress string) (net.Listener, error) {
	l, err := NewListener()
	if err == nil {
		return l, nil
	}

	dl, err := net.Listen(proto, defaultBindAddress)
	if err == nil {
		return dl, nil
	}
	return nil, err
}

/*

NewListener creates a new listener from SERVER_STARTER_PORT environment variable

Note that this binds to only ONE file descriptor.
If multiple ports are specified in the environment variable, the first one is used

*/
func NewListener() (net.Listener, error) {
	portmap, err := Ports()
	if err != nil {
		return nil, err
	}
	return NewListenerOn(portmap[0])
}

/*

AllListeners creates new listeners from SERVER_STARTER_PORT environment
variable.

This binds to ALL file descriptors in SERVER_STARTER_PORT

*/
func AllListeners() ([]net.Listener, error) {
	portmap, err := Ports()
	if err != nil {
		return nil, err
	}
	return NewListenersOn(portmap)
}

/*

NewListenersOn creates listeners for each ListenTarget given

*/
func NewListenersOn(list []ListenTarget) ([]net.Listener, error) {
	ret := []net.Listener{}
	for _, t := range list {
		l, err := NewListenerOn(t)
		if err != nil {
			return nil, err
		}
		ret = append(ret, l)
	}
	return ret, nil
}

/*

NewListenerOn creates a listener for given ListenTarget

*/
func NewListenerOn(t ListenTarget) (net.Listener, error) {
	f := os.NewFile(t.Fd, t.Name)
	return net.FileListener(f)
}
