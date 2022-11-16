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
	c := &connection{
		id:   id,
		conn: conn,
	}

	return c
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

func (c *connection) receiver(s *Server) {
	for {
		data, _, err := wsutil.ReadClientData(c.conn)
		if err != nil {
			return
		}

		xy := struct {
			X  int    `json:"x"`
			Y  int    `json:"y"`
			Id string `json:"id"`
		}{}
		err = json.Unmarshal(data, &xy)
		if err != nil {
			return
		}

		s.broadcast(morphData{
			Type: "morph_data",
			Id:   "cursor_" + xy.Id,
			Html: fmt.Sprintf("<div id=\"cursor_%s\" class=\"cursor\" style=\"--x: %d; --y: %d;\">%[1]s</div>", xy.Id, xy.X, xy.Y),
		})
	}
}
