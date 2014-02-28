package main

import (
	"flag"
	"github.com/dpapathanasiou/go-modbus"
	"log"
)

func main() {

	// get the device serial port from the command line
	var serialDevice string
	const defaultInput = ""
	flag.StringVar(&serialDevice, "serial", defaultInput, "Serial port (RS485) to use, e.g., /dev/ttyS0 (try \"dmesg | grep tty\" to find)")
	// IRL: use more command-line arguments to define which registers to read or write...
	flag.Parse()

	if len(serialDevice) > 0 {

		// attempt to read the 200 address register on
		// slave device number 1 (0x01) at serialDevice
		readResult, readErr := modbusclient.RTURead(serialDevice, 0x01, modbusclient.FUNCTION_READ_COILS, 200, 1)
		if readErr != nil {
			log.Println(readErr)
		}
		log.Println(readResult)

		// attempt to write to a single coil at address 300
		// on slave device number 1 (0x01) at serialDevice
		writeResult, writeErr := modbusclient.RTUWrite(serialDevice, 0x01, modbusclient.FUNCTION_WRITE_SINGLE_COIL, 300, 1, []byte{0xff})
		if writeErr != nil {
			log.Println(writeErr)
		}
		log.Println(writeResult)

	} else {

		// display the command line usage requirements
		flag.PrintDefaults()

	}

}
