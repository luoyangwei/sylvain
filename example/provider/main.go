package main

import (
	"time"

	"github.com/luoyangwei/sylvain"
)

func main() {
	endpoints := sylvain.NewEtcdEndpoints([]string{"127.0.0.1:2379"}, 5*time.Second)
	_ = sylvain.NewNamedServerDiscover(endpoints, sylvain.WithServerNamed("provider"), sylvain.WithServerPort(8888))

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// log.Println("running...")
		}
	}
}
