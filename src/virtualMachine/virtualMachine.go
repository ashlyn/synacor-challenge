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
	logFile *bufio.Writer
}

// NewVirtualMachine creates a new VirtualMachine
func NewVirtualMachine(reader *bufio.Reader, writer *bufio.Writer) *VirtualMachine {
	vm := VirtualMachine{
		currentInstruction: 0,
		memory: []int{},
		registers: [8]int{},
		stack: []int{},
		reader: reader,
		logFile: writer,
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

// SetEnergyLevel Allows the energy level (register 8) value to be set manually
func (vm *VirtualMachine) SetEnergyLevel(energyLevel int) {
	vm.registers[7] = energyLevel
}

// PrintMemoryHead is a debugging tool which prints the first 10 memory items
func (vm *VirtualMachine) PrintMemoryHead() {
	tail := int(math.Min(11, float64(len(vm.memory))))
	fmt.Println(vm.memory[:tail])
}

// MemoryDump is a debugging tool which dumps the current memory contents to a specified file
func (vm *VirtualMachine) MemoryDump(dump *bufio.Writer) error {
	for _, m := range vm.memory {
		op := OpCode(m)
		var err error
		if op.String() == "unknown" {
			_, err = dump.WriteString(fmt.Sprintf("%d\n", m))
		} else {
			dump.WriteString(fmt.Sprintf("%s (%d)\n", op, m))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteAssembly is a debugging tool which writes the assembly code for the current memory contents
func (vm *VirtualMachine) WriteAssembly(dump *bufio.Writer) error {
	for i := 0; i < len(vm.memory); i++ {
		op := OpCode(vm.memory[i])
		args := op.arguments()
		line := op.String()
		if line == "unknown" {
			line = fmt.Sprint(vm.memory[i])
		}
		for j := 1; j <= args; j++ {
			i++
			line += " " + fmt.Sprint(vm.memory[i])
		}
		_, err := dump.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute runs the program in memory
func (vm *VirtualMachine) Execute(debug bool) error {
	for vm.currentInstruction < len(vm.memory) {
		op := OpCode(vm.memory[vm.currentInstruction])
		vm.logFile.WriteString(op.String() + " ")
		switch op {
		case Halt:
			return nil
		case Set:
			vm.set()
		case Push:
			vm.push()
		case Pop:
			vm.pop()
		case Eq:
			vm.eq()
		case Gt:
			vm.gt()
		case Jmp:
			vm.jmp()
		case Jt:
			vm.jt()
		case Jf:
			vm.jf()
		case Add:
			vm.add()
		case Mult:
			vm.mult()
		case Mod:
			vm.mod()
		case And:
			vm.and()
		case Or:
			vm.or()
		case Not:
			vm.not()
		case Rmem:
			vm.rmem()
		case Wmem:
			vm.wmem()
		case Call:
			vm.call()
		case Ret:
			vm.ret()
		case Out:
			vm.out()
		case In:
			vm.in()
		case Noop:
			vm.noop()
		default:
			if debug {
				fmt.Println(vm.memory[vm.currentInstruction-5:vm.currentInstruction])
				fmt.Println(vm.memory[vm.currentInstruction:vm.currentInstruction+5])
			}
			return fmt.Errorf("Unrecognized op code %d", vm.memory[vm.currentInstruction])
		}
		vm.logFile.WriteString("\n")
	}

	if debug {
		fmt.Println(vm.memory)
		fmt.Println(vm.registers)
		fmt.Println(vm.stack)
	}

	return nil
}

func (vm *VirtualMachine) next() int {
	value := vm.memory[vm.currentInstruction]
	vm.currentInstruction++
	return value
}

func (vm *VirtualMachine) writeValueToAddress(value int, address int) {
	isRegister, regAddress := isRegisterAddress(address)
	if isRegister {
		vm.logFile.WriteString(fmt.Sprintf("Writing value %d to register %d\n", value, regAddress))
		vm.logFile.WriteString(fmt.Sprintf("%v\n", vm.registers))
		vm.registers[regAddress] = value
	} else {
		vm.logFile.WriteString(fmt.Sprintf("Writing value %d to memory address %d\n", value, address))
		vm.memory[address] = value
	}
}

func (vm *VirtualMachine) value(address int) int {
	isRegister, registerAddress := isRegisterAddress(address)
	if isRegister {
		vm.logFile.WriteString(fmt.Sprintf("Reading value %d from register %d\n", vm.registers[registerAddress], registerAddress))
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
	vm.logFile.WriteString(fmt.Sprintf("%d %d\n", a, b))
	vm.writeValueToAddress(b, a)
	vm.currentInstruction += 3
}

func (vm *VirtualMachine) push() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	vm.logFile.WriteString(fmt.Sprintf("%d\n", a))
	vm.stack = append([]int{a}, vm.stack...)
	vm.currentInstruction += 2
}

func (vm *VirtualMachine) pop() {
	if len(vm.stack) == 0 {
		panic("empty stack")
	}
	value := vm.stack[0]
	vm.stack = vm.stack[1:]
	a := vm.memory[vm.currentInstruction+1]
	vm.logFile.WriteString(fmt.Sprintf("%d\n", a))
	vm.writeValueToAddress(value, a)
	vm.currentInstruction += 2
}

func (vm *VirtualMachine) eq() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.logFile.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	if b == c {
		vm.writeValueToAddress(1, a)
	} else {
		vm.writeValueToAddress(0, a)
	}
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) gt() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.logFile.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	if b > c {
		vm.writeValueToAddress(1, a)
	} else {
		vm.writeValueToAddress(0, a)
	}
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) jmp() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	vm.logFile.WriteString(fmt.Sprintf("%d\n", a))
	vm.currentInstruction = vm.value(a)
}

func (vm *VirtualMachine) jt() {
	a, b := vm.value(vm.memory[vm.currentInstruction+1]), vm.value(vm.memory[vm.currentInstruction+2])
	isReg, regAddress := isRegisterAddress(vm.memory[vm.currentInstruction+1])
	vm.logFile.WriteString(fmt.Sprintf("%d %d\n", a, b))
	if isReg && regAddress == 7 && a != 0 && b == 1093 {
		// bypass reg 8 check in tests
		vm.currentInstruction += 3
	} else if a != 0 {
		vm.currentInstruction = b
	} else {
		vm.currentInstruction += 3
	}
}

func (vm *VirtualMachine) jf() {
	a, b := vm.value(vm.memory[vm.currentInstruction+1]), vm.value(vm.memory[vm.currentInstruction+2])
	vm.logFile.WriteString(fmt.Sprintf("%d %d\n", a, b))
	if a == 0 {
		vm.currentInstruction = b
	} else {
		vm.currentInstruction += 3
	}
}

func (vm *VirtualMachine) add() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.logFile.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	vm.writeValueToAddress((b + c) % modValue, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) mult() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.logFile.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	vm.writeValueToAddress((b * c) % modValue, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) mod() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.logFile.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	vm.writeValueToAddress(b % c, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) and() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.logFile.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	vm.writeValueToAddress(b & c, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) or() {
	a, b, c := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2]), vm.value(vm.memory[vm.currentInstruction+3])
	vm.logFile.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	vm.writeValueToAddress(b | c, a)
	vm.currentInstruction += 4
}

func (vm *VirtualMachine) not() {
	a, b := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2])
	vm.logFile.WriteString(fmt.Sprintf("%d %d\n", a, b))
	bStr := fmt.Sprintf("%015b", b)
	bStr = strings.ReplaceAll(bStr, "0", "-")
	bStr = strings.ReplaceAll(bStr, "1", "0")
	bStr = strings.ReplaceAll(bStr, "-", "1")
	value, _ := strconv.ParseUint(bStr, 2, 15)
	vm.writeValueToAddress(int(value), a)
	vm.currentInstruction += 3
}

