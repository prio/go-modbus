// Package modbusclient provides modbus Serial Line/RTU and TCP/IP access
// for client (master) applications to communicate with server (slave)
// devices. Logic specifically in this file implements the TCP/IP protocol.

package modbusclient

import (
    "fmt"
    "io/ioutil"
    "net"
)

// GenerateTCPFrame is a method corresponding to a TCPFrame object which
// returns a byte array representing the associated TCP/IP application data
// unit (ADU)
func (frame *TCPFrame) GenerateTCPFrame() []byte {
    packetLen := len(frame.Data) + 8             // 7 bytes for the header + 1 for the function code
    packet := make([]byte, packetLen)            //len(m.Data)+8) // 7 bytes for the header + 1 for the function code
    packet[0] = byte(frame.TransactionID >> 8)   // Transaction ID (High Byte)
    packet[1] = byte(frame.TransactionID & 0xff) //                (Low Byte)
    packet[2] = 0x00                             // Protocol ID (2 bytes) -- always 00
    packet[3] = 0x00
    packet[4] = byte((packetLen - 4) >> 8)   // Remaining length of packet (High Byte)
    packet[5] = byte((packetLen - 4) & 0xff) //                            (Low Byte)
    packet[6] = 0x01                         // Unit ID (1 byte)
    packet[7] = frame.FunctionCode
    packet = append(packet, frame.Data...)
    fmt.Println(fmt.Sprintf("%x", packet))

    return packet
}

// TransmitAndReceive is a method corresponding to a TCPFrame object which
// generates the corresponding ADU, transmits it to the modbus server
// (slave device) specified by the TCP address, and returns a byte array
// of the slave device's reply, and error (if any)
func (frame *TCPFrame) TransmitAndReceive(slaveAddress string) ([]byte, error) {
    adu := frame.GenerateTCPFrame()
    addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", slaveAddress, MODBUS_PORT))
    if err == nil {
        conn, err := net.DialTCP("tcp", nil, addr)
        defer conn.Close()
        if err == nil {
            _, err = conn.Write(adu)
            if err == nil {
                res, err := ioutil.ReadAll(conn)
                if err == nil {
                    return res, nil
                }
            }
        }
    }
    return []byte{}, err
}
