package main
type Animal struct{

}

func (a *Animal) eat(){

}

func (a *Animal) sleep(){
	
}

type dogInterface interface{
	eat()
}

type dog struct {
	d dogInterface //任何类型 只要是拥有eat方法 就是 dogInterface类型
}

func New(d dogInterface) *dog{
	return &dog{d:d}
}

func testInterface(){
	a := Animal{}
	tom := New(&a)
	tom.d.eat()
}