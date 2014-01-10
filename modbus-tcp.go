// Package modbusclient provides modbus Serial Line/RTU and TCP/IP access
// for client (master) applications to communicate with server (slave)
// devices. Logic specifically in this file implements the TCP/IP protocol.

package modbusclient

import (
    "io/ioutil"
    "net"
)

func send(a string, d []byte) ([]byte, error) {
    addr, err := net.ResolveTCPAddr("tcp4", a)
    if err == nil {
        c, err := net.DialTCP("tcp", nil, addr)
        if err == nil {
            _, err = c.Write(d)
            if err == nil {
                r, err := ioutil.ReadAll(c)
                if err == nil {
                    return r, nil
                }
            }
        }
    }
    return []byte{}, err
}

type MbTcp struct {
    Addr byte
    Code byte
    Data []byte
}

func (m MbTcp) generate() []byte {
    head := make([]byte, 8, 8)
    l := byte(len(m.Data) + 2)
    head[0] = 0x00
    head[1] = 0x00
    head[2] = 0x00
    head[3] = 0x00
    head[4] = 0x00
    head[5] = l
    head[6] = m.Addr
    head[7] = m.Code
    body := make([]byte, 260)
    body = append(body, head...)
    body = append(body, m.Data...)
    return body
}

func (m *MbTcp) TCPSend(addr string) ([]byte, error) {
    req := m.generate()
    return send(addr, req)
}
