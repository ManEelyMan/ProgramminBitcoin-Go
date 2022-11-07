package btc

import (
	"bitcoin-go/collections"
	"bitcoin-go/ecc"
	"bitcoin-go/utility"
	"bytes"
	"math"
	"math/big"
)

type ExecutionContext struct {
	Stack    *collections.Stack
	AltStack *collections.Stack
	Hash     *big.Int
}

type opFxn func(*ExecutionContext) bool

type twoIntOpComparator func(int64, int64) bool
type oneIntOpChooser func(int64) int64
type twoIntOpChooser func(int64, int64) int64
type oneBinaryOpChooser func([]byte) []byte

var opCodeFxns map[byte](opFxn)
var opCodeNames map[byte](string)

func init() {

	// Set up map of opcode to method
	opCodeFxns = make(map[byte]opFxn)
	opCodeFxns[0x00] = opFalse
	// opCodeFxns[0x4c] = opPushData1
	// opCodeFxns[0x4d] = opPushData2
	// opCodeFxns[0x4e] = opPushData4
	// opCodeFxns[0x4f] = opNeg1
	// opCodeFxns[0x51] = op1
	// opCodeFxns[0x52] = op2
	// opCodeFxns[0x53] = op3
	// opCodeFxns[0x54] = op4
	// opCodeFxns[0x55] = op5
	// opCodeFxns[0x56] = op6
	// opCodeFxns[0x57] = op7
	// opCodeFxns[0x58] = op8
	// opCodeFxns[0x59] = op9
	// opCodeFxns[0x5a] = op10
	// opCodeFxns[0x5b] = op11
	// opCodeFxns[0x5c] = op12
	// opCodeFxns[0x5d] = op13
	// opCodeFxns[0x5e] = op14
	// opCodeFxns[0x5f] = op15
	// opCodeFxns[0x60] = op16
	opCodeFxns[0x61] = opNop
	// opCodeFxns[0x63] = opIf
	// opCodeFxns[0x64] = opNotIf
	// opCodeFxns[0x67] = opElse
	// opCodeFxns[0x68] = opEndIf
	opCodeFxns[0x69] = opVerify
	opCodeFxns[0x6a] = opReturn
	opCodeFxns[0x6b] = opToAltStack
	opCodeFxns[0x6c] = opFromAltStack
	opCodeFxns[0x73] = opIfDupe
	opCodeFxns[0x74] = opDepth
	opCodeFxns[0x75] = opDrop
	opCodeFxns[0x76] = opDupe
	opCodeFxns[0x77] = opNip
	opCodeFxns[0x78] = opOver
	opCodeFxns[0x79] = opPick
	opCodeFxns[0x7a] = opRoll
	opCodeFxns[0x7b] = opRot
	opCodeFxns[0x7c] = opSwap
	opCodeFxns[0x7d] = opTuck
	opCodeFxns[0x6d] = opDrop2
	opCodeFxns[0x6e] = opDupe2
	opCodeFxns[0x6f] = opDupe3
	opCodeFxns[0x70] = opOver2
	opCodeFxns[0x71] = opRot2
	opCodeFxns[0x72] = opSwap2
	opCodeFxns[0x82] = opSize
	opCodeFxns[0x87] = opEqual
	opCodeFxns[0x88] = opEqualVerify
	opCodeFxns[0x8b] = opAdd1
	opCodeFxns[0x8c] = opSub1
	opCodeFxns[0x8f] = opNegate
	opCodeFxns[0x90] = opAbs
	opCodeFxns[0x91] = opNot
	opCodeFxns[0x92] = opNotEqual0
	opCodeFxns[0x93] = opAdd
	opCodeFxns[0x94] = opSub
	opCodeFxns[0x9a] = opBoolAnd
	opCodeFxns[0x9b] = opBoolOr
	opCodeFxns[0x9c] = opNumEqual
	opCodeFxns[0x9d] = opNumEqualVerify
	opCodeFxns[0x9e] = opNumNotEqual
	opCodeFxns[0x9f] = opLessThan
	opCodeFxns[0xa0] = opGreaterThan
	opCodeFxns[0xa1] = opLessThanOrEqual
	opCodeFxns[0xa2] = opGreaterThanOrEqual
	opCodeFxns[0xa3] = opMin
	opCodeFxns[0xa4] = opMax
	opCodeFxns[0xa5] = opWithin
	opCodeFxns[0xa6] = opRipeMd160
	opCodeFxns[0xa7] = opSha1
	opCodeFxns[0xa8] = opSha256
	opCodeFxns[0xa9] = opHash160
	opCodeFxns[0xaa] = opHash256
	opCodeFxns[0xab] = opCodeSeparator
	opCodeFxns[0xac] = opCheckSig
	opCodeFxns[0xad] = opCheckSigVerify
	opCodeFxns[0xae] = opCheckMultiSig
	opCodeFxns[0xaf] = opCheckMultiSigVerify
	opCodeFxns[0xb1] = opCheckLockTimeVerify
	opCodeFxns[0xb2] = opCheckSequenceVerify

	// Reserved
	opCodeFxns[0x50] = opAutoFail
	opCodeFxns[0x62] = opAutoFail
	opCodeFxns[0x65] = opAutoFail
	opCodeFxns[0x66] = opAutoFail
	opCodeFxns[0x89] = opAutoFail
	opCodeFxns[0x8a] = opAutoFail
	opCodeFxns[0xb0] = opNop
	opCodeFxns[0xb3] = opNop
	opCodeFxns[0xb4] = opNop
	opCodeFxns[0xb5] = opNop
	opCodeFxns[0xb6] = opNop
	opCodeFxns[0xb7] = opNop
	opCodeFxns[0xb8] = opNop
	opCodeFxns[0xb9] = opNop

	opCodeNames = make(map[byte]string)
	opCodeNames[0x00] = "OP_0"
	opCodeNames[0x4f] = "OP_1NEGATE"
	opCodeNames[0x51] = "OP_1"
	opCodeNames[0x52] = "OP_2"
	opCodeNames[0x53] = "OP_3"
	opCodeNames[0x54] = "OP_4"
	opCodeNames[0x55] = "OP_5"
	opCodeNames[0x56] = "OP_6"
	opCodeNames[0x57] = "OP_7"
	opCodeNames[0x58] = "OP_8"
	opCodeNames[0x59] = "OP_9"
	opCodeNames[0x5a] = "OP_10"
	opCodeNames[0x5b] = "OP_11"
	opCodeNames[0x5c] = "OP_12"
	opCodeNames[0x5d] = "OP_13"
	opCodeNames[0x5e] = "OP_14"
	opCodeNames[0x5f] = "OP_15"
	opCodeNames[0x60] = "OP_16"
	opCodeNames[0x61] = "OP_NOP"
	opCodeNames[0x69] = "OP_VERIFY"
	opCodeNames[0x6a] = "OP_RETURN"
	opCodeNames[0x6b] = "OP_TOALTSTACK"
	opCodeNames[0x6c] = "OP_FROMALTSTACK"
	opCodeNames[0x73] = "OP_IFDUP"
	opCodeNames[0x74] = "OP_DEPTH"
	opCodeNames[0x75] = "OP_DROP"
	opCodeNames[0x76] = "OP_DUP"
	opCodeNames[0x77] = "OP_NIP"
	opCodeNames[0x78] = "OP_OVER"
	opCodeNames[0x79] = "OP_PICK"
	opCodeNames[0x7a] = "OP_ROLL"
	opCodeNames[0x7b] = "OP_ROT"
	opCodeNames[0x7c] = "OP_SWAP"
	opCodeNames[0x7d] = "OP_TUCK"
	opCodeNames[0x6d] = "OP_2DROP"
	opCodeNames[0x6e] = "OP_2DUP"
	opCodeNames[0x6f] = "OP_3DUP"
	opCodeNames[0x70] = "OP_2OVER"
	opCodeNames[0x71] = "OP_2ROT"
	opCodeNames[0x72] = "OP_2SWAP"
	opCodeNames[0x82] = "OP_SIZE"
	opCodeNames[0x87] = "OP_EQUAL"
	opCodeNames[0x88] = "OP_EQUALVERIFY"
	opCodeNames[0x8b] = "OP_1ADD"
	opCodeNames[0x8c] = "OP_1SUB"
	opCodeNames[0x8f] = "OP_NEGATE"
	opCodeNames[0x90] = "OP_ABS"
	opCodeNames[0x91] = "OP_NOT"
	opCodeNames[0x92] = "OP_0NOTEQUAL"
	opCodeNames[0x93] = "OP_ADD"
	opCodeNames[0x94] = "OP_SUB"
	opCodeNames[0x9a] = "OP_BOOLAND"
	opCodeNames[0x9b] = "OP_BOOLOR"
	opCodeNames[0x9c] = "OP_NUMEQUAL"
	opCodeNames[0x9d] = "OP_NUMEQUALVERIFY"
	opCodeNames[0x9e] = "OP_NUMNOTEQUAL"
	opCodeNames[0x9f] = "OP_LESSTHAN"
	opCodeNames[0xa0] = "OP_GREATERTHAN"
	opCodeNames[0xa1] = "OP_LESSTHANOREQUAL"
	opCodeNames[0xa2] = "OP_GREATERTHANOREQUAL"
	opCodeNames[0xa3] = "OP_MIN"
	opCodeNames[0xa4] = "OP_MAX"
	opCodeNames[0xa5] = "OP_WITHIN"
	opCodeNames[0xa6] = "OP_RIPEMD160"
	opCodeNames[0xa7] = "OP_SHA1"
	opCodeNames[0xa8] = "OP_SHA256"
	opCodeNames[0xa9] = "OP_HASH160"
	opCodeNames[0xaa] = "OP_HASH256"
	opCodeNames[0xab] = "OP_CODESEPARATOR"
	opCodeNames[0xac] = "OP_CHECKSIG"
	opCodeNames[0xad] = "OP_CHECKSIGVERIFY"
	opCodeNames[0xae] = "OP_CHECKMULTISIG"
	opCodeNames[0xaf] = "OP_CHECKMULTISIGVERIFY"
	opCodeNames[0xb1] = "OP_CHECKLOCKTIMEVERIFY"
	opCodeNames[0xb2] = "OP_CHECKSEQUENCEVERIFY"

	// Reserved
	opCodeNames[0x50] = "OP_RESERVED"
	opCodeNames[0x62] = "OP_VER"
	opCodeNames[0x65] = "OP_VERIF"
	opCodeNames[0x66] = "OP_VERNOTIF"
	opCodeNames[0x89] = "OP_RESERVED1"
	opCodeNames[0x8a] = "OP_RESERVED2"
	opCodeNames[0xb0] = "OP_NOP1"
	opCodeNames[0xb3] = "OP_NOP4"
	opCodeNames[0xb4] = "OP_NOP5"
	opCodeNames[0xb5] = "OP_NOP6"
	opCodeNames[0xb6] = "OP_NOP7"
	opCodeNames[0xb7] = "OP_NOP8"
	opCodeNames[0xb8] = "OP_NOP9"
	opCodeNames[0xb9] = "OP_NOP10"
}

