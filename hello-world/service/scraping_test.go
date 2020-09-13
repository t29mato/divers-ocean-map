package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ScrapingService ...
var ss ScrapingServiceImpl

func TestSetURL(t *testing.T) {
	url := "https://google.com"
	ss.setURL(url)
	assert.Equal(t, url, ss.url)
}
