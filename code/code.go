package code

import (
	"bytes"
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

	// 피연산자를 부호화
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

// "\x00\x00\x00\x00\x00\x01" 이런 바이트 대신 사람이 읽을 수 있는 형태로 출력하기 위한 함수
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// 부호화된 피연산자를 복호화
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	// 피연산자를 모두 담을 수 있는 크기로 슬라이스 할당
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	// 가상 머신이 직접 호출할 수 있도록 공용 함수로 구현
	// 이를 통해 ReadOperands 사용 시 정의 탐색을 건너뛸 수 있음
	return binary.BigEndian.Uint16(ins)
}
