package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

// 컴파일 시 상수를 만나면 이를 따로 저장해놓는다.
// 가상 머신은 이를 참조에 값을 가져올 수도 있다.
const OpConstant Opcode = iota

// 편리한 디버깅을 위한 tooling.
type Definition struct {
	Name          string // 사람이 읽을 수 있는 이름
	OperandWidths []int  // 피연산자가 차지하는 바이트 크기
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2} /* 피연산자는 2바이트, 즉 16비트를 갖는다. */},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}
