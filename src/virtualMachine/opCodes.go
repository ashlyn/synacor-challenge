package virtualmachine

// OpCode type represents supported opcodes specified in `arch-spec`
type OpCode int

// Enum of all OpCodes
const (
	Halt OpCode = 0
	Set = 1
	Push = 2
	Pop = 3
	Eq = 4
	Gt = 5
	Jmp = 6
	Jt = 7
	Jf = 8
	Add = 9
	Mult = 10
	Mod = 11
	And = 12
	Or = 13
	Not = 14
	Rmem = 15
	Wmem = 16
	Call = 17
	Ret = 18
	Out = 19
	In = 20
	Noop = 21
)

func (op OpCode) String() string {
	switch op {
	case Halt:
		return "halt"
	case Set:
		return "set"
	case Push:
		return "push"
	case Pop:
		return "pop"
	case Eq:
		return "eq"
	case Gt:
		return "gt"
	case Jmp:
		return "jmp"
	case Jt:
		return "jt"
	case Jf:
		return "jf"
	case Add:
		return "add"
	case Mult:
		return "mult"
	case Mod:
		return "mod"
	case And:
		return "and"
	case Or:
		return "or"
	case Not:
		return "not"
	case Rmem:
		return "rmem"
	case Wmem:
		return "wmem"
	case Call:
		return "call"
	case Ret:
		return "ret"
	case Out:
		return "out"
	case In:
		return "in"
	case Noop:
		return "noop"
	default:
		return "unknown"
	}
}

func (op OpCode) arguments() int {
	switch op {
	case Halt:
		return 0
	case Set:
		return 2
	case Push:
		return 1
	case Pop:
		return 1
	case Eq:
		return 3
	case Gt:
		return 3
	case Jmp:
		return 1
	case Jt:
		return 2
	case Jf:
		return 2
	case Add:
		return 3
	case Mult:
		return 3
	case Mod:
		return 3
	case And:
		return 3
	case Or:
		return 3
	case Not:
		return 2
	case Rmem:
		return 2
	case Wmem:
		return 2
	case Call:
		return 1
	case Ret:
		return 0
	case Out:
		return 1
	case In:
		return 1
	case Noop:
		return 0
	default:
		return 0
	}
}
