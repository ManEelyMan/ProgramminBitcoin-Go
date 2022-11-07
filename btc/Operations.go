package btc

import (
	"bitcoin-go/utility"
	"fmt"
	"io"
)

type Operation interface {
	GetOpCode() byte
	GetOpName() string
	Execute(context *ExecutionContext) bool
}

type GenericOperation struct {
	OpCode byte
	OpName string
	OpFxn  opFxn
}

func (op GenericOperation) GetOpCode() byte {
	return op.OpCode
}
func (op GenericOperation) GetOpName() string {
	return op.OpName
}
func (op GenericOperation) Execute(context *ExecutionContext) bool {
	return op.OpFxn(context)
}

type AddDataToStackOperation struct {
	OpCode byte
	OpName string
	Data   []byte
}

func (op AddDataToStackOperation) GetOpCode() byte {
	return op.OpCode
}
func (op AddDataToStackOperation) GetOpName() string {
	return op.OpName
}
func (op AddDataToStackOperation) Execute(context *ExecutionContext) bool {
	context.Stack.Push(op.Data)
	return true
}

func NewOperation(reader io.Reader) (Operation, error) {
	tmp := make([]byte, 1)
	_, err := reader.Read(tmp)
	if err != nil {
		return nil, err
	}

	opCode := tmp[0]
	opCodeName := opCodeNames[opCode]
	if len(opCodeName) == 0 {
		opCodeName = "[none]"
	}

	var data []byte = nil

	// First let's see if this is an "add data to stack" operation.
	if opCode == 0x00 {
		// This is an OP_0
		data = make([]byte, 0)

	} else if opCode >= 0x01 && opCode <= 0x4e {
		opCodeName = "[StackData]"
		// We have to decode the varint to see how much data to push on the stack.
		b, err := extractScriptData(opCode, reader)
		if err != nil {
			return nil, err
		}
		data = b
	} else if opCode == 0x49 {
		data = encodeNumber(-1)
	} else if opCode >= 0x51 && opCode <= 0x60 {
		num := opCode - 0x50
		data = encodeNumber(int64(num))
	}

	if data != nil {
		return AddDataToStackOperation{OpCode: opCode, OpName: opCodeName, Data: data}, nil
	}

	// If we made it this far, we don't have that.
	return GenericOperation{OpCode: opCode, OpName: opCodeNames[opCode], OpFxn: opCodeFxns[opCode]}, nil
}

func extractScriptData(op byte, reader io.Reader) ([]byte, error) {

	var dataLength uint32 = (uint32)(op)

	if op == 0x4c {
		dataLength = (uint32)(utility.ReadByte(reader))
	} else if op == 0x4d {
		dataLength = uint32(utility.ReadUInt16(reader, true))
	} else if op == 0x4e {
		dataLength = utility.ReadUint32(reader, true)
	}

	data := make([]byte, dataLength)
	n, err := reader.Read(data)

	if err != nil {
		return nil, err
	}

	if n != int(dataLength) {
		return nil, fmt.Errorf("only was able to read in %v bytes, not %v", n, dataLength)
	}

	return data, nil
}
