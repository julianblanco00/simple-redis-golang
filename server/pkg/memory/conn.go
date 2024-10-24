package memory

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

const CONN_ID_LENGTH = 36

func handleReadFromConn(conn net.Conn, data *Data) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("connection closed")
				conn.Close()
				break
			}
		}

		connId := buf[:CONN_ID_LENGTH]
		cmd := buf[CONN_ID_LENGTH:n]

		result, error := parseCommand(strings.TrimSpace(string(cmd)), data)
		if error != nil {
			conn.Write([]byte(error.Error()))
			continue
		}

		fmt.Fprint(conn, string(connId), result)
	}
}

func HandleConnection(listener net.Listener) {
	data := NewData()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection %v", err)
			continue
		}

		fmt.Printf("new tcp connection %s \n", conn.RemoteAddr())

		go handleReadFromConn(conn, data)
	}
}
