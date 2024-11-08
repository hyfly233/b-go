package main

import "fmt"

type Yarn struct {
	Name string
}

func NewYarn(name string) *Yarn {
	return &Yarn{Name: name}
}

func (y *Yarn) Initialize() {
	fmt.Printf("Initializing Yarn: %s\n", y.Name)
}
