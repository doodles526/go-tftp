package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/doodles526/go-tftp/errors"
	"io"
	_ "net"
)

const (
	ReadRequestPacketOpcode  = 1
	WriteRequestPacketOpcode = 2
	DataPacketOpcode         = 3
	AckPacketOpcode          = 4
	ErrorPacketOpcode        = 5
)

type Packet interface {
	Encode() ([]byte, error)
}

type WriteRequestPacket struct {
	Filename string
	Mode     string
}

type ReadRequestPacket struct {
	Filename string
	Mode     string
}

type DataPacket struct {
	BlockNumber uint16
	Data        []byte
}

type ErrorPacket struct {
	ErrorCode    uint16
	ErrorMessage string
}

type AckPacket struct {
	BlockNumber uint16
}

func (w *WriteRequestPacket) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	// BigEndian equivelant to Network Byte Order
	err := binary.Write(buffer, binary.BigEndian, uint16(WriteRequestPacketOpcode))
	if err != nil {
		return nil, err
	}
	l, err := buffer.WriteString(w.Filename)
	// docs say err is always nil, but doing this to stay idiomatic
	if err != nil {
		return nil, err
	}
	if l != len(w.Filename) {
		return nil, fmt.Errorf("Length of filename did not match that written to buffer")
	}

	if err = buffer.WriteByte(0x00); err != nil {
		return nil, err
	}

	l, err = buffer.WriteString(w.Mode)
	if err != nil {
		return nil, err
	}
	if l != len(w.Mode) {
		return nil, fmt.Errorf("Length of mode did not match that written to buffer")
	}

	if err = buffer.WriteByte(0x00); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (r *ReadRequestPacket) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	// BigEndian equivelant to Network Byte Order
	err := binary.Write(buffer, binary.BigEndian, uint16(ReadRequestPacketOpcode))
	if err != nil {
		return nil, err
	}
	l, err := buffer.WriteString(r.Filename)
	// docs say err is always nil, but doing this to stay idiomatic
	if err != nil {
		return nil, err
	}
	if l != len(r.Filename) {
		return nil, fmt.Errorf("Length of filename did not match that written to buffer")
	}

	if err = buffer.WriteByte(0x00); err != nil {
		return nil, err
	}

	l, err = buffer.WriteString(r.Mode)
	if err != nil {
		return nil, err
	}
	if l != len(r.Mode) {
		return nil, fmt.Errorf("Length of mode did not match that written to buffer")
	}

	if err = buffer.WriteByte(0x00); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (d *DataPacket) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	// BigEndian equivelant to Network Byte Order
	err := binary.Write(buffer, binary.BigEndian, uint16(DataPacketOpcode))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, d.BlockNumber)
	if err != nil {
		return nil, err
	}

	l, err := buffer.Write(d.Data)
	if l != len(d.Data) {
		return nil, fmt.Errorf("Length of data did not match that written to buffer")
	}
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (a *AckPacket) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	// BigEndian equivelant to Network Byte Order
	err := binary.Write(buffer, binary.BigEndian, uint16(AckPacketOpcode))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, a.BlockNumber)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (e *ErrorPacket) Encode() ([]byte, error) {
	if e.ErrorCode > 7 || e.ErrorCode < 0 {
		return nil, fmt.Errorf("Invalid Error Code")
	}

	buffer := new(bytes.Buffer)
	// BigEndian equivelant to Network Byte Order
	err := binary.Write(buffer, binary.BigEndian, uint16(ErrorPacketOpcode))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, e.ErrorCode)
	if err != nil {
		return nil, err
	}

	l, err := buffer.WriteString(e.ErrorMessage)
	if l != len(e.ErrorMessage) {
		return nil, fmt.Errorf("Length of error message does not match that written to buffer")
	}
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(0x00)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode will decode data received in a packet and return a Packet object
// Note: we are not defining Decode for the Packet interface
// because it is not useful for that to be exported
func Decode(packetByte []byte) (Packet, error) {
	if len(packetByte) < 2 {
		return nil, fmt.Errorf("No data in packet")
	}

	opcode := binary.BigEndian.Uint16(packetByte)
	switch opcode {
	case ReadRequestPacketOpcode:
		return decodeReadRequestPacket(packetByte)
	case WriteRequestPacketOpcode:
		return decodeWriteRequestPacket(packetByte)
	case DataPacketOpcode:
		return decodeDataPacket(packetByte)
	case AckPacketOpcode:
		return decodeAckPacket(packetByte)
	case ErrorPacketOpcode:
		return decodeErrorPacket(packetByte)
	default:
		return nil, errors.ErrorIllegalOperation{
			Message: fmt.Sprintf("An illegal operation was attempted: Unknown Opcode - %d", opcode),
		}
	}
}

