package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/gobwas/ws/wsutil"
)

type connection struct {
	id   int
	conn net.Conn
}

func newConnection(id int, conn net.Conn) *connection {
	return &connection{
		id:   id,
		conn: conn,
	}
}

// Serialize the data to JSON and send it to the client
func (c *connection) send(m morphData) error {
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(c.conn, data)
	if err != nil {
		return err
	}

	return nil
}
