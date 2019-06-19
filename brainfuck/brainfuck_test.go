package brainfuck

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NewInterpreterTestSuite struct {
	suite.Suite
	testCode          string
	bracketsTestCases []bracketTestCases
}

func TestNewInterpreterTestSuite(t *testing.T) {
	suite.Run(t, new(NewInterpreterTestSuite))
}

type bracketTestCases struct {
	code           string
	errorMessage   string
	openBrackets   []int
	closedBrackets []int
}

func (suite *NewInterpreterTestSuite) SetupTest() {
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
		{
			code:           "[[[][[][]][[[]]]][[]][][]]",
			openBrackets:   []int{0, 1, 2, 4, 5, 7, 10, 11, 12, 17, 18, 21, 23},
			closedBrackets: []int{25, 16, 3, 9, 6, 8, 15, 14, 13, 20, 19, 22, 24},
		},
		{
			code:           "[+.[]]>-[.-]",
			openBrackets:   []int{0, 3, 8},
			closedBrackets: []int{5, 4, 11},
		},
		{code: "+[>,.<", errorMessage: notAllClosingErr},
		{code: "+>,.<]", errorMessage: notAllOpeningErr},
		{code: "+[[>,.<]", errorMessage: notAllClosingErr},
		{code: "+[>,.<]]", errorMessage: notAllOpeningErr},
		{code: "+[>[,.]<]--[++++", errorMessage: notAllClosingErr},
		{code: "+[>[,.]<]--]++++", errorMessage: notAllOpeningErr},
	}
}

func (suite *NewInterpreterTestSuite) TestNewInterpreter() {
	interpreter, err := NewInterpreter(suite.testCode)
	suite.Nil(err)

	suite.Equal(1, len(interpreter.closedBrackets))
	suite.Equal(1, len(interpreter.openBrackets))
	suite.Equal(suite.testCode, string(interpreter.code))
	suite.Equal(0, interpreter.pointer)
	suite.Greater(len(interpreter.data), 0)
}

func (suite *NewInterpreterTestSuite) TestNewInterpreterBracketsValidation() {
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

type OperatorTestSuite struct {
	suite.Suite
}

func TestOperatorSuite(t *testing.T) {
	suite.Run(t, new(OperatorTestSuite))
}

func (suite *OperatorTestSuite) TestPlus() {
	interpreter, err := NewInterpreter("+++++.")
	suite.Nil(err)

	var buffer bytes.Buffer
	resultWriter := bufio.NewWriter(&buffer)
	interpreter.Run(resultWriter, nil)

	resultWriter.Flush()
	output := buffer.Bytes()
	suite.EqualValues(5, output[0])
}

func (suite *OperatorTestSuite) TestMinus() {
	interpreter, err := NewInterpreter("-----.")
	suite.Nil(err)

	var buffer bytes.Buffer
	resultWriter := bufio.NewWriter(&buffer)
	interpreter.Run(resultWriter, nil)

	resultWriter.Flush()
	output := buffer.Bytes()
	suite.EqualValues(251, output[0])
}

func (suite *OperatorTestSuite) TestPlusAndMinus() {
	interpreter, err := NewInterpreter("+++++.-----.+++++---.+---.-----.++.")
	suite.Nil(err)

	var buffer bytes.Buffer
	resultWriter := bufio.NewWriter(&buffer)
	interpreter.Run(resultWriter, nil)

	resultWriter.Flush()
	output := buffer.Bytes()
	suite.EqualValues(5, output[0])
	suite.EqualValues(0, output[1])
	suite.EqualValues(2, output[2])
	suite.EqualValues(0, output[3])
	suite.EqualValues(251, output[4])
	suite.EqualValues(253, output[5])
}

func (suite *OperatorTestSuite) TestForward() {
	interpreter, err := NewInterpreter(">>>>>>>>>>")
	suite.Nil(err)

	interpreter.Run(nil, nil)

	suite.EqualValues(10, interpreter.pointer)
}

func (suite *OperatorTestSuite) TestBackward() {
	interpreter, err := NewInterpreter("<<<<<<<<<<")
	suite.Nil(err)

	interpreter.Run(nil, nil)

	suite.EqualValues(len(interpreter.data)-10, interpreter.pointer)
}

func (suite *OperatorTestSuite) TestForwardAndBackward() {
	interpreter, err := NewInterpreter(">+>++>+++>++++>+++++>++++++>+++++++>++++++++>+++++++++>++++++++++" + // setup cells with numbers from 0-10
		"<<<<<<<<<<.>>>>><<.>>>>>>>.<<<<<<<<<<.<<<<<>>>")
	suite.Nil(err)

	var buffer bytes.Buffer
	resultWriter := bufio.NewWriter(&buffer)
	interpreter.Run(resultWriter, nil)

	resultWriter.Flush()
	output := buffer.Bytes()

	suite.EqualValues(0, output[0])
	suite.EqualValues(3, output[1])
	suite.EqualValues(10, output[2])
	suite.EqualValues(0, output[3])
	suite.EqualValues(len(interpreter.data)-2, interpreter.pointer) //check last position
}

func (suite *OperatorTestSuite) TestInputOutput() {
	interpreter, err := NewInterpreter(",.>-,+.>,-..>,,.")
	suite.Nil(err)

	inputReader := bytes.NewReader([]byte("5a" + string(byte(1)) + "fl"))
	var buffer bytes.Buffer
	resultWriter := bufio.NewWriter(&buffer)
	interpreter.Run(resultWriter, inputReader)

	resultWriter.Flush()
	output := buffer.Bytes()
	suite.EqualValues(byte('5'), output[0])
	suite.EqualValues(byte('b'), output[1])
	suite.EqualValues(byte(0), output[2])
	suite.EqualValues(byte(0), output[3])
	suite.EqualValues(byte('l'), output[4])
}
