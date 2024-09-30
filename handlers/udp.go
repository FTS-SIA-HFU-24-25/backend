package handlers

import (
	"fmt"
	"sia/backend/db"
	"sia/backend/models"
	"sia/backend/tools"

	"github.com/xtaci/kcp-go/v5"
)

type (
	ConnectionChunk struct {
		Chunk      []byte
		Length     int
		CurrLength int
		IsStream   bool
	}
)

var (
	MAX_ERR_THRESHOLD                               = 5
	current_conn      map[[8]byte]models.Connection = make(map[[8]byte]models.Connection)
)

func HandleUDPConn(session *kcp.UDPSession) {
	err_threshold := 0
	buf := make([]byte, 256)
	for {
		if err_threshold >= MAX_ERR_THRESHOLD {
			tools.Log("[UDP Handler]", "Error threshold reached. Closing connection.")
			session.Close()
			break
		}
		n, err := session.Read(buf)
		if err != nil {
			tools.Log("[UDP Handler]", err)
			err_threshold++
			continue
		}

		fmt.Println("Received ", buf[:n])
		var uuid [8]byte
		var conn models.Connection
		raw_uuid := buf[:8]
		copy(uuid[:], raw_uuid[:])
		conn, ok := current_conn[uuid]
		if !ok {
			new_conn, err := db.RedisDB.GetConnection(uuid)
			if err != tools.OK {
				tools.Log("[UDP Handler]", err)
				err_threshold++
				continue
			}
			current_conn[uuid] = new_conn
			conn = new_conn
		}

		ver := buf[8]
		if int(ver) != tools.SERVICE_VERSION {
			tools.Log("[UDP Handler]", "Version mismatch")
			err_threshold++
			continue
		}

		data_length := int(buf[9])
		if data_length != n-3 {
			tools.Log("[UDP Handler]", "Lenght missmatch")
		}

		datas := buf[9:]
		conn.Data = make([]byte, n-3)
		conn.Data = append(conn.Data, datas...)

		n, err = session.Write(buf[:n])
		if err != nil {
			tools.Log("[UDP Handler]", err)
			err_threshold++
			continue
		}
	}
}
