// Package modbusclient provides modbus Serial Line/RTU and TCP/IP access
// for client (master) applications to communicate with server (slave)
// devices. Logic specifically in this file implements the TCP/IP protocol.

package modbusclient

import (
	"fmt"
	"net"
)

// GenerateTCPFrame is a method corresponding to a TCPFrame object which
// returns a byte array representing the associated TCP/IP application data
// unit (ADU)
func (frame *TCPFrame) GenerateTCPFrame() []byte {
	packetLen := len(frame.Data) + 8 // 7 bytes for the header + 1 for the function code
	packet := make([]byte, packetLen)
	packet[0] = byte(frame.TransactionID >> 8)   // Transaction ID (High Byte)
	packet[1] = byte(frame.TransactionID & 0xff) //                (Low Byte)
	packet[2] = 0x00                             // Protocol ID (2 bytes) -- always 00
	packet[3] = 0x00
	packet[4] = byte((packetLen - 6) >> 8)   // Remaining length of packet, beyond this point (High Byte)
	packet[5] = byte((packetLen - 6) & 0xff) //                                               (Low Byte)

	/* Unit ID (1 byte):
	   If the slave device is using an Ethernet-to-serial bridge, set this to the
	   corresponding SlaveAddress. Otherwise, use 0x00 for "do not bridge".
	*/
	if frame.EthernetToSerialBridge {
		packet[6] = frame.SlaveAddress
	} else {
		packet[6] = 0x00
	}
	packet[7] = frame.FunctionCode
	packet = append(packet[:8], frame.Data...)

	return packet
}

// TransmitAndReceive is a method corresponding to a TCPFrame object which
// generates the corresponding ADU, transmits it to the modbus server
// (slave device) specified by the TCP address+port, and returns a byte array
// of the slave device's reply, and error (if any)
func (frame *TCPFrame) TransmitAndReceive(server string, port int) ([]byte, error) {
	adu := frame.GenerateTCPFrame() // generate the ADU
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", server, port))

	if err == nil {
		// attempt to connect to the slave device (server)
		conn, err := net.DialTCP("tcp", nil, addr)
		defer conn.Close()

		if err == nil {
			// transmit the ADU
			_, err = conn.Write(adu)

			if err == nil {
				// read the response
				response := make([]byte, TCP_FRAME_MAXSIZE)
				n, err := conn.Read(response)

				if err == nil {
					// return only the number of bytes read
					return response[:n], nil
				}
			}
		}
	}
	return []byte{}, err
}

// viaTCP is a private method which applies the given function validator, to
// make sure the functionCode passed is valid for the operation desired. If
// correct, it creates a TCPFrame given the corresponding information,
// calls TransmitAndReceive, returning the result. Otherwise, it returns
// an illegal function error.
func viaTCP(fnValidator func(byte) bool, h string, p, transactionID int, functionCode byte, serialBridge bool, slaveAddress byte, data []byte) ([]byte, error) {
	if fnValidator(functionCode) {
		frame := new(TCPFrame)
		frame.TransactionID = transactionID
		frame.FunctionCode = functionCode
		frame.EthernetToSerialBridge = serialBridge
		frame.SlaveAddress = slaveAddress
		frame.Data = data
		result, err := frame.TransmitAndReceive(h, p)
		return result, err
	}
	return []byte{}, MODBUS_EXCEPTIONS[EXCEPTION_ILLEGAL_FUNCTION]
}

// TCPRead performs the given modbus Read function over TCP to the given
// host/port combination, using the given frame data
func TCPRead(h string, p, transactionID int, functionCode byte, serialBridge bool, slaveAddress byte, data []byte) ([]byte, error) {
	return viaTCP(ValidReadFunction, h, p, transactionID, functionCode, serialBridge, slaveAddress, data)
}

// TCPWrite performs the given modbus Write function over TCP to the given
// host/port combination, using the given frame data
func TCPWrite(h string, p, transactionID int, functionCode byte, serialBridge bool, slaveAddress byte, data []byte) ([]byte, error) {
	return viaTCP(ValidWriteFunction, h, p, transactionID, functionCode, serialBridge, slaveAddress, data)
}
