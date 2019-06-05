package brainfuck

import (
	"github.com/pkg/errors"
	"io"
)

type Interpreter struct {
	code           string
	pointer        int
	openBrackets   []int
	closedBrackets []int
	data           []byte
}

func NewInterpreter(code string) (*Interpreter, error) {
	i := Interpreter{
		code:           code,
		pointer:        0,
		openBrackets:   nil,
		closedBrackets: nil,
		data:           make([]byte, 30000),
	}

	err := i.prepareCode()
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (i *Interpreter) prepareCode() error {
	const noClosingBracket = "not all opening brackets have a closing bracket"
	const noOpeningBracket = "there is a closing bracket before an opening bracket"

	i.openBrackets = []int{}
	i.closedBrackets = []int{}

	var bracketStack []int

	for codePos, cmd := range i.code {
		if cmd == '[' {
			bracketStack = append(bracketStack, len(i.openBrackets))
			i.openBrackets = append(i.openBrackets, codePos)
			i.closedBrackets = append(i.closedBrackets, -1) // -1: dummy
		} else if cmd == ']' {
			if len(bracketStack) == 0 {
				return errors.New(noOpeningBracket)
			}

			lastOpen := bracketStack[len(bracketStack)-1]
			i.closedBrackets[lastOpen] = codePos
			bracketStack = bracketStack[:len(bracketStack)-1]
		}
	}

	if len(bracketStack) != 0 {
		return errors.New(noClosingBracket)
	}

	return nil
}

func (i *Interpreter) Run(w io.Writer, r io.Reader) {
	codePointer := 0
	for codePointer < len(i.code) {
		//fmt.Println(codePointer, "----- ", string(i.code[codePointer]), "data:", i.getCurrentCell(), string(i.getCurrentCell()), "pointer", i.pointer)
		//time.Sleep(time.Second / 20)
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
				// count brackets until current pos
				j := 0
				for i.openBrackets[j] != codePointer {
					j++
				}
				// jump to closed
				codePointer = i.closedBrackets[j]
			}
		case ']':
			// count brackets until current pos
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
