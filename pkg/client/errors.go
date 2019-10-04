package client

import (
	"fmt"
)

type ReportingError struct {
	message string
}

func (r *ReportingError) Error() string {
	return r.message
}

func ErrorReadingPayload(err error) *ReportingError {
	return &ReportingError{message: fmt.Sprintf("Encountered error while reading payload: %s", err.Error())}
}

func ErrorConnecting(err error) *ReportingError {
	return &ReportingError{message: fmt.Sprintf("Encountered error while connecting to the grpc server: %s", err.Error())}
}

func ErrorSendingUsage(err error) *ReportingError {
	return &ReportingError{message: fmt.Sprintf("Encountered error while reporting usage: %s", err.Error())}
}

func IsReportingError(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(*ReportingError)
	return ok
}
