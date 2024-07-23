package main

import (
	"log"
	"time"

	"github.com/luoyangwei/sylvain"
)

func main() {
	endpoints := sylvain.NewEtcdEndpoints([]string{"127.0.0.1:2379"}, 5*time.Second)
	serverDiscover := sylvain.NewNamedServerDiscover(endpoints, sylvain.WithServerNamed("customer"), sylvain.WithServerPort(9999))
	providerServer := serverDiscover.ElectionServerEndpoint("provider")
	log.Printf("discover: %s\n", providerServer)
	// ticker := time.NewTicker(3 * time.Second)
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		providerServer := serverDiscover.ElectionServerEndpoint("provider")
	// 		log.Printf("discover: %s\n", providerServer)
	// 	}
	// }
}
