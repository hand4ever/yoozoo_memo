package main

import "fmt"

func main() {
	person := Person{"lpan", 50}
	fmt.Printf("person<%s:%d>\n", person.name, person.age)
	person.sayHi()
	person.ModifyAge(100)
	person.sayHi()
}

type Person struct {
	name string
	age  int
}

func (p Person) sayHi() {
	fmt.Printf("SayHi -- This is %s, my age is %d\n", p.name, p.age)
}

func (p Person) ModifyAge(age int) {
	fmt.Printf("ModifyAge")
	p.age = age
}