func encodeNumber(num int64) []byte {
	buffer := make([]byte, 0)
	if num == 0 {
		return buffer
	}

	negative := num < 0
	abs := (int64)(math.Abs(float64(num)))

	for abs != 0 {
		tmp := abs & 0xFF
		buffer = append(buffer, (byte)(tmp))
		abs = abs >> 8
	}

	if buffer[len(buffer)-1]&0x80 != 0 {
		buffer = append(buffer, utility.IIF(negative, 0x80, 0x00).(byte))
	} else if negative {
		buffer[len(buffer)-1] = buffer[len(buffer)-1] | 0x80
	}

	return buffer
}

func decodeNumber(buffer []byte) int64 {
	len := len(buffer)
	if len == 0 {
		return 0
	}

	negative := buffer[len-1]&0x80 != 0
	buffer = utility.ReverseBytes(buffer) // Make big-endian
	buffer[0] = buffer[0] & 0x7f

	var result int64 = 0
	for _, b := range buffer {
		result = (result << 8) + (int64)(b)
	}

	if negative {
		return -result
	} else {
		return result
	}
}

func pushNumberOnStack(stack *collections.Stack, num int64) bool {
	// The number -1 is pushed onto the stack.
	buff := encodeNumber(num)
	stack.Push(buff)
	return true
}