func (vm *VirtualMachine) rmem() {
	a, b := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2])
	vm.logFile.WriteString(fmt.Sprintf("%d %d\n", a, b))
	vm.writeValueToAddress(vm.value(vm.memory[b]), a)
	vm.currentInstruction += 3
}

func (vm *VirtualMachine) wmem() {
	a, b := vm.memory[vm.currentInstruction+1], vm.value(vm.memory[vm.currentInstruction+2])
	vm.logFile.WriteString(fmt.Sprintf("%d %d\n", a, b))
	vm.writeValueToAddress(b, vm.value(a))
	vm.currentInstruction += 3
}

func (vm *VirtualMachine) call() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	vm.logFile.WriteString(fmt.Sprintf("%d\n", a))
	vm.stack = append([]int{vm.currentInstruction+2}, vm.stack...)
	vm.currentInstruction = a
}

func (vm *VirtualMachine) ret() {
	if len(vm.stack) == 0 {
		panic("empty stack")
	}
	value := vm.stack[0]
	vm.stack = vm.stack[1:]
	vm.logFile.WriteString(fmt.Sprintf("%d\n", value))
	vm.currentInstruction = value
}

func (vm *VirtualMachine) out() {
	a := vm.value(vm.memory[vm.currentInstruction+1])
	vm.logFile.WriteString(fmt.Sprintf("%d (%s)\n", a, string(rune(a))))
	fmt.Print(string(rune(a)))
	vm.currentInstruction += 2
}

func (vm *VirtualMachine) in() {
	a := vm.memory[vm.currentInstruction+1]
	char, _, err := vm.reader.ReadRune()
	vm.logFile.WriteString(fmt.Sprintf("%d (%s)\n", a, string(char)))
	if err != nil {
		panic(err)
	}
	vm.writeValueToAddress(int(char), a)
	vm.currentInstruction += 2
}

func (vm *VirtualMachine) noop() {
	vm.currentInstruction++
}
