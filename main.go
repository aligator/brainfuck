package main

import (
	"fmt"
	"github.com/aligator/brainfuck/brainfuck"
	"os"
)

// a simple Writer to StdOut using fmt.Println
type StdWriter struct{}

func (StdWriter) Write(p []byte) (n int, err error) {
	return fmt.Print(string(p))
}

// a reader for readline which only reads one char and removes newlines
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
	fmt.Println("Hello World example:")
	fmt.Println()

	brfck := brainfuck.NewBrainfuckInterpreter(">>>++++++++++" +
		"[" +
		">+++++++>++++++++++>+++>+<<<<-" +
		"]   Schleife zur Vorbereitung der Textausgabe" +
		">++.                    Ausgabe von 'H'" +
		">+.                     Ausgabe von 'e'" +
		"+++++++.                'l'" +
		".                       'l'" +
		"+++.                    'o'" +
		">++.                    Leerzeichen" +
		"<<+++++++++++++++.      'W'" +
		">.                      'o'" +
		"+++.                    'r'" +
		"------.                 'l'" +
		"--------.               'd'" +
		">+.                     '!'" +
		">.                      Zeilenvorschub" +
		"+++.                    Wagenrücklauf")
	brfck.Run(StdWriter{}, StdCharReader{})
	fmt.Println()
	fmt.Println()

	fmt.Println("Echo example:")
	fmt.Println()

	// example which just echos out the input
	brfck = brainfuck.NewBrainfuckInterpreter("+[>,.<]")
	brfck.Run(StdWriter{}, StdCharReader{})
}
