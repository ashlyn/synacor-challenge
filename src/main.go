package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	virtualmachine "github.com/ashlyn/synacor-challenge/src/virtualMachine"
)

func main() {
	vm := createAutomatedGame()
	// vm := createManualGame()
	err := vm.Execute(false)
	if err != nil {
		panic(err)
	}
}

func createAutomatedGame() *virtualmachine.VirtualMachine {
	f, err := os.Open("challenge.bin")
	if err != nil {
		panic(err)
	}

	logFileName := fmt.Sprintf("notes/logs/%v.log", time.Now().Format("2006-01-02 15:04:05"))
	logFile, logFileErr := os.Create(logFileName)
	if logFileErr != nil {
		panic(err)
	}
	writer := bufio.NewWriter(logFile)

	automated, _ := os.Open("notes/Game.md")
  reader := bufio.NewReader(automated) // bufio.NewReader(os.Stdin)
	vm := virtualmachine.NewVirtualMachine(reader, writer)
	// vm.LoadTestProgram("9,32768,32769,43,19,32768")
	vm.LoadProgram(f)

	/*
	dump, dumpErr := os.Create("notes/memoryDump.log")
	dumpWriter := bufio.NewWriter(dump)
	dumpErr = vm.MemoryDump(dumpWriter)
	if dumpErr != nil {
		panic(dumpErr)
	}
	dump.Close()
	*/
	
	vm.SetEnergyLevel(1000)
	
	defer f.Close()
	defer logFile.Close()

	return vm
}

func createManualGame() *virtualmachine.VirtualMachine {
	f, err := os.Open("challenge.bin")
	if err != nil {
		panic(err)
	}

	logFileName := fmt.Sprintf("notes/logs/%v.log", time.Now().Format("2006-01-02 15:04:05"))
	logFile, logFileErr := os.Create(logFileName)
	if logFileErr != nil {
		panic(err)
	}
	writer := bufio.NewWriter(logFile)

  reader := bufio.NewReader(os.Stdin)
	vm := virtualmachine.NewVirtualMachine(reader, writer)
	// vm.LoadTestProgram("9,32768,32769,43,19,32768")
	vm.LoadProgram(f)
	
	defer f.Close()
	defer logFile.Close()

	return vm
}