package main

import (
	"errors"
	"fmt"
)

func testError(){
	err1:= errors.New("1111")
	err2:= errors.New("2222")
	
	fmt.Errorf("err1 %v, err2 %v", err1,err2)
	fmt.Println(err1)
}