func twoIntCompareOp(stack *collections.Stack, comparer twoIntOpComparator) bool {
	if stack.Length() < 2 {
		return false
	}

	top, _ := stack.Pop()
	bottom, _ := stack.Pop()

	b := decodeNumber(top)
	a := decodeNumber(bottom)
	var res int64 = 0

	if comparer(a, b) {
		res = 1
	}

	stack.Push(encodeNumber(res))
	return true
}

func twoIntChooseOp(stack *collections.Stack, chooser twoIntOpChooser) bool {
	if stack.Length() < 2 {
		return false
	}

	top, _ := stack.Pop()
	bottom, _ := stack.Pop()

	b := decodeNumber(top)
	a := decodeNumber(bottom)
	res := chooser(a, b)

	stack.Push(encodeNumber(res))
	return true
}

func oneIntChooseOp(stack *collections.Stack, chooser oneIntOpChooser) bool {
	top, ok := stack.Pop()
	if !ok || top == nil {
		return false
	}

	a := decodeNumber(top)
	res := chooser(a)

	stack.Push(encodeNumber(res))
	return true
}

func oneBinaryChooseOp(stack *collections.Stack, chooser oneBinaryOpChooser) bool {
	item, ok := stack.Pop()
	if !ok || item == nil {
		return false
	}
	res := chooser(item)

	stack.Push(res)
	return true
}

