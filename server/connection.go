package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/gobwas/ws/wsutil"
)

type connection struct {
	conn net.Conn
}

func newConnection(conn net.Conn) *connection {
	return &connection{
		conn: conn,
	}
}

// Serialize the data to JSON and send it to the client
func (c *connection) send(m morphData) {
	if c.conn == nil {
		return
	}
	data, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}

	err = wsutil.WriteServerText(c.conn, data)
	if err != nil {
		fmt.Println(err)
	}
}
