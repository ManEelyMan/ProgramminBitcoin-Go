package transaction

import (
	"bitcoin-go/utility"
	"bytes"
	"io"
)

type Script struct {
	RawData    []byte
	operations []Operation
}

func ParseScript(reader io.Reader) Script {
	scriptLength := utility.ReadVarInt(reader)

	s := Script{}
	s.RawData = make([]byte, scriptLength)
	reader.Read(s.RawData)
	return s
}

func (script *Script) Serialize(writer io.Writer) {
	utility.WriteVarInt(writer, (uint64)(len(script.RawData)))
	writer.Write(script.RawData)
}

func (script *Script) parseOperations() ([]Operation, error) {
	ops := make([]Operation, 0)
	reader := bytes.NewBuffer(script.RawData)

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

func (script *Script) GetOperations() []Operation {
	if script.operations == nil {
		ops, err := script.parseOperations()
		if err != nil {
			panic("AHHHH!")
		}

		script.operations = ops
	}
	return script.operations
}

func (script *Script) IsPayToScriptHash() bool {

	ops := script.GetOperations()

	if len(ops) != 3 {
		return false
	}

	if ops[0].GetOpCode() != 0xa9 || ops[2].GetOpCode() != 0x87 {
		return false
	}

	dataOp, ok := ops[1].(AddDataToStackOperation)
	if !ok || len(dataOp.Data) != 20 {
		return false
	}

	return true
}

func (script *Script) GetRedeemScriptHash() *Script {
	if !script.IsPayToScriptHash() {
		return nil
	}

	dataOp, ok := script.operations[1].(AddDataToStackOperation)
	if !ok || len(dataOp.Data) != 20 {
		return nil
	}

	reader := bytes.NewBuffer(dataOp.Data)
	s := ParseScript(reader)
	return &s
}
