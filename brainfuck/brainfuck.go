package brainfuck

import (
	"io"
)

type BrainfuckInterpreter struct {
	code           []rune
	pointer        int
	openBrackets   []int
	closedBrackets []int
	input          rune
	data           []rune
}

func NewBrainfuckInterpreter(code string) *BrainfuckInterpreter {
	i := BrainfuckInterpreter{
		code:           []rune(code),
		pointer:        0,
		openBrackets:   nil,
		closedBrackets: nil,
		input:          ' ',
		data:           make([]rune, 30000),
	}

	i.prepareCode()
	return &i
}

func (i *BrainfuckInterpreter) prepareCode() {
	counter := 0

	// count occurences to initialize slice in correct size
	for _, cmd := range i.code {
		if cmd == '[' {
			counter++
		}
	}

	i.openBrackets = make([]int, counter)
	i.closedBrackets = make([]int, counter)

	counter = 0
	for pos, cmd := range i.code {
		if cmd == '[' {
			i.openBrackets[counter] = pos
			loops := 1
			loopPointer := pos

			for loops > 0 {
				loopPointer++
				if i.code[loopPointer] == '[' {
					loops++
				}
				if i.code[loopPointer] == ']' {
					loops--
				}
			}
			i.closedBrackets[counter] = loopPointer
			counter++
		}
	}
}

func (i *BrainfuckInterpreter) Run(w io.Writer, r io.Reader) {
	codePointer := 0
	for codePointer < len(i.code) {
		switch i.code[codePointer] {
		case '>':
			i.incrementPointer()
		case '<':
			i.decrementPointer()
		case '+':
			i.incrementData()
		case '-':
			i.decrementData()
		case '.':
			i.write(w)
		case ',':
			i.read(r)
		case '[':
			if i.getCurrentCell() == 0 {
				// count brackets except current pos
				j := 0
				for i.openBrackets[j] != codePointer {
					j++
				}
				// jump to closed
				codePointer = i.closedBrackets[j]
			}
		case ']':
			// count brackets except current pos
			j := 0
			for i.closedBrackets[j] != codePointer {
				j++
			}
			// jump back
			codePointer = i.openBrackets[j] - 1
		}

		codePointer++
	}
}

func (i *BrainfuckInterpreter) incrementData() {
	i.setCurrentCell(i.getCurrentCell() + 1)
}

func (i *BrainfuckInterpreter) decrementData() {
	i.setCurrentCell(i.getCurrentCell() - 1)
}

func (i *BrainfuckInterpreter) incrementPointer() {
	i.pointer++
	if i.pointer >= len(i.data) {
		i.pointer = 0
	}
}

func (i *BrainfuckInterpreter) decrementPointer() {
	i.pointer--
	if i.pointer < 0 {
		i.pointer = len(i.data) - 1
	}
}

func (i *BrainfuckInterpreter) write(w io.Writer) {
	_, err := w.Write([]byte(string(i.getCurrentCell())))
	if err != nil {
		panic(err)
	}
}

func (i *BrainfuckInterpreter) read(r io.Reader) {
	buff := make([]byte, 1)
	n, err := r.Read(buff)
	if err != nil {
		panic(err)
	}

	if n > 0 {
		i.setCurrentCell(rune(buff[0]))
	}
}

func (i *BrainfuckInterpreter) getCurrentCell() rune {
	return i.data[i.pointer]
}

func (i *BrainfuckInterpreter) setCurrentCell(newVal rune) {
	i.data[i.pointer] = newVal
}