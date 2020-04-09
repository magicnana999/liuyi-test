package utils

import (
	"strings"
)

func StringIsBlank(src string) bool{
	return src == "" || len(strings.TrimSpace(src))==0
}