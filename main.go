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
	reader := os.Stdin

	for {
		n, err = reader.Read(p)

		if err != nil {
			break
		}

		if n > 0 {
			p = p[0:1]

			if int(p[0]) != 10 && int(p[0]) != 13 {
				break
			}
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

	brfck := brainfuck.NewBrainfuckInterpreter(string(file))
	brfck.Run(StdWriter{}, StdCharReader{})
}
