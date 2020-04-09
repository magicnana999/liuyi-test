package id

import (
	"github.com/kjk/betterguid"
)

func StrID() string{
	return betterguid.New()
}
