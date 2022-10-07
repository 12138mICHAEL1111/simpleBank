package main

import "fmt"

type  Father struct{
	 *Son //匿名字段 父类可以直接访问子类所有方法
}

type Son struct{
	
}

func (s *Son) test(){
	fmt.Println("111")
}

func testStruct(){
	s := &Son{}
	f := Father{s}
	f.test()
}
