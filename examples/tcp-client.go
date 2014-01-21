package main

import (
	"flag"
	"fmt"
	"github.com/dpapathanasiou/go-modbus"
	"log"
)

func TCPRead(h string, p, transactionID int, functionCode byte, serialBridge bool, slaveAddress byte) ([]byte, error) {
	if modbusclient.ValidReadFunction(functionCode) {
		frame := new(modbusclient.TCPFrame)
		frame.TransactionID = transactionID
		frame.FunctionCode = functionCode
		frame.EthernetToSerialBridge = serialBridge
		frame.SlaveAddress = slaveAddress
		frame.Data = []byte{0, 1}
		result, err := frame.TransmitAndReceive(h, p)
		return result, err
	}
	return []byte{}, modbusclient.MODBUS_EXCEPTIONS[modbusclient.EXCEPTION_ILLEGAL_FUNCTION]
}

func TCPWrite(h string, p, transactionID int, functionCode byte, serialBridge bool, slaveAddress byte, data []byte) ([]byte, error) {
	if modbusclient.ValidWriteFunction(functionCode) {
		frame := new(modbusclient.TCPFrame)
		frame.TransactionID = transactionID
		frame.FunctionCode = functionCode
		frame.EthernetToSerialBridge = serialBridge
		frame.SlaveAddress = slaveAddress
		frame.Data = data
		result, err := frame.TransmitAndReceive(h, p)
		return result, err
	}
	return []byte{}, modbusclient.MODBUS_EXCEPTIONS[modbusclient.EXCEPTION_ILLEGAL_FUNCTION]
}

func main() {

	// get device host (url or ip address) and port from the command line
	var (
		host string
		port int
	)
	const (
		defaultHost = "127.0.0.1"
		defaultPort = modbusclient.MODBUS_PORT
	)

	flag.StringVar(&host, "host", defaultHost, "Slave device host (url or ip address)")
	flag.IntVar(&port, "port", defaultPort, fmt.Sprintf("Slave device port (the default is %d)", defaultPort))
	flag.Parse()

	// attempt to read the holding registers
	readResult, readErr := TCPRead(host, port, 1, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x00)
	if readErr != nil {
		log.Println(readErr)
	}
	log.Println(readResult)

	// attempt to write to a single coil
	writeResult, writeErr := TCPWrite(host, port, 1, modbusclient.FUNCTION_WRITE_SINGLE_COIL, false, 0x00, []byte{0, 1})
	if writeErr != nil {
		log.Println(writeErr)
	}
	log.Println(writeResult)

}
