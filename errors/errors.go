package errors

import (
	"fmt"
)

type ErrorFileNotFound struct {
	File string
}

func (e ErrorFileNotFound) Error() string {
	return fmt.Sprintf("Error File Not Found - %s", e.File)
}

type ErrorAccessViolation struct {
}

func (e ErrorAccessViolation) Error() string {
	return fmt.Sprintf("Error Access Violation")
}

type ErrorDiskFull struct {
}

func (e ErrorDiskFull) Error() string {
	return fmt.Sprintf("Error Disk Full")
}

type ErrorIllegalOperation struct {
	Message string
}

func (e ErrorIllegalOperation) Error() string {
	return fmt.Sprintf("Error Illegal Operation -  %s", e.Message)
}

type ErrorUnknownTransferID struct {
	TransferID string
}

func (e ErrorUnknownTransferID) Error() string {
	return fmt.Sprintf("Error Unknown Transfer ID - %s", e.TransferID)
}

type ErrorFileExists struct {
	File string
}

func (e ErrorFileExists) Error() string {
	return fmt.Sprintf("Error File Exists: %s", e.File)
}

type ErrorNoSuchUser struct {
	User string
}

func (e ErrorNoSuchUser) Error() string {
	return fmt.Sprintf("Error No Such User: %s", e.User)
}
