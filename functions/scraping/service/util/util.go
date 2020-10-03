package util

import (
	"strconv"

	"golang.org/x/text/width"
)

// ConvertIntFromFullWidthString
func ConvertIntFromFullWidthString(s *string) (int, error) {
	return strconv.Atoi(width.Narrow.String(*s))
}
