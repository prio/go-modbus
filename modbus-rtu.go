// Package modbusclient provides modbus Serial Line/RTU and TCP/IP access
// for client (master) applications to communicate with server (slave)
// devices. Logic specifically in this file implements the Serial Line/RTU
// protocol.

package modbusclient

import (
	"os"
	"syscall"
	"time"
)

// crc computes and returns a cyclic redundancy check of the given byte array
func crc(data []byte) uint16 {
	var crc16 uint16 = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		crc16 ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if crc16&0x0001 > 0 {
				crc16 = (crc16 >> 1) ^ 0xA001
			} else {
				crc16 >>= 1
			}
		}
	}
	return crc16
}

// GenerateRTUFrame is a method corresponding to a RTUFrame object which
// returns a byte array representing the associated serial line/RTU
// application data unit (ADU)
func (frame *RTUFrame) GenerateRTUFrame() []byte {

	packetLen := 8
	if len(frame.Data) > 0 {
		packetLen = RTU_FRAME_MAXSIZE
	}

	packet := make([]byte, packetLen)
	packet[0] = frame.SlaveAddress
	packet[1] = frame.FunctionCode
	packet[2] = byte(frame.StartRegister >> 8)       // (High Byte)
	packet[3] = byte(frame.StartRegister & 0xff)     // (Low Byte)
	packet[4] = byte(frame.NumberOfRegisters >> 8)   // (High Byte)
	packet[5] = byte(frame.NumberOfRegisters & 0xff) // (Low Byte)
	bytesUsed := 6

	for i := 0; i < len(frame.Data); i++ {
		packet[(bytesUsed + i)] = frame.Data[i]
	}
	bytesUsed += len(frame.Data)

	// add the crc to the end
	packet_crc := crc(packet[:bytesUsed])
	packet[bytesUsed] = byte(packet_crc & 0xff)
	packet[(bytesUsed + 1)] = byte(packet_crc >> 8)

	return packet[:bytesUsed]
}

// TransmitAndReceive is a method corresponding to an RTUFrame object
// which generates the corresponding ADU, transmits it to the modbus server
// (slave device) specified by the given file pointer (serial port), and
// returns a byte array of th e slave device's reply, and error (if any)
func (frame *RTUFrame) TransmitAndReceive(fd *os.File) ([]byte, error) {
	// generate the ADU from the RTU frame
	adu := frame.GenerateRTUFrame()

	// transmit the ADU to the slave device via the
	// serial port represented by the fd pointer
	_, err := fd.Write(adu)
	if err != nil {
		return []byte{}, err
	}

	// allow the slave device adequate time to respond
	time.Sleep(300 * time.Millisecond)

	// then attempt to read the reply
	response := make([]byte, TCP_FRAME_MAXSIZE)
	_, err = fd.Read(response)
	if err != nil {
		return []byte{}, err
	}

	// check the validity of the response
	if response[0] != frame.SlaveAddress || response[1] != frame.FunctionCode {
		if response[0] == frame.SlaveAddress && (response[1]&0x7f) == frame.FunctionCode {
			switch response[2] {
			case EXCEPTION_ILLEGAL_FUNCTION:
				return []byte{}, MODBUS_EXCEPTIONS[EXCEPTION_ILLEGAL_FUNCTION]
			case EXCEPTION_DATA_ADDRESS:
				return []byte{}, MODBUS_EXCEPTIONS[EXCEPTION_DATA_ADDRESS]
			case EXCEPTION_DATA_VALUE:
				return []byte{}, MODBUS_EXCEPTIONS[EXCEPTION_DATA_VALUE]
			case EXCEPTION_SLAVE_DEVICE_FAILURE:
				return []byte{}, MODBUS_EXCEPTIONS[EXCEPTION_SLAVE_DEVICE_FAILURE]
			}
		}
		return []byte{}, MODBUS_EXCEPTIONS[EXCEPTION_UNSPECIFIED]
	}

	// confirm the checksum (crc)
	response_length := response[2]
	response_crc := crc(response[:3+response_length])
	if response[(3+response_length)] != byte((response_crc&0xff)) ||
		response[(3+response_length+1)] != byte((response_crc>>8)) {
		// crc failed (odd that there's no specific code for it)
		return []byte{}, MODBUS_EXCEPTIONS[EXCEPTION_UNSPECIFIED]
	}
	return response[3:(response_length + 3)], nil
}

// viaRTU is a private method which applies the given function validator,
// to make sure the functionCode passed is valid for the operation
// desired. If correct, it creates an RTUFrame given the corresponding
// information, attempts to open the serialDevice, and if successful, calls
// TransmitAndReceive, returning the result. Otherwise, it returns an illegal
// function error, or the I/O device access error, whichever it encountered.
func viaRTU(fnValidator func(byte) bool, serialDevice string, slaveAddress, functionCode byte, startRegister, numRegisters uint16, data []byte) ([]byte, error) {
	if fnValidator(functionCode) {
		frame := new(RTUFrame)
		frame.SlaveAddress = slaveAddress
		frame.FunctionCode = functionCode
		frame.StartRegister = startRegister
		frame.NumberOfRegisters = numRegisters
		if len(data) > 0 {
			frame.Data = data
		}

		// make sure we can access the serial device
		fd, err := os.OpenFile(serialDevice, os.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY, 0666)
		if err != nil {
			return []byte{}, err
		} else {
			defer fd.Close()
			result, err := frame.TransmitAndReceive(fd)
			return result, err
		}
	}
	return []byte{}, MODBUS_EXCEPTIONS[EXCEPTION_ILLEGAL_FUNCTION]
}

// RTURead performs the given modbus Read function over RTU to the given
// serialDevice, using the given frame data
func RTURead(serialDevice string, slaveAddress, functionCode byte, startRegister, numRegisters uint16) ([]byte, error) {
	return viaRTU(ValidReadFunction, serialDevice, slaveAddress, functionCode, startRegister, numRegisters, []byte{})
}

// RTUWrite performs the given modbus Write function over RTU to the given
// serialDevice, using the given frame data
func RTUWrite(serialDevice string, slaveAddress, functionCode byte, startRegister, numRegisters uint16, data []byte) ([]byte, error) {
	return viaRTU(ValidWriteFunction, serialDevice, slaveAddress, functionCode, startRegister, numRegisters, data)
}
