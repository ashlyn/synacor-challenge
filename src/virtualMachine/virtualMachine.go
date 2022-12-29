package virtualmachine

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const modValue = 32768

// VirtualMachine can execute instructions per
// architecutre described in `arch-spec`
type VirtualMachine struct {
	currentInstruction int
	memory []int
	registers [8]int
	stack []int
	reader *bufio.Reader
}

// NewVirtualMachine creates a new VirtualMachine
func NewVirtualMachine(reader *bufio.Reader) *VirtualMachine {
	vm := VirtualMachine{
		currentInstruction: 0,
		memory: []int{},
		registers: [8]int{},
		stack: []int{},
		reader: reader,
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

// LoadTestProgram loads a program into memory from a comma-separated string
func (vm *VirtualMachine) LoadTestProgram(input string) {
	program := strings.Split(input, ",")
	for _, p := range program {
		value, _ := strconv.Atoi(p)
		vm.memory = append(vm.memory, value)
	}
}

// PrintMemoryHead is a debugging tool which prints the first 10 memory items
func (vm *VirtualMachine) PrintMemoryHead() {
	tail := int(math.Min(11, float64(len(vm.memory))))
	fmt.Println(vm.memory[:tail])
}

// Execute runs the program in memory
func (vm *VirtualMachine) Execute(debug bool) error {
	for vm.currentInstruction < len(vm.memory) {
		switch vm.memory[vm.currentInstruction] {
		case 0:
			return nil
		case 1:
			vm.set()
		case 2:
			vm.push()
		case 3:
			vm.pop()
		case 4:
			vm.eq()
		case 5:
			vm.gt()
		case 6:
			vm.jmp()
		case 7:
			vm.jt()
		case 8:
			vm.jf()
		case 9:
			vm.add()
		case 10:
			vm.mult()
		case 11:
			vm.mod()
		case 12:
			vm.and()
		case 13:
			vm.or()
		case 14:
			vm.not()
		case 15:
			vm.rmem()
		case 16:
			vm.wmem()
		case 17:
			vm.call()
		case 18:
			vm.ret()
		case 19:
			vm.out()
		case 20:
			vm.in()
		case 21:
			vm.noop()
		default:
			return fmt.Errorf("Unrecognized op code %d", vm.memory[vm.currentInstruction])
		}
	}

	if debug {
		fmt.Println(vm.memory)
		fmt.Println(vm.registers)
		fmt.Println(vm.stack)
	}

	return nil
}

func (vm *VirtualMachine) writeValueToAddress(value int, address int) {
	isRegister, regAddress := isRegisterAddress(address)
	if isRegister {
		vm.registers[regAddress] = value
	} else {
		vm.memory[address] = value
	}
}

func (vm *VirtualMachine) value(address int) int {
	isRegister, registerAddress := isRegisterAddress(address)
	if isRegister {
		return vm.registers[registerAddress]
	} else if address >= 0 && address <= 32767 {
		return address
	} else {
		panic("invalid index")
	}
}

func isRegisterAddress(address int) (bool, int) {
	if address >= 32768 && address <= 32775 {
		return true, address - 32768
	}
	return false, address
}

func (vm *VirtualMachine) set() {
	a, b := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2])
	vm.writeValueToAddress(b, a)
	vm.currentInstruction += 3
}

func (vm *VirtualMachine) push() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	vm.stack = append([]int{a}, vm.stack...)
	vm.currentInstruction += 2
}

func (vm *VirtualMachine) pop() {
	if len(vm.stack) == 0 {
		panic("empty stack")
	}
	value := vm.stack[0]
	vm.stack = vm.stack[1:]
	vm.writeValueToAddress(value, vm.memory[vm.currentInstruction+1])
	vm.currentInstruction += 2
}

func (vm *VirtualMachine) eq() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	if b == c {
		vm.writeValueToAddress(1, a)
	} else {
		vm.writeValueToAddress(0, a)
	}
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) gt() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	if b > c {
		vm.writeValueToAddress(1, a)
	} else {
		vm.writeValueToAddress(0, a)
	}
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) jmp() {
	vm.currentInstruction = vm.value(vm.memory[vm.currentInstruction+1])
}

func (vm *VirtualMachine) jt() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	if a != 0 {
		vm.currentInstruction = vm.value(vm.memory[vm.currentInstruction+2])
	} else {
		vm.currentInstruction += 2
	}
}

func (vm *VirtualMachine) jf() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	if a == 0 {
		vm.currentInstruction = vm.value(vm.memory[vm.currentInstruction+2])
	} else {
		vm.currentInstruction += 2
	}
}

func (vm *VirtualMachine) add() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.writeValueToAddress((b + c) % modValue, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) mult() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.writeValueToAddress((b * c) % modValue, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) mod() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.writeValueToAddress((b % c) % modValue, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) and() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.writeValueToAddress(b & c, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) or() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.writeValueToAddress((b | c) % modValue, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) not() {
	a, b := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2])
	vm.writeValueToAddress((^b) % modValue, a)
	vm.currentInstruction += 3
}

func (vm *VirtualMachine) rmem() {
	a, b := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2])
	vm.writeValueToAddress(vm.value(vm.memory[b]), a)
	vm.currentInstruction += 3
}

func (vm *VirtualMachine) wmem() {
	a, b := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2])
	vm.writeValueToAddress(b, a)
	vm.currentInstruction += 3
}

func (vm *VirtualMachine) call() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	vm.stack = append([]int{vm.currentInstruction+2}, vm.stack...)
	vm.currentInstruction = a
}

func (vm *VirtualMachine) ret() {
	if len(vm.stack) == 0 {
		panic("empty stack")
	}
	value := vm.stack[0]
	vm.stack = vm.stack[1:]
	vm.currentInstruction = value
}

func (vm *VirtualMachine) out() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	fmt.Print(string(rune(a)))
	vm.currentInstruction += 2
}

func (vm *VirtualMachine) in() {
	a := vm.memory[vm.currentInstruction+1]
	char, _, err := vm.reader.ReadRune()
	if err != nil {
		panic(err)
	}
	vm.writeValueToAddress(int(char), a)
	vm.currentInstruction += 2
}

func (vm *VirtualMachine) noop() {
	vm.currentInstruction++
}
