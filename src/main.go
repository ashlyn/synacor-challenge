package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	virtualmachine "github.com/ashlyn/synacor-challenge/src/virtualMachine"
)

var challengeBinSrc string = "bin/challenge.bin"
func main() {
	vm := createAutomatedGame(false)
	// vm := createManualGame(challengeBinSrc, false)

	// energylevel.GetMaxEnergyLevel(4, 1) = 25734 but it takes a while to brute-force
	vm.SetEnergyLevel(25734)
	err := vm.Execute()
	
	if err != nil {
		panic(err)
	}
}

// automated games are hard-coded to use the challenge bin because
// it's the only bin with steps automation
func createAutomatedGame(debug bool) *virtualmachine.VirtualMachine {
	f, err := os.Open(challengeBinSrc)
	if err != nil {
		panic(err)
	}

	var writer *bufio.Writer = nil
	if debug {
		writer = createLogWriter()
	}

	automated, _ := os.Open("notes/Game.md")
  reader := bufio.NewReader(automated)
	vm := virtualmachine.NewVirtualMachine(reader, writer, debug)
	// vm.LoadTestProgram("9,32768,32769,43,19,32768")
	vm.LoadProgram(f)

	/*
	dump, dumpErr := os.Create("notes/memoryDump.log")
	dumpWriter := bufio.NewWriter(dump)
	dumpErr = vm.MemoryDump(dumpWriter)
	// OR	dumpErr = vm.WriteAssembly(dumpWriter)
	if dumpErr != nil {
		panic(dumpErr)
	}
	dump.Close()
	*/
	
	defer f.Close()

	return vm
}

func createManualGame(programSrc string, debug bool) *virtualmachine.VirtualMachine {
	f, err := os.Open(programSrc)
	if err != nil {
		panic(err)
	}

	var writer *bufio.Writer = nil
	if debug {
		writer = createLogWriter()
	}

  reader := bufio.NewReader(os.Stdin)
	vm := virtualmachine.NewVirtualMachine(reader, writer, debug)
	// vm.LoadTestProgram("9,32768,32769,43,19,32768")
	vm.LoadProgram(f)
	
	defer f.Close()

	return vm
}

func createLogWriter() *bufio.Writer {
	logFileName := fmt.Sprintf("notes/logs/%v.log", time.Now().Format("2006-01-02 15:04:05"))
	logFile, logFileErr := os.Create(logFileName)
	if logFileErr != nil {
		panic(logFileErr)
	}
	writer := bufio.NewWriter(logFile)
	defer logFile.Close()
	return writer
}