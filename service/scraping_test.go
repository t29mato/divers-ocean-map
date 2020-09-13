package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ScrapingService ...
var ss ScrapingServiceImpl

func TestSetURL(t *testing.T) {
	ss.setURL("https://google.com")
	assert.Equal(t, "https://google.com", ss.url)
}