func opFalse(context *ExecutionContext) bool {
	// An empty array of bytes is pushed onto the stack. (This is not a no-op: an item is added to the stack.)
	context.Stack.Push(make([]byte, 0))
	return true
}

// func opNeg1(context *ExecutionContext) bool {
// 	// The number -1 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, -1)
// }
// func op1(context *ExecutionContext) bool {
// 	// The number 1 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 1)
// }
// func op2(context *ExecutionContext) bool {
// 	// The number 2 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 2)
// }
// func op3(context *ExecutionContext) bool {
// 	// The number 3 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 3)
// }
// func op4(context *ExecutionContext) bool {
// 	// The number 4 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 4)
// }
// func op5(context *ExecutionContext) bool {
// 	// The number 5 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 5)
// }
// func op6(context *ExecutionContext) bool {
// 	// The number 6 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 6)
// }
// func op7(context *ExecutionContext) bool {
// 	// The number 7 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 7)
// }
// func op8(context *ExecutionContext) bool {
// 	// The number 8 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 8)
// }
// func op9(context *ExecutionContext) bool {
// 	// The number 9 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 9)
// }
// func op10(context *ExecutionContext) bool {
// 	// The number 10 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 10)
// }
// func op11(context *ExecutionContext) bool {
// 	// The number 11 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 11)
// }
// func op12(context *ExecutionContext) bool {
// 	// The number 12 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 12)
// }
// func op13(context *ExecutionContext) bool {
// 	// The number 13 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 13)
// }
// func op14(context *ExecutionContext) bool {
// 	// The number 14 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 14)
// }
// func op15(context *ExecutionContext) bool {
// 	// The number 15 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 15)
// }
// func op16(context *ExecutionContext) bool {
// 	// The number 16 is pushed onto the stack.
// 	return pushNumberOnStack(context.Stack, 16)
// }

