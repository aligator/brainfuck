package main

type BrainfuckInterpreter struct {
	code           []rune
	pointer        uint
	openBrackets   []uint
	closedBrackets []uint
	input          rune
	data           []rune
}

func NewBrainfuckInterpreter(code string) *BrainfuckInterpreter {
	return &BrainfuckInterpreter{
		code:           []rune(code),
		pointer:        0,
		openBrackets:   make([]uint, 0),
		closedBrackets: make([]uint, 0),
		data:           make([]rune, 30000),
	}
}

func (i *BrainfuckInterpreter) Run() {

}

func (i *BrainfuckInterpreter) incrementData() {

}

func (i *BrainfuckInterpreter) decrementData() {

}

func (i *BrainfuckInterpreter) incrementPointer() {

}

func (i *BrainfuckInterpreter) decrementPointer() {

}

func (i *BrainfuckInterpreter) write() {

}
