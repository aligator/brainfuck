package brainfuck

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type NewInterpreterTestSuit struct {
	suite.Suite
	testCode          string
	bracketsTestCases []bracketTestCases
}

func TestNewInterpreterTestSuit(t *testing.T) {
	suite.Run(t, new(NewInterpreterTestSuit))
}

type bracketTestCases struct {
	code           string
	errorMessage   string
	openBrackets   []int
	closedBrackets []int
}

func (suite *NewInterpreterTestSuit) SetupTest() {
	const notAllClosingErr = "not all opening brackets have a closing bracket"
	const notAllOpeningErr = "there is a closing bracket before an opening bracket"

	suite.testCode = "+[>,.<]"
	suite.bracketsTestCases = []bracketTestCases{
		{
			code:           "+[>,.<]",
			openBrackets:   []int{1},
			closedBrackets: []int{6},
		},
		{
			code:           "[+[>,[.]<] [++++] []]",
			openBrackets:   []int{0, 2, 5, 11, 18},
			closedBrackets: []int{20, 9, 7, 16, 19},
		},
		{code: "+[>,.<", errorMessage: notAllClosingErr},
		{code: "+>,.<]", errorMessage: notAllOpeningErr},
		{code: "+[[>,.<]", errorMessage: notAllClosingErr},
		{code: "+[>,.<]]", errorMessage: notAllOpeningErr},
		{code: "+[>[,.]<]--[++++", errorMessage: notAllClosingErr},
		{code: "+[>[,.]<]--]++++", errorMessage: notAllOpeningErr},
	}
}

func (suite *NewInterpreterTestSuit) TestNewInterpreter() {
	interpreter, err := NewInterpreter(suite.testCode)
	suite.Nil(err)

	suite.Equal(1, len(interpreter.closedBrackets))
	suite.Equal(1, len(interpreter.closedBrackets))
	suite.Equal(suite.testCode, string(interpreter.code))
	suite.Equal(0, interpreter.pointer)
	suite.Greater(len(interpreter.data), 0)
}

func (suite *NewInterpreterTestSuit) TestNewInterpreterBracketsValidation() {
	for _, testCase := range suite.bracketsTestCases {
		interpreter, err := NewInterpreter(testCase.code)
		if testCase.errorMessage != "" {
			suite.EqualError(err, testCase.errorMessage, testCase)
		} else {
			suite.NoError(err, testCase.errorMessage, testCase)
			suite.EqualValues(testCase.openBrackets, interpreter.openBrackets, testCase)
			suite.EqualValues(testCase.closedBrackets, interpreter.closedBrackets, testCase)
		}
	}
}
