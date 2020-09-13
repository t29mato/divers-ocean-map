package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
)

type GoqueryMock struct {
	mock.Mock
}

func TestScrapeIOP(t *testing.T) {
	main()
	goqueryMock := new(GoqueryMock)
	gp := filepath.Join("testdata", t.Name()+".golden")
	fmt.Println(gp)
	goqueryMock.On("NewDocument", "https://iop-dc.com/").Return()

}
