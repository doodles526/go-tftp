package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/doodles526/go-tftp/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeReadPacket(t *testing.T) {
	readRQ := ReadRequestPacket{
		Filename: "./testfile",
		Mode:     "octet",
	}

	// Expected Val:
	//    2 bytes     string       1 byte     string      1 byte
	//   --------------------------------------------------------
	//   |   1   |  "./testfile"  |   0  |    "octet"    |   0  |
	//   --------------------------------------------------------
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, uint16(1))
	assert.NoError(t, err, "There should be no error writing to buffer")

	l, err := buffer.WriteString("./testfile")
	assert.Len(t, "./testfile", l, "We should write the same number of bytes as in string")
	assert.NoError(t, err, "There should be no error writing to buffer")

	err = buffer.WriteByte(0x00)
	assert.NoError(t, err, "There should be no error writing to buffer")

	l, err = buffer.WriteString("octet")
	assert.Len(t, "octet", l, "We should write the same number of bytes as in string")
	assert.NoError(t, err, "There should be no error writing to buffer")

	err = buffer.WriteByte(0x00)
	assert.NoError(t, err, "There should be no error writing to buffer")

	expected := buffer.Bytes()

	actual, err := readRQ.Encode()
	assert.NoError(t, err, "There should be no error when encoding a read packet")

	assert.Equal(t, expected, actual)
}

func TestEncodeWritePacket(t *testing.T) {
	writeRQ := WriteRequestPacket{
		Filename: "./testfile",
		Mode:     "octet",
	}

	// Expected Val:
	//    2 bytes     string       1 byte     string      1 byte
	//   --------------------------------------------------------
	//   |   2   |  "./testfile"  |   0  |    "octet"    |   0  |
	//   --------------------------------------------------------
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, uint16(2))
	assert.NoError(t, err, "There should be no error writing to buffer")

	l, err := buffer.WriteString("./testfile")
	assert.Len(t, "./testfile", l, "We should write the same number of bytes as in string")
	assert.NoError(t, err, "There should be no error writing to buffer")

	err = buffer.WriteByte(0x00)
	assert.NoError(t, err, "There should be no error writing to buffer")

	l, err = buffer.WriteString("octet")
	assert.Len(t, "octet", l, "We should write the same number of bytes as in string")
	assert.NoError(t, err, "There should be no error writing to buffer")

	err = buffer.WriteByte(0x00)
	assert.NoError(t, err, "There should be no error writing to buffer")

	expected := buffer.Bytes()

	actual, err := writeRQ.Encode()
	assert.NoError(t, err, "There should be no error when encoding a write packet")

	assert.Equal(t, expected, actual)
}

func TestEncodeDataPacket(t *testing.T) {
	dataRQ := DataPacket{
		BlockNumber: 50,
		Data:        []byte{84, 101, 115, 116}, // []byte representation of "Test"
	}

	// Expected Val:
	//      2 bytes     2 bytes      n bytes
	//      --------------------------------------------
	//     |   3   |   50   |  [84, 101, 115, 116]     |
	//      --------------------------------------------
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, uint16(3))
	assert.NoError(t, err, "There should be no error writing to buffer")

	err = binary.Write(buffer, binary.BigEndian, uint16(50))
	assert.NoError(t, err, "There should be no error when writing block number")

	l, err := buffer.Write([]byte{84, 101, 115, 116})
	assert.Len(t, []byte{84, 101, 115, 116}, l, "We should write the same number of bytes")
	assert.NoError(t, err, "There should be no error writing to buffer")

	expected := buffer.Bytes()

	actual, err := dataRQ.Encode()
	assert.NoError(t, err, "There should be no error when encoding a write packet")

	assert.Equal(t, expected, actual)
}

func TestAckPacket(t *testing.T) {
	ackRQ := AckPacket{
		BlockNumber: 50,
	}

	// Expected Val:
	//      2 bytes  2 bytes
	//      -----------------
	//     |   4   |   50   |
	//      -----------------
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, uint16(4))
	assert.NoError(t, err, "There should be no error writing to buffer")

	err = binary.Write(buffer, binary.BigEndian, uint16(50))
	assert.NoError(t, err, "There should be no error when writing block number")

	expected := buffer.Bytes()

	actual, err := ackRQ.Encode()
	assert.NoError(t, err, "There should be no error when encoding a write packet")

	assert.Equal(t, expected, actual)
}

