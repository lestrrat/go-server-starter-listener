go-server-starter-listener
==========================

Creates a net.Listener that works with start_server (Server::Starter) utility

```go
import (
    "net"
    "net/http"
    "os"
    "github.com/lestrrat/go-server-starter-listener"
)

func main () {
    pwd, _ := os.Getwd()

    l, _ := ss.NewListener()
    if l == nil {
        // Fallback if not running under Server::Starter
        l, err := net.Listen("tcp", ":8080")
        if err != nil {
            panic("Failed to listen to port 8080")
        }
    }

    http.Serve(l, http.FileServer(http.Dir(pwd))
}
```
