// Package modbusclient provides modbus Serial Line/RTU and TCP/IP access
// for client (master) applications to communicate with server (slave)
// devices. This file specifies core definitions and data structures.

package modbusclient

import (
	"errors"
)

const (
	MODBUS_PORT       = 502
	RTU_FRAME_MAXSIZE = 512
	TCP_FRAME_MAXSIZE = 260

	FUNCTION_READ_COILS                    = 0x01
	FUNCTION_READ_DISCRETE_INPUTS          = 0x02
	FUNCTION_READ_HOLDING_REGISTERS        = 0x03
	FUNCTION_READ_INPUT_REGISTERS          = 0x04
	FUNCTION_WRITE_SINGLE_COIL             = 0x05
	FUNCTION_WRITE_SINGLE_REGISTER         = 0x06
	FUNCTION_WRITE_MULTIPLE_REGISTERS      = 0x10
	FUNCTION_MODBUS_ENCAPSULATED_INTERFACE = 0x2B

	EXCEPTION_ILLEGAL_FUNCTION                        = 0x01
	EXCEPTION_DATA_ADDRESS                            = 0x02
	EXCEPTION_DATA_VALUE                              = 0x03
	EXCEPTION_SLAVE_DEVICE_FAILURE                    = 0x04
	EXCEPTION_ACKNOWLEDGE                             = 0x05
	EXCEPTION_SLAVE_DEVICE_BUSY                       = 0x06
	EXCEPTION_MEMORY_PARITY_ERROR                     = 0x08
	EXCEPTION_GATEWAY_PATH_UNAVAILABLE                = 0x0A
	EXCEPTION_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND = 0x0B
)

var MODBUS_EXCEPTIONS = map[uint16]error{
	EXCEPTION_ILLEGAL_FUNCTION:                        errors.New("Modbus Error: Illegal Function (0x01)"),
	EXCEPTION_DATA_ADDRESS:                            errors.New("Modbus Error: Data Address (0x02)"),
	EXCEPTION_DATA_VALUE:                              errors.New("Modbus Error: Data Value (0x03)"),
	EXCEPTION_SLAVE_DEVICE_FAILURE:                    errors.New("Modbus Error: Slave Device Failure (0x04)"),
	EXCEPTION_ACKNOWLEDGE:                             errors.New("Modbus Error: Acknowledge (0x05)"),
	EXCEPTION_SLAVE_DEVICE_BUSY:                       errors.New("Modbus Error: Slave Device Busy (0x06)"),
	EXCEPTION_MEMORY_PARITY_ERROR:                     errors.New("Modbus Error: Memory Parity Error (0x08)"),
	EXCEPTION_GATEWAY_PATH_UNAVAILABLE:                errors.New("Modbus Error: Gateway Path Unavailable (0x0A)"),
	EXCEPTION_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND: errors.New("Modbus Error: Gateway Target Device Failed to Respond (0x0B)"),
}

type TCPFrame struct {
	TransactionID int
	FunctionCode  byte
	Data          []byte
}

type RTUFrame struct {
	SlaveAddress      byte
	FunctionCode      byte
	StartRegister     uint16
	NumberOfRegisters uint16
	Data              []byte
}
