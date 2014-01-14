// Package modbusclient provides modbus Serial Line/RTU and TCP/IP access
// for client (master) applications to communicate with server (slave)
// devices. This file specifies core definitions and data structures.

package modbusclient

const (
    MODBUS_PORT = 502

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

type TCPFrame struct {
    TransactionID int
    FunctionCode  byte
    Data          []byte
}
