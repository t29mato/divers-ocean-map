package logging

import "fmt"

// OceanLogging ...
type OceanLogging interface {
	Info(message []string)
}

// OceanLoggingImpl ...
type OceanLoggingImpl struct {
	requestID string
}

// NewOceanLoggingImpl ...
func NewOceanLoggingImpl(requestID string) *OceanLoggingImpl {
	return &OceanLoggingImpl{
		requestID: requestID,
	}
}

// Info ...
func (o *OceanLoggingImpl) Info(message []string) {
	fmt.Println(message)
}
