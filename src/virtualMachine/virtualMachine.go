package virtualmachine

import (
	"fmt"
	"os"
)

const modValue = 32768

// VirtualMachine can execute instructions per
// architecutre described in `arch-spec`
type VirtualMachine struct {
	memory []int
	registers [8]int
	stack []int
}

// NewVirtualMachine creates a new VirtualMachine
func NewVirtualMachine() *VirtualMachine {
	vm := VirtualMachine{
		memory: []int{},
		registers: [8]int{},
		stack: []int{},
	}
	return &vm
}

// LoadProgram loads a program into memory from the supplied bin file
func (vm *VirtualMachine) LoadProgram(file *os.File) error {
	stats, statErr := file.Stat()
	if statErr != nil {
		return statErr
	}
	stats.Size()
	input := make([]byte, stats.Size())
	program := []int{}
	file.Read(input)
	for i := 0; i < len(input); i += 2 {
		little, big := int(input[i]), int(input[i+1])
		program = append(program, big << 8 + little)
	}
	vm.memory = program
	return nil
}

func (vm *VirtualMachine) printMemoryHead() {
	fmt.Println(vm.memory[:11])
}

// Execute runs the program in memory
func (vm *VirtualMachine) Execute() error {
	for inst := 0; inst < len(vm.memory); {
		switch vm.memory[inst] {
		case 0:
			return nil
		case 1:
			
		}
	}

	return nil
}

func (vm *VirtualMachine) set(a int, b int) {

}

func (vm *VirtualMachine) push(a int) {

}

func (vm *VirtualMachine) pop(a int) {

}

func (vm *VirtualMachine) eq(a int, b int, c int) {

}

func (vm *VirtualMachine) gt(a int, b int, c int) {

}

func (vm *VirtualMachine) jmp(a int) {

}

func (vm *VirtualMachine) jt(a int, b int) {

}

func (vm *VirtualMachine) jf(a int, b int) {
	
}

func (vm *VirtualMachine) add(a int, b int, c int) {

}

func (vm *VirtualMachine) mult(a int, b int, c int) {
	
}

func (vm *VirtualMachine) mod(a int, b int, c int) {

}

func (vm *VirtualMachine) and(a int, b int, c int) {

}

func (vm *VirtualMachine) or(a int, b int, c int) {

}

func (vm *VirtualMachine) not(a int, b int) {

}

func (vm *VirtualMachine) rmem(a int, b int) {

}

func (vm *VirtualMachine) wmem(a int, b int) {

}

func (vm *VirtualMachine) call(a int) {

}

func (vm *VirtualMachine) ret() {

}

func (vm *VirtualMachine) out(a int) {

}

func (vm *VirtualMachine) int(a int) {

}

func (vm *VirtualMachine) noop() {
	
}

