package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Conn struct {
	conn          net.Conn
	controlReader *bufio.Reader
	controlWriter *bufio.Writer
	logger        Logger
	server        *Server
	sessionID     string
	closed        bool
}

func (conn *Conn) Serve() {
	conn.logger.Print(conn.sessionID, "Connection Established")
	for {
		//line, err := conn.controlReader.ReadString('\n')
		dt := make([]byte, 72)

		_, err := io.ReadFull(conn.controlReader, dt)
		if err != nil {
			if err != io.EOF {
				conn.logger.Print(conn.sessionID, fmt.Sprint("read error:", err))
			}

			break
		}
		conn.receiveLine(dt)
		// QUIT command closes connection, break to avoid error on reading from
		// closed socket
		if conn.closed == true {
			break
		}
	}
	conn.Close()
	conn.logger.Print(conn.sessionID, "Connection Terminated")
}
func (conn *Conn) parseLine(line []byte) ([]byte, []byte) {
	//4位终端号 +1位命令+ 1 位长度 + 64 位数据 + 2 位校验 = 72位
	ter := line[:4]
	dat := line[6:70]
	return ter, dat
}

func (conn *Conn) receiveLine(line []byte) {
	command, param := conn.parseLine(line)
	conn.logger.PrintCommand(conn.sessionID, fmt.Sprintf("%x", command), fmt.Sprintf("%x", param))
}

func (conn *Conn) writeMessage(code int, message []byte) (wrote int, err error) {
	conn.logger.PrintResponse(conn.sessionID, code, fmt.Sprintf("%x", message))
	//line := fmt.Sprintf("%d %s\r\n", code, message)
	//wrote, err = conn.controlWriter.WriteString(line)
	conn.controlWriter.Write(message)
	conn.controlWriter.Flush()
	return
}

func (conn *Conn) Close() {
	conn.conn.Close()
	conn.closed = true
}
