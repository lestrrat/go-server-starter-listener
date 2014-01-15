package ss

import(
  "os"
  "testing"
)

func TestPort(t *testing.T) {
  os.Setenv("SERVER_STARTER_PORT", "8080=4")
  ports, err := Ports()
  if err != nil {
    t.Errorf("Failed to parse ports from env: %s", err)
  }

  if len(ports) <= 0 {
    t.Errorf("no ports found?!")
  }

  if ports[0].Fd != 4 {
    t.Errorf("fd is not what we expected")
  }
}