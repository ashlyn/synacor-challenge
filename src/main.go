package main

import (
	"bufio"
	"os"

	virtualmachine "github.com/ashlyn/synacor-challenge/src/virtualMachine"
)

func main() {
	f, err := os.Open("challenge.bin")
	if err != nil {
		panic(err)
	}

  reader := bufio.NewReader(os.Stdin)
	vm := virtualmachine.NewVirtualMachine(reader)
	// vm.LoadTestProgram("9,32768,32769,43,19,32768")
	vm.LoadProgram(f)
	vm.Execute(false)
	
	defer f.Close()
}