package btc

import (
	"bitcoin-go/collections"
	"bytes"
	"fmt"
	"io"
	"math/big"
)

type ScriptExecutor struct {
	scriptPubKey    *Script
	scriptSignature *Script
	stack           collections.Stack
	altStack        collections.Stack
	hash            *big.Int
}

func NewScriptExecutor(pubkey *Script, sig *Script, hash *big.Int) ScriptExecutor {
	return ScriptExecutor{scriptPubKey: pubkey, scriptSignature: sig, hash: hash}
}

func (ex *ScriptExecutor) Execute() bool {

	executionContext := ExecutionContext{Stack: &ex.stack, AltStack: &ex.altStack, Hash: ex.hash}

	// 1. Parse, load and execute script signature
	ok := executeScript(ex.scriptSignature.Data, &executionContext)
	if !ok {
		return false
	}
	// 2. Parse, load and execute script pub key
	ok = executeScript(ex.scriptPubKey.Data, &executionContext)
	if !ok {
		return false
	}

	// 3. Verify completion
	if ex.stack.Length() == 0 {
		return false
	}

	if b, _ := ex.stack.Pop(); len(b) == 0 {
		return false
	}

	return true
}

func executeScript(script []byte, context *ExecutionContext) bool {

	operations, err := PrepareScript(script)

	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, op := range operations {
		opCode := op.GetOpCode()

		// TODO: Make some decisions based on the op code!
		if opCode == 0x63 || opCode == 0x64 || opCode == 0x67 || opCode == 0x69 {
			// TODO: Figure out how to track the if/else/endif stuff.
		} else {
			fmt.Printf("Performing op %+v (%v)...\n", op.GetOpName(), opCode)
			ok := op.Execute(context)
			if !ok {
				fmt.Printf("Failed processing op %+v.\n", op.GetOpName())
				return false
			}
		}

	}

	fmt.Printf("%+v\n", operations)
	return true

	// for {
	// 	opCode, err := reader.ReadByte()
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			return true
	// 		} else {
	// 			fmt.Println(err)
	// 			return false
	// 		}
	// 	}

	// 	opCodeFxn, ok := opCodeFxns[opCode]
	// 	if !ok {
	// 		// This byte doesn't map to an op code. Determine what to do!
	// 		if opCode >= 0x01 && opCode <= 0x4e {
	// 			extractDataAndPushToStack(opCode, reader, stack)
	// 		} else if opCode == 0x63 || opCode == 0x64 || opCode == 0x67 || opCode == 0x69 {
	// 			// TODO: Figure out how to track the if/else/endif stuff.
	// 		}
	// 	} else {
	// 		opName := opCodeNames[opCode]
	// 		fmt.Printf("Performing op %+v (%v)...\n", opName, opCode)
	// 		success := opCodeFxn(reader, stack, altStack)
	// 		if !success {
	// 			fmt.Printf("Failed processing op %+v.\n", opName)
	// 			return false
	// 		}
	// 	}
	// }
}

// func extractDataAndPushToStack(op byte, reader io.Reader, stack *collections.Stack) bool {

// 	var dataLength uint32 = (uint32)(op)

// 	if op == 0x4c {
// 		dataLength = (uint32)(utility.ReadByte(reader))
// 	} else if op == 0x4d {
// 		dataLength = uint32(utility.ReadUInt16(reader, true))
// 	} else if op == 0x4e {
// 		dataLength = utility.ReadUint32(reader, true)
// 	}

// 	data := make([]byte, dataLength)
// 	n, err := reader.Read(data)

// 	if err != nil {
// 		fmt.Printf("Error reading in data: %+v\n", err)
// 		return false
// 	}

// 	if n != int(dataLength) {
// 		fmt.Printf("Only was able to read in %v bytes, not %v\n", n, dataLength)
// 		return false
// 	}

// 	stack.Push(data)
// 	return true
// }

func PrepareScript(script []byte) ([]Operation, error) {

	ops := make([]Operation, 0)
	reader := bytes.NewBuffer(script)

	for {
		op, err := NewOperation(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		ops = append(ops, op)
	}

	return ops, nil
}
