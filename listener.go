package ss
// Is ss too short? 
// Oh well.. http://golang.org/doc/effective_go.html#package-names
// Alias it if it bugs you

import (
  "errors"
  "fmt"
  "net"
  "os"
  "strconv"
  "strings"
)

type ListenTarget struct {
  Name  string // host:port | port. Currently unix sockets are not supported
  Fd    uintptr
}

func Ports() ([]ListenTarget, error) {
  ssport := os.Getenv("SERVER_STARTER_PORT")
  if ssport == "" {
    return nil, errors.New("No environment variable SERVER_STARTER_PORT available")
  }

  return ParsePorts(ssport)
}

func ParsePorts(ssport string) ([]ListenTarget, error) {
  ret := []ListenTarget{}
  for _, pairstring := range strings.Split(ssport, ";") {
    pair := strings.Split(pairstring, "=")
    port, err := strconv.ParseUint(pair[1], 10, 0)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Failed to parse '%s'", pairstring))
    }
    ret = append(ret, ListenTarget { pair[0], uintptr(port) })
  }
  return ret, nil
}

func NewListener() (net.Listener, error) {
  portmap, err := Ports()
  if err != nil {
    return nil, err
  }
  return ListenOn(portmap[0])
}

func AllListeners() ([]net.Listener, error) {
  portmap, err := Ports()
  if err != nil {
    return nil, err
  }
  return ListenersIn(portmap)
}

func ListenersIn (list []ListenTarget) ([]net.Listener, error) {
  ret := []net.Listener {}
  for _, t := range list {
    l, err := ListenOn(t)
    if err != nil {
      return nil, err
    }
    ret = append(ret, l)
  }
  return ret, nil
}

func ListenOn (t ListenTarget) (net.Listener, error) {
  f := os.NewFile(t.Fd, "foo")
  return net.FileListener(f)
}
