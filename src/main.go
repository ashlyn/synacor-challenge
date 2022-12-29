package main

import (
	"os"

	virtualmachine "github.com/ashlyn/synacor-challenge/src/virtualMachine"
)

func main() {
	f, err := os.Open("challenge.bin")
	if err != nil {
		panic(err)
	}

	vm := virtualmachine.NewVirtualMachine()
	vm.LoadProgram(f)
	vm.Execute()
	
	defer f.Close()
}