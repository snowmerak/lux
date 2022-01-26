package context

import (
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type WSContext struct {
	Conn net.Conn
}

func (w *WSContext) Close() error {
	return w.Conn.Close()
}

func (w *WSContext) ReadBinary() ([]byte, error) {
	return wsutil.ReadClientBinary(w.Conn)
}

func (w *WSContext) ReadText() ([]byte, error) {
	return wsutil.ReadClientText(w.Conn)
}

func (w *WSContext) ReadData() ([]byte, ws.OpCode, error) {
	return wsutil.ReadClientData(w.Conn)
}

func (w *WSContext) WriteBinary(bin []byte) error {
	return wsutil.WriteServerBinary(w.Conn, bin)
}

func (w *WSContext) WriteText(text []byte) error {
	return wsutil.WriteServerText(w.Conn, text)
}

func (w *WSContext) WriteData(data []byte, opCode ws.OpCode) error {
	return wsutil.WriteServerMessage(w.Conn, opCode, data)
}
