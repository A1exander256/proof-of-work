package tcputils

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func ReadMessage(conn net.Conn) ([]byte, error) {
	var (
		lng uint64
		res []byte
	)

	if err := binary.Read(conn, binary.BigEndian, &lng); err != nil {
		return nil, fmt.Errorf("binary reading message: %w", err)
	}

	if _, err := io.ReadFull(conn, res); err != nil {
		return nil, fmt.Errorf("reading message: %w", err)
	}

	return res, nil
}

func WriteMessage(conn net.Conn, msg []byte) error {
	if err := binary.Write(conn, binary.BigEndian, uint64(len(msg))); err != nil {
		return fmt.Errorf("binary writing message: %w", err)
	}

	if _, err := conn.Write(msg); err != nil {
		return fmt.Errorf("writing message: %w", err)
	}

	return nil
}
