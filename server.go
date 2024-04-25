package sylvain

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

var _ encoding.BinaryMarshaler = (*Server)(nil)
var _ encoding.BinaryUnmarshaler = (*Server)(nil)

type Server struct {
	leaseId string // lease id of the server

	ip   string
	port int

	Addr string `json:"addr"`
	Name string `json:"name,omitempty"` // name of the server
}

func (svr *Server) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, svr)
}

func (svr *Server) MarshalBinary() (data []byte, err error) {
	return json.Marshal(svr)
}

func (svr *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", svr.ip, svr.port)
}

func getHostIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(addr.String(), ":")[0]
	return ip
}
