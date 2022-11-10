package btc

import (
	"bitcoin-go/collections"
	"bitcoin-go/utility"
	"bytes"
	"fmt"
	"math/big"
)

type ScriptExecutor struct {
	scriptPubKey    *Script
	scriptSignature *Script
	hash            *big.Int
}

func NewScriptExecutor(pubkey *Script, sig *Script, hash *big.Int) ScriptExecutor {
	return ScriptExecutor{scriptPubKey: pubkey, scriptSignature: sig, hash: hash}
}

func (ex *ScriptExecutor) Execute() bool {

	// Allocate new stacks for this execution run.
	stack := collections.NewStack()
	altStack := collections.NewStack()

	executionContext := ExecutionContext{Stack: &stack, AltStack: &altStack, Hash: ex.hash}

	// 1. Parse, load and execute script signature
	ok := executeScript(ex.scriptSignature, &executionContext)
	if !ok {
		return false
	}

	// If the pub key is a P2SH, the last thing on the stack is the redeem script. Save it for later.
	var redeemScriptBytes []byte = nil
	if ex.scriptPubKey.IsPayToScriptHash() {
		redeemScriptBytes, _ = stack.Peek()
	}

	// 2. Parse, load and execute script pub key
	ok = executeScript(ex.scriptPubKey, &executionContext)
	if !ok {
		return false
	}

	// 3. Verify completion
	if stack.Length() == 0 {
		return false
	}

	// 4. Ensure a non-zero value is here.
	if b, _ := stack.Pop(); len(b) == 0 {
		return false
	}

	// 5. Hack for P2SH, start executing the redeem script
	if ex.scriptPubKey.IsPayToScriptHash() {

		buffer := bytes.NewBuffer(make([]byte, 0))
		utility.WriteVarInt(buffer, uint64(len(redeemScriptBytes)))
		buffer.Write(redeemScriptBytes)
		newScript := ParseScript(buffer)

		ok := executeScript(&newScript, &executionContext)
		if !ok {
			return false
		}

		// Verify completion (again)
		if stack.Length() == 0 {
			return false
		}

		// Ensure a non-zero value is here. (again)
		if b, _ := stack.Pop(); len(b) == 0 {
			return false
		}
	}

	return true
}

func executeScript(script *Script, context *ExecutionContext) bool {

	operations := script.GetOperations()

	for i := 0; i < len(operations); i++ {
		op := operations[i]
		opCode := op.GetOpCode()

		// Make some decisions based on the op code!
		if opCode == 0x63 || opCode == 0x64 || opCode == 0x67 || opCode == 0x69 {
			// TODO: Figure out how to track the if/else/endif stuff.
			return false // NOT IMPLEMENTED
		} else {

			fmt.Printf("Performing op %+v (%v)...\n", op.GetOpName(), opCode)
			ok := op.Execute(context)
			if !ok {
				fmt.Printf("Failed processing op %+v.\n", op.GetOpName())
				return false
			}
		}
	}

	return true
}