func opVerify(context *ExecutionContext) bool {
	// Marks transaction as invalid if top stack value is not true. The top stack value is removed.
	b, ok := context.Stack.Pop()
	if !ok || b == nil {
		return false
	}

	num := decodeNumber(b)
	return num != 0
}
func opReturn(context *ExecutionContext) bool {
	// Marks transaction as invalid.
	return false
}
func opToAltStack(context *ExecutionContext) bool {
	// Puts the input onto the top of the alt stack. Removes it from the main stack.
	top, ok := context.Stack.Pop()
	if !ok || top == nil {
		return false
	}

	context.AltStack.Push(top)
	return true
}
func opFromAltStack(context *ExecutionContext) bool {
	// Puts the input onto the top of the main stack. Removes it from the alt stack.
	top, ok := context.AltStack.Pop()
	if !ok || top == nil {
		return false
	}

	context.Stack.Push(top)
	return true
}
func opIfDupe(context *ExecutionContext) bool {
	// If the top stack value is not 0, duplicate it.
	b, ok := context.Stack.Peek()
	if !ok {
		return false
	}

	if decodeNumber(b) != 0 {
		context.Stack.Push(b)
	}

	return true
}
func opDepth(context *ExecutionContext) bool {
	// Puts the number of stack items onto the stack.
	size := context.Stack.Length()
	context.Stack.Push(encodeNumber((int64)(size)))
	return true
}
func opDrop(context *ExecutionContext) bool {
	// Removes the top stack item.
	_, ok := context.Stack.Pop()
	return ok
}
func opDupe(context *ExecutionContext) bool {
	// Duplicates the top stack item.
	b, ok := context.Stack.Peek()
	if !ok {
		return false
	}

	context.Stack.Push(b)
	return true
}
func opNip(context *ExecutionContext) bool {
	// Removes the second-to-top stack item.
	if context.Stack.Length() < 2 {
		return false
	}

	first, _ := context.Stack.Pop()
	context.Stack.Pop()       // Burn the 2nd item
	context.Stack.Push(first) // Put the first back.
	return true
}
func opOver(context *ExecutionContext) bool {
	// Copies the second-to-top stack item to the top.
	if context.Stack.Length() < 2 {
		return false
	}

	item, ok := context.Stack.PeekAt(1)
	if !ok {
		return false
	}

	context.Stack.Push(item)
	return true
}
func opPick(context *ExecutionContext) bool {
	// The item n back in the stack is copied to the top.
	top, ok := context.Stack.Pop()
	if !ok || top == nil {
		return false
	}

	n := decodeNumber(top)
	index := (uint32)(n + 1)
	if context.Stack.Length() < index {
		return false
	}

	item, ok := context.Stack.PeekAt(index)
	if !ok {
		return false
	}

	context.Stack.Push(item)
	return true
}
func opRoll(context *ExecutionContext) bool {
	// The item n back in the stack is moved to the top.

	// Determine the index of the item to move to the top
	item, ok := context.Stack.Pop()
	if !ok || item == nil {
		return false
	}

	n := decodeNumber(item)
	index := (uint32)(n + 1)

	// Verify there's enough items on the stack.
	if context.Stack.Length() < index {
		return false
	}

	// Briefly move n-1 items to altstack.
	var i uint32 = 0
	for i = 0; i < index-1; i++ {
		top, _ := context.Stack.Pop()
		context.AltStack.Push(top)
	}

	// Grab the item we want.
	item, _ = context.Stack.Pop()

	// Move the n-1 items back off the altstack.
	for i = 0; i < index-1; i++ {
		top, _ := context.Stack.Pop()
		context.Stack.Push(top)
	}

	// Move our desired item back to the top of the stack.
	context.Stack.Push(item)
	return true
}
func opRot(context *ExecutionContext) bool {
	// The 3rd item down the stack is moved to the top.
	if context.Stack.Length() < 3 {
		return false
	}

	one, _ := context.Stack.Pop()
	two, _ := context.Stack.Pop()
	three, _ := context.Stack.Pop()

	context.Stack.Push(two)
	context.Stack.Push(one)
	context.Stack.Push(three)
	return true
}
func opSwap(context *ExecutionContext) bool {
	// The top two items on the stack are swapped.
	if context.Stack.Length() < 2 {
		return false
	}

	one, _ := context.Stack.Pop()
	two, _ := context.Stack.Pop()

	context.Stack.Push(one)
	context.Stack.Push(two)
	return true
}
func opTuck(context *ExecutionContext) bool {
	// The item at the top of the stack is copied and inserted before the second-to-top item.
	if context.Stack.Length() < 2 {
		return false
	}

	one, _ := context.Stack.Pop()
	two, _ := context.Stack.Pop()

	context.Stack.Push(one)
	context.Stack.Push(two)
	context.Stack.Push(one)
	return true
}
func opDrop2(context *ExecutionContext) bool {
	// Removes the top two stack items.
	if context.Stack.Length() < 2 {
		return false
	}

	context.Stack.Pop()
	context.Stack.Pop()
	return true
}
func opDupe2(context *ExecutionContext) bool {
	// Duplicates the top two stack items.
	if context.Stack.Length() < 2 {
		return false
	}

	one, ok := context.Stack.PeekAt(0)
	if !ok {
		return false
	}
	two, ok := context.Stack.PeekAt(1)
	if !ok {
		return false
	}

	context.Stack.Push(two)
	context.Stack.Push(one)
	return true
}
func opDupe3(context *ExecutionContext) bool {
	// Duplicates the top three stack items.
	if context.Stack.Length() < 3 {
		return false
	}

	one, ok := context.Stack.PeekAt(0)
	if !ok {
		return false
	}
	two, ok := context.Stack.PeekAt(1)
	if !ok {
		return false
	}
	three, ok := context.Stack.PeekAt(2)
	if !ok {
		return false
	}

	context.Stack.Push(three)
	context.Stack.Push(two)
	context.Stack.Push(one)
	return true
}
func opOver2(context *ExecutionContext) bool {
	// Copies the pair of items two spaces back in the stack to the front.
	if context.Stack.Length() < 4 {
		return false
	}

	three, ok := context.Stack.PeekAt(2)
	if !ok {
		return false
	}
	four, ok := context.Stack.PeekAt(3)
	if !ok {
		return false
	}

	context.Stack.Push(four)
	context.Stack.Push(three)
	return true
}
func opRot2(context *ExecutionContext) bool {
	// The fifth and sixth items back are moved to the top of the stack.
	if context.Stack.Length() < 6 {
		return false
	}

	one, _ := context.Stack.Pop()
	two, _ := context.Stack.Pop()
	three, _ := context.Stack.Pop()
	four, _ := context.Stack.Pop()
	five, _ := context.Stack.Pop()
	six, _ := context.Stack.Pop()

	context.Stack.Push(four)
	context.Stack.Push(three)
	context.Stack.Push(two)
	context.Stack.Push(one)
	context.Stack.Push(six)
	context.Stack.Push(five)
	return true
}
func opSwap2(context *ExecutionContext) bool {
	// Swaps the top two pairs of items.
	if context.Stack.Length() < 4 {
		return false
	}

	one, _ := context.Stack.Pop()
	two, _ := context.Stack.Pop()
	three, _ := context.Stack.Pop()
	four, _ := context.Stack.Pop()

	context.Stack.Push(two)
	context.Stack.Push(one)
	context.Stack.Push(four)
	context.Stack.Push(three)

	return true
}
func opSize(context *ExecutionContext) bool {
	// Pushes the string length of the top element of the stack (without popping it).
	b, ok := context.Stack.Peek()
	if !ok {
		return false
	}

	var len int64 = (int64)(len(b))
	context.Stack.Push(encodeNumber(len))
	return true
}
func opEqual(context *ExecutionContext) bool {
	// Returns 1 if the inputs are exactly equal, 0 otherwise
	if context.Stack.Length() < 2 {
		return false
	}

	a, _ := context.Stack.Pop()
	b, _ := context.Stack.Pop()

	var res int64 = 0

	if bytes.Equal(a, b) {
		res = 1
	}

	context.Stack.Push(encodeNumber(res))
	return true
}
func opEqualVerify(context *ExecutionContext) bool {
	// Same as OP_EQUAL, but runs OP_VERIFY afterward.
	return opEqual(context) && opVerify(context)
}
func opAdd1(context *ExecutionContext) bool {
	// 1 is added to the input.
	return oneIntChooseOp(context.Stack, func(i int64) int64 { return i + 1 })
}
func opSub1(context *ExecutionContext) bool {
	// 1 is subtracted from the input.
	return oneIntChooseOp(context.Stack, func(i int64) int64 { return i - 1 })
}
func opNegate(context *ExecutionContext) bool {
	// The sign of the input is flipped.
	return oneIntChooseOp(context.Stack, func(i int64) int64 { return -i })
}
func opAbs(context *ExecutionContext) bool {
	// The input is made positive.
	return oneIntChooseOp(context.Stack, func(i int64) int64 { return int64(math.Abs(float64(i))) })
}
func opNot(context *ExecutionContext) bool {
	// If the input is 0 or 1, it is flipped. Otherwise the output will be 0.
	return oneIntChooseOp(context.Stack, func(i int64) int64 {
		if i == 0 {
			return 1
		} else {
			return 0
		}
	})
}
func opNotEqual0(context *ExecutionContext) bool {
	// Returns 0 if the input is 0. 1 otherwise.
	return oneIntChooseOp(context.Stack, func(i int64) int64 {
		if i == 0 {
			return 0
		} else {
			return 1
		}
	})
}
func opAdd(context *ExecutionContext) bool {
	// a is added to b.
	return twoIntChooseOp(context.Stack, func(a, b int64) int64 { return a + b })
}
func opSub(context *ExecutionContext) bool {
	// b is subtracted from a.
	return twoIntChooseOp(context.Stack, func(a, b int64) int64 { return a - b })
}
func opBoolAnd(context *ExecutionContext) bool {
	// If both a and b are not 0, the output is 1. Otherwise 0
	return twoIntCompareOp(context.Stack, func(a, b int64) bool { return a != 0 && b != 0 })
}
func opBoolOr(context *ExecutionContext) bool {
	// If a or b is not 0, the output is 1. Otherwise 0.
	return twoIntCompareOp(context.Stack, func(a, b int64) bool { return a != 0 || b != 0 })
}
func opNumEqual(context *ExecutionContext) bool {
	// Returns 1 if the numbers are equal, 0 otherwise.
	return twoIntCompareOp(context.Stack, func(a, b int64) bool { return a == b })
}
func opNumEqualVerify(context *ExecutionContext) bool {
	// Same as OP_NUMEQUAL, but runs OP_VERIFY afterward.
	return opNumEqual(context) && opVerify(context)
}
func opNumNotEqual(context *ExecutionContext) bool {
	// Returns 1 if the numbers are not equal, 0 otherwise.
	return twoIntCompareOp(context.Stack, func(a, b int64) bool { return a != b })
}
func opLessThan(context *ExecutionContext) bool {
	// Returns 1 if a is less than b, 0 otherwise.
	return twoIntCompareOp(context.Stack, func(a, b int64) bool { return a < b })
}
func opGreaterThan(context *ExecutionContext) bool {
	// Returns 1 if a is greater than b, 0 otherwise.
	return twoIntCompareOp(context.Stack, func(a, b int64) bool { return a > b })
}
func opLessThanOrEqual(context *ExecutionContext) bool {
	// Returns 1 if a is less than or equal to b, 0 otherwise.
	return twoIntCompareOp(context.Stack, func(a, b int64) bool { return a <= b })
}
func opGreaterThanOrEqual(context *ExecutionContext) bool {
	// Returns 1 if a is greater than or equal to b, 0 otherwise.
	return twoIntCompareOp(context.Stack, func(a, b int64) bool { return a >= b })
}
func opMin(context *ExecutionContext) bool {
	// Returns the smaller of a and b.
	return twoIntChooseOp(context.Stack, func(a, b int64) int64 {
		if a <= b {
			return a
		} else {
			return b
		}
	})
}
func opMax(context *ExecutionContext) bool {
	// Returns the larger of a and b.
	return twoIntChooseOp(context.Stack, func(a, b int64) int64 {
		if a >= b {
			return a
		} else {
			return b
		}
	})
}
func opWithin(context *ExecutionContext) bool {
	// Returns 1 if x is within the specified range (left-inclusive), 0 otherwise.
	if context.Stack.Length() < 3 {
		return false
	}

	one, _ := context.Stack.Pop()
	two, _ := context.Stack.Pop()
	three, _ := context.Stack.Pop()

	x := decodeNumber(one)
	min := decodeNumber(two)
	max := decodeNumber(three)
	var res int64 = 0
	if x >= min && x < max {
		res = 1
	}

	context.Stack.Push(encodeNumber(res))
	return true
}
func opRipeMd160(context *ExecutionContext) bool {
	// The input is hashed using RIPEMD-160.
	return oneBinaryChooseOp(context.Stack, func(b []byte) []byte { return utility.HashRipemd160(b) })
}
func opSha1(context *ExecutionContext) bool {
	//The input is hashed using SHA-1.
	return oneBinaryChooseOp(context.Stack, func(b []byte) []byte { return utility.Sha1(b) })
}
func opSha256(context *ExecutionContext) bool {
	// The input is hashed using SHA-256.
	return oneBinaryChooseOp(context.Stack, func(b []byte) []byte { return utility.Sha256(b) })
}
func opHash160(context *ExecutionContext) bool {
	// The input is hashed twice: first with SHA-256 and then with RIPEMD-160.
	return oneBinaryChooseOp(context.Stack, func(b []byte) []byte { return utility.Hash160(b) })
}
func opHash256(context *ExecutionContext) bool {
	// The input is hashed two times with SHA-256.
	return oneBinaryChooseOp(context.Stack, func(b []byte) []byte { return utility.Hash256(b) })
}
func opCodeSeparator(context *ExecutionContext) bool {
	// All of the signature checking words will only match signatures to the data after the most recently-executed OP_CODESEPARATOR.
	return false
}
func opCheckSig(context *ExecutionContext) bool {
	// The entire transaction's outputs, inputs, and script (from the most recently-executed OP_CODESEPARATOR to the end) are hashed.
	// The signature used by OP_CHECKSIG must be a valid signature for this hash and public key. If it is, 1 is returned, 0 otherwise.

	if context.Stack.Length() < 2 {
		return false
	}

	secPubKey, _ := context.Stack.Pop()
	derSignature, _ := context.Stack.Pop() // Do I have to reverse these?

	point := ecc.NewPointFromSEC(secPubKey)
	sig := ecc.NewSignatureFromDER(derSignature)

	if point.Verify(context.Hash, sig) {
		context.Stack.Push(encodeNumber(1))
	} else {
		context.Stack.Push(encodeNumber(0))
	}

	return true
}
func opCheckSigVerify(context *ExecutionContext) bool {
	// Same as OP_CHECKSIG, but OP_VERIFY is executed afterward.
	return opCheckSig(context) && opVerify(context)
}
func opCheckMultiSig(context *ExecutionContext) bool {
	// Compares the first signature against each public key until it finds an ECDSA match.
	// Starting with the subsequent public key, it compares the second signature against each remaining public key until it finds an ECDSA match.
	// The process is repeated until all signatures have been checked or not enough public keys remain to produce a successful result.
	// All signatures need to match a public key. Because public keys are not checked again if they fail any signature comparison,
	// signatures must be placed in the scriptSig using the same order as their corresponding public keys were placed in the scriptPubKey or redeemScript.
	// If all signatures are valid, 1 is returned, 0 otherwise. Due to a bug, one extra unused value is removed from the stack.

	// Get 'n'
	tmp, ok := context.Stack.Pop()
	if !ok {
		return false
	}
	n := decodeNumber(tmp)

	// Get n+1 elements off the stack and convert to points.
	if int64(context.Stack.Length()) < n+1 {
		return false
	}

	pubKeys := make([]ecc.Point, n)
	for i := 0; int64(i) < n; i++ {
		tmp, _ := context.Stack.Pop()
		pubKeys[i] = ecc.NewPointFromSEC(tmp)
	}

	// Get 'm'
	tmp, ok = context.Stack.Pop()
	if !ok {
		return false
	}
	m := decodeNumber(tmp)

	// Get m+1 elements off of the stack and convert to signatures.
	if int64(context.Stack.Length()) < m+1 {
		return false
	}

	sigs := make([]ecc.Signature, m)
	for i := 0; int64(i) < m; i++ {
		tmp, _ := context.Stack.Pop()
		sigs[i] = ecc.NewSignatureFromDER(tmp[:len(tmp)-1]) // signature is assumed to be using SIGHASH_ALL
	}

	// OP_CHECKMULTISIG bug: Pop off one additional, unused element.
	_, ok = context.Stack.Pop()
	if !ok {
		return false
	}

	pointCounter := 0
	for i := 0; i < len(sigs); i++ {
		if pointCounter >= len(pubKeys) {
			return false
		}

		for len(pubKeys) > 0 {
			pk := pubKeys[pointCounter]
			pointCounter++

			if pk.Verify(context.Hash, sigs[i]) {
				break
			}
		}
	}

	context.Stack.Push(encodeNumber(1))
	return true
}
func opCheckMultiSigVerify(context *ExecutionContext) bool {
	// Same as OP_CHECKMULTISIG, but OP_VERIFY is executed afterward.
	return opCheckMultiSig(context) && opVerify(context)
}
func opCheckLockTimeVerify(context *ExecutionContext) bool {
	// Marks transaction as invalid if the top stack item is greater than the transaction's nLockTime field,
	// otherwise script evaluation continues as though an OP_NOP was executed. Transaction is also invalid if:
	//     1. the stack is empty;
	//  or 2. the top stack item is negative;
	//  or 3. the top stack item is greater than or equal to 500000000 while the transaction's nLockTime field is less than 500000000, or vice versa;
	//  or 4. the input's nSequence field is equal to 0xffffffff.
	// The precise semantics are described in BIP 0065.
	return false
}
func opCheckSequenceVerify(context *ExecutionContext) bool {
	// Marks transaction as invalid if the relative lock time of the input (enforced by BIP 0068 with nSequence) is not equal to or longer than the value of the top stack item. The precise semantics are described in BIP 0112.
	return false
}

func opNop(context *ExecutionContext) bool {
	// It's a no-op
	return true
}

func opAutoFail(context *ExecutionContext) bool {
	// If we attempt to use any of these, automatically fail.
	return false
}
