package main

import (
	"flag"
	"fmt"
	"github.com/aligator/brainfuck/brainfuck"
	"io/ioutil"
	"os"
)

// a simple Writer to StdOut using fmt.Println
type StdWriter struct{}

func (StdWriter) Write(p []byte) (n int, err error) {
	return fmt.Print(string(p))
}

// a reader for stdin which only reads one char and removes newlines
type StdCharReader struct{}

func (StdCharReader) Read(p []byte) (n int, err error) {
	// shrink buffer to only one byte, as we only need one
	p = p[0:1]
	reader := os.Stdin

	n, err = reader.Read(p)
	if err != nil {
		return
	}

	if n > 0 {
		p = p[0:1]
		// if first char is not newline: remove the next char as should be a newline
		// if first char is already a newline, it is an intended user input
		if int(p[0]) != 10 && int(p[0]) != 13 {
			// remove newline from stdio-buffer
			tmpBuff := make([]byte, 1)
			reader.Read(tmpBuff)
		}
	}

	return
}

func main() {
	filePtr := flag.String("file", "", "a filename")

	flag.Parse()

	file, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		panic(err)
	}

	brfck, err := brainfuck.NewInterpreter(string(file))
	if err != nil {
		panic(err)
	}
	brfck.Run(StdWriter{}, StdCharReader{})
}
