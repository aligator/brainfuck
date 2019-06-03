package main

import (
	"fmt"
	"github.com/aligator/brainfuck/brainfuck"
	"os"
)

// a simple Writer to StdOut using fmt.Println
type StdWriter struct{}

func (StdWriter) Write(p []byte) (n int, err error) {
	return fmt.Println(string(p))
}

func main() {
	brfck := brainfuck.NewBrainfuckInterpreter("++++++++++" +
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
	brfck.Run(StdWriter{}, os.Stdin)
}
