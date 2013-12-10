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

/*
 * Parses SERVER_STARTER_PORT environment variable, and returns a list of
 * of ListenTarget structs that can be passed to NewListenerOn()
 */
func Ports() ([]ListenTarget, error) {
  ssport := os.Getenv("SERVER_STARTER_PORT")
  if ssport == "" {
    return nil, errors.New("No environment variable SERVER_STARTER_PORT available")
  }

  return ParsePorts(ssport)
}

/*
 * Parses the given string and returns a list of
 * of ListenTarget structs that can be passed to NewListenerOn()
 */
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

/*
 * Creates a new listener from SERVER_STARTER_PORT environment variable
 *
 * Note that this binds to only ONE file descriptor (the first one found)
 */
func NewListener() (net.Listener, error) {
  portmap, err := Ports()
  if err != nil {
    return nil, err
  }
  return NewListenerOn(portmap[0])
}

/* 
 * Creates new listeners from SERVER_STARTER_PORT environment variable.
 *
 * This binds to ALL file descriptors in SERVER_STARTER_PORT
 */
func AllListeners() ([]net.Listener, error) {
  portmap, err := Ports()
  if err != nil {
    return nil, err
  }
  return NewListenersOn(portmap)
}

/*
 * Given a list of ListenTargets, creates listeners for each one
 */
func NewListenersOn (list []ListenTarget) ([]net.Listener, error) {
  ret := []net.Listener {}
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
 * Given a ListenTarget, creates a listener
 */
func NewListenerOn (t ListenTarget) (net.Listener, error) {
  f := os.NewFile(t.Fd, t.Name)
  return net.FileListener(f)
}
