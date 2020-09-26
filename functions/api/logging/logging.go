package logging

import (
	"fmt"

	"github.com/google/uuid"
)

// OceanLogging ...
type OceanLogging interface {
	Info(message ...string)
}

// OceanLoggingImpl ...
type OceanLoggingImpl struct {
	requestID string
}

// NewOceanLoggingImpl ...
func NewOceanLoggingImpl() *OceanLoggingImpl {
	return &OceanLoggingImpl{
		requestID: uuid.New().String(),
	}
}

// Info ...
func (o *OceanLoggingImpl) Info(message ...string) {
	fmt.Printf("[INFO] [%s] %s\n", o.requestID, message)
}