func TestErrorPacket(t *testing.T) {
	errorRQ := ErrorPacket{
		ErrorCode:    1,
		ErrorMessage: "Test",
	}

	// Expected Val:
	//      2 bytes  2 bytes  String   1 byte
	//      -----------------------------------
	//     |   5   |   1    |  "Test" |   0   |
	//      -----------------------------------
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, uint16(5))
	assert.NoError(t, err, "There should be no error writing to buffer")

	err = binary.Write(buffer, binary.BigEndian, uint16(1))
	assert.NoError(t, err, "There should be no error when writing block number")

	l, err := buffer.WriteString("Test")
	assert.Len(t, "Test", l, "We should write the same number of bytes as in string")
	assert.NoError(t, err, "There should be no error writing to buffer")

	err = buffer.WriteByte(0x00)
	assert.NoError(t, err, "There should be no error writing to buffer")

	expected := buffer.Bytes()

	actual, err := errorRQ.Encode()
	assert.NoError(t, err, "There should be no error when encoding a write packet")

	assert.Equal(t, expected, actual)
}

func TestErrorToPacket(t *testing.T) {

	// FileNotFound
	fileNotFound := errors.ErrorFileNotFound{
		File: "test",
	}

	fileNotFoundPacket := ErrorToPacket(fileNotFound)
	assert.EqualValues(t, 1, fileNotFoundPacket.ErrorCode, "FileNotFound packet code should be 1")
	assert.EqualError(t, fileNotFound, fileNotFoundPacket.ErrorMessage)

	// AccessViolation
	accessViolation := errors.ErrorAccessViolation{}
	accessViolationPacket := ErrorToPacket(accessViolation)

	assert.EqualValues(t, 2, accessViolationPacket.ErrorCode)
	assert.EqualError(t, fileNotFound, fileNotFoundPacket.ErrorMessage)

	// DiskFull
	diskFull := errors.ErrorDiskFull{}
	diskFullPacket := ErrorToPacket(diskFull)

	assert.EqualValues(t, 3, diskFullPacket.ErrorCode)
	assert.EqualError(t, diskFull, diskFullPacket.ErrorMessage)

	// IllegalOperation
	illegalOperation := errors.ErrorIllegalOperation{
		Message: "test",
	}
	illegalOperationPacket := ErrorToPacket(illegalOperation)

	assert.EqualValues(t, 4, illegalOperationPacket.ErrorCode)
	assert.EqualError(t, illegalOperation, illegalOperationPacket.ErrorMessage)

	// UnknownTransferID
	unknownID := errors.ErrorUnknownTransferID{
		TransferID: "hello",
	}
	unknownIDPacket := ErrorToPacket(unknownID)

	assert.EqualValues(t, 5, unknownIDPacket.ErrorCode)
	assert.EqualError(t, unknownID, unknownIDPacket.ErrorMessage)

	// FileExists
	fileExists := errors.ErrorFileExists{
		File: "./test",
	}
	fileExistsPacket := ErrorToPacket(fileExists)

	assert.EqualValues(t, 6, fileExistsPacket.ErrorCode)
	assert.EqualError(t, fileExists, fileExistsPacket.ErrorMessage)

	// NoSuchUser
	noSuchUser := errors.ErrorNoSuchUser{
		User: "user",
	}
	noSuchUserPacket := ErrorToPacket(noSuchUser)

	assert.EqualValues(t, 7, noSuchUserPacket.ErrorCode)
	assert.EqualError(t, noSuchUser, noSuchUserPacket.ErrorMessage)

	// Arbitrary Error
	arbitraryError := fmt.Errorf("Random Error")
	arbitraryErrorPacket := ErrorToPacket(arbitraryError)

	assert.EqualValues(t, 0, arbitraryErrorPacket.ErrorCode)
	assert.EqualError(t, arbitraryError, arbitraryErrorPacket.ErrorMessage)
}