func decodeErrorPacket(errorData []byte) (*ErrorPacket, error) {
	if len(errorData) < 5 {
		return nil, errors.ErrorIllegalOperation{
			Message: "Invalid Error packet length - must be 5 bytes",
		}
	}

	errorCode := binary.BigEndian.Uint16(errorData[2:])

	if errorCode > 8 || errorCode < 0 {
		return nil, errors.ErrorIllegalOperation{
			Message: "Invalid error code - must be between 0 and 7",
		}
	}

	buffer := bytes.NewBuffer(errorData)
	errorMsg, err := buffer.ReadString(0x00)
	if err != nil {
		switch err {
		case io.EOF:
			return nil, errors.ErrorIllegalOperation{
				Message: "Error message not 0x0 terminated",
			}
		default:
			return nil, err
		}
	}

	return &ErrorPacket{
		ErrorCode:    errorCode,
		ErrorMessage: errorMsg,
	}, nil
}

func decodeAckPacket(ackData []byte) (*AckPacket, error) {
	if len(ackData) != 4 {
		return nil, errors.ErrorIllegalOperation{
			Message: "Invalid ACK packet length - must be 4 bytes",
		}
	}

	blockNumber := binary.BigEndian.Uint16(ackData[2:])
	return &AckPacket{
		BlockNumber: blockNumber,
	}, nil
}

func decodeDataPacket(dataByte []byte) (*DataPacket, error) {
	if len(dataByte) < 4 {
		return nil, errors.ErrorIllegalOperation{
			Message: "Data packet too short",
		}
	}

	blockNumber := binary.BigEndian.Uint16(dataByte[2:])
	return &DataPacket{
		BlockNumber: blockNumber,
		Data:        dataByte[4:],
	}, nil
}

func decodeWriteRequestPacket(readData []byte) (*WriteRequestPacket, error) {
	filename, mode, err := decodeRequest(readData)
	if err != nil {
		return nil, err
	}

	return &WriteRequestPacket{
		Filename: filename,
		Mode:     mode,
	}, nil
}

func decodeReadRequestPacket(readData []byte) (*ReadRequestPacket, error) {
	filename, mode, err := decodeRequest(readData)
	if err != nil {
		return nil, err
	}

	return &ReadRequestPacket{
		Filename: filename,
		Mode:     mode,
	}, nil
}

// decodeRequest exists simply to reduce code duplication in
// decodeXXXXRequestPacket functions
func decodeRequest(reqData []byte) (string, string, error) {
	// 2 byte opcode + non-empty string filename(1+ bytes) +  1 byte stop
	// + non-empty string filename(1+ bytes) + 1 byte stop
	if len(reqData) < 6 {
		return "", "", errors.ErrorIllegalOperation{
			Message: "RRQ not long enough",
		}
	}

	buffer := bytes.NewBuffer(reqData[2:])

	filename, err := buffer.ReadString(0x00)
	if err != nil {
		switch err {
		case io.EOF:
			return "", "", errors.ErrorIllegalOperation{
				Message: "Non 0x0 terminated Filename",
			}
		default:
			return "", "", err
		}
	}

	// 2 since the ReadString includes the termination
	if len(filename) < 2 {
		return "", "", errors.ErrorIllegalOperation{
			Message: "Blank Filename",
		}
	}

	// trimming the termination byte
	filename = filename[:len(filename)-1]

	mode, err := buffer.ReadString(0x00)
	if err != nil {
		switch err {
		case io.EOF:
			return "", "", errors.ErrorIllegalOperation{
				Message: "Non 0x0 terminated Mode",
			}
		default:
			return "", "", err
		}
	}

	// 2 since the ReadString includes the termination
	if len(mode) < 2 {
		return "", "", errors.ErrorIllegalOperation{
			Message: "Blank Mode",
		}
	}

	// trimming the termination byte
	mode = mode[:len(mode)-1]

	return filename, mode, nil
}

// ErrorToPacket takes an arbitrary error and converts it to an
// ErrorPacket
func ErrorToPacket(err error) *ErrorPacket {
	switch err := err.(type) {
	case errors.ErrorFileNotFound:
		return &ErrorPacket{
			ErrorCode:    1,
			ErrorMessage: err.Error(),
		}
	case errors.ErrorAccessViolation:
		return &ErrorPacket{
			ErrorCode:    2,
			ErrorMessage: err.Error(),
		}
	case errors.ErrorDiskFull:
		return &ErrorPacket{
			ErrorCode:    3,
			ErrorMessage: err.Error(),
		}
	case errors.ErrorIllegalOperation:
		return &ErrorPacket{
			ErrorCode:    4,
			ErrorMessage: err.Error(),
		}
	case errors.ErrorUnknownTransferID:
		return &ErrorPacket{
			ErrorCode:    5,
			ErrorMessage: err.Error(),
		}
	case errors.ErrorFileExists:
		return &ErrorPacket{
			ErrorCode:    6,
			ErrorMessage: err.Error(),
		}
	case errors.ErrorNoSuchUser:
		return &ErrorPacket{
			ErrorCode:    7,
			ErrorMessage: err.Error(),
		}
	default:
		return &ErrorPacket{
			ErrorCode:    0,
			ErrorMessage: err.Error(),
		}
	}

}
