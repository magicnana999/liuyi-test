package utils

import (
	"errors"
	"fmt"
)

func PanicError(err error){
	if err != nil {
		panic(err)
	}
}

func NewError1(msg string) error{
	return errors.New(msg)
}

func NewError2(src string,msg string)error{
	return errors.New(fmt.Sprintf(src,msg))
}


