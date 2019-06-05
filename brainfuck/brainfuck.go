package brainfuck

import (
	"io"
)

type Interpreter struct {
	code           []rune
	pointer        int
	openBrackets   []int
	closedBrackets []int
	data           []byte
}

func NewInterpreter(code string) *Interpreter {
	i := Interpreter{
		code:           []rune(code),
		pointer:        0,
		openBrackets:   nil,
		closedBrackets: nil,
		data:           make([]byte, 30000),
	}

	i.prepareCode()
	return &i
}

func (i *Interpreter) prepareCode() {
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

func (i *Interpreter) Run(w io.Writer, r io.Reader) {
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

func (i *Interpreter) incrementData() {
	i.setCurrentCell(i.getCurrentCell() + 1)
}

func (i *Interpreter) decrementData() {
	i.setCurrentCell(i.getCurrentCell() - 1)
}

func (i *Interpreter) incrementPointer() {
	i.pointer++
	if i.pointer >= len(i.data) {
		i.pointer = 0
	}
}

func (i *Interpreter) decrementPointer() {
	i.pointer--
	if i.pointer < 0 {
		i.pointer = len(i.data) - 1
	}
}

func (i *Interpreter) write(w io.Writer) {
	data := make([]byte, 1)
	data[0] = i.getCurrentCell()
	_, err := w.Write(data)
	if err != nil {
		panic(err)
	}
}

func (i *Interpreter) read(r io.Reader) {
	buff := make([]byte, 1)
	n, err := r.Read(buff)
	if err != nil {
		if err != io.EOF {
			panic(err)
		} else {
			// Do nothing
			// maybe add option to use 0 or \n (10) instead
			// buff[0] = 0
		}
	}

	if n > 0 {
		i.setCurrentCell(buff[0])
	}
}

func (i *Interpreter) getCurrentCell() byte {
	return i.data[i.pointer]
}

func (i *Interpreter) setCurrentCell(newVal byte) {
	i.data[i.pointer] = newVal
}
