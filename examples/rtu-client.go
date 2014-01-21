package main

import (
	"flag"
	"fmt"
	"github.com/dpapathanasiou/go-modbus"
	"log"
	"os"
)

func RTURead(fd *os.File, slaveAddress, functionCode byte, startRegister, numRegisters uint16) ([]byte, error) {
	if modbusclient.ValidReadFunction(functionCode) {
		frame := new(modbusclient.RTUFrame)
		frame.SlaveAddress = slaveAddress
		frame.FunctionCode = functionCode
		frame.StartRegister = startRegister
		frame.NumberOfRegisters = numRegisters
		result, err := frame.TransmitAndReceive(fd)
		return result, err
	}
	return []byte{}, modbusclient.MODBUS_EXCEPTIONS[modbusclient.EXCEPTION_ILLEGAL_FUNCTION]
}

func RTUWrite(fd *os.File, slaveAddress, functionCode byte, startRegister, numRegisters uint16, data []byte) ([]byte, error) {
	if modbusclient.ValidWriteFunction(functionCode) {
		frame := new(modbusclient.RTUFrame)
		frame.SlaveAddress = slaveAddress
		frame.FunctionCode = functionCode
		frame.StartRegister = startRegister
		frame.NumberOfRegisters = numRegisters
		frame.Data = data
		result, err := frame.TransmitAndReceive(fd)
		return result, err
	}
	return []byte{}, modbusclient.MODBUS_EXCEPTIONS[modbusclient.EXCEPTION_ILLEGAL_FUNCTION]
}

func main() {

	// get the serial port from the command line
	var serialDevice string
	const defaultInput = ""
	flag.StringVar(&serialDevice, "serial", defaultInput, "Serial port (RS485) to use, e.g., /dev/ttyS0 (try \"dmesg | grep tty\" to find)")
	flag.Parse()

	if len(serialDevice) > 0 {
		// make sure we can access the serial device
		fd, err := os.Open(serialDevice)
		if err != nil {
			log.Println(fmt.Sprintf("Unable to open serial port: %s", serialDevice))
		} else {
			// attempt to read the third register (0x03) on slave device number 2 (0x02)
			// and log the result
			readResult, readErr := RTURead(fd, 0x02, modbusclient.FUNCTION_READ_COILS, 0x03, 0x01)
			if readErr != nil {
				log.Println(readErr)
			}
			log.Println(readResult)

			// attempt to write a single byte to the third register (0x03) on slave
			// device number 2 (0x02) and log the result
			writeResult, writeErr := RTUWrite(fd, 0x02, modbusclient.FUNCTION_WRITE_SINGLE_REGISTER, 0x03, 0x01, []byte{0, 1})
			if writeErr != nil {
				log.Println(writeErr)
			}
			log.Println(writeResult)
		}
	} else {
		// display the command line usage requirements
		flag.PrintDefaults()
	}

}
