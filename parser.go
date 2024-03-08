package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// we want to tokenize the input

// func parseSimpleString(input []byte) RESP {
// 	// input will be of for "+[INPUT]/r/n", this function should return whatever [INPUT] is
// 	value := input[1 : len(input)-2]
// 	return RESP{Type: String, Raw: input, Data: value, Count: len(input)}
// }

// func parseSimpleError(input []byte) RESP {
// 	// input will be of for "-[INPUT]/r/n", this function should return whatever [INPUT] is
// 	value := input[1 : len(input)-2]
// 	return RESP{Type: Error, Raw: input, Data: value, Count: len(input)}
// }

type Parser struct {
	input    []byte
	position int
	ch       byte
}

// constructor for a parser
func Parse(bulkString []byte) ([]string, error) {
	p := Parser{
		input:    bulkString,
		position: 0,
	}
	if len(p.input) == 0 {
		return nil, errors.New("zero len string passed into parser")
	}
	p.ch = p.input[p.position]

	return p.parse()
}

func (p *Parser) readChar() {
	p.position++
	if p.position >= len(p.input) {
		return
	}
	p.ch = p.input[p.position]
}

func (p *Parser) currentAscii() {
	displayRepresentation := fmt.Sprintf("%q", p.ch)
	fmt.Println("current char: ", displayRepresentation)
}

func (p *Parser) isNewLine() bool {
	if p.ch != '\r' {
		return false
	}
	if p.position+1 >= len(p.input) {
		return false
	}
	if p.input[p.position+1] != '\n' {
		return false
	}
	return true
}

func (p *Parser) skipNewLine() error {
	if !p.isNewLine() {
		return errors.New("invalid newline")
	}
	p.readChar()
	p.readChar()
	return nil
}

func (p *Parser) parseInteger() (int, error) {

	var integerString strings.Builder

	for unicode.IsDigit(rune(p.ch)) {
		integerString.WriteByte(p.ch)
		p.readChar()
	}
	integer, error := strconv.Atoi(integerString.String())

	if error != nil {
		return 0, error
	}

	return integer, nil
}

func (p *Parser) parseBulkString() (string, error) {
	if p.ch != '$' {
		return "", errors.New("invalid bulk string")
	}

	p.readChar()

	bulkStringLength, error := p.parseInteger()

	if error != nil {
		return "", error
	}

	bulkStringStartingError := p.skipNewLine()

	if bulkStringStartingError != nil {
		return "", bulkStringStartingError
	}

	var string strings.Builder

	for i := 0; i < bulkStringLength; i++ {
		string.WriteByte(p.ch)
		p.readChar()
	}

	endOfStringError := p.skipNewLine()

	if endOfStringError != nil {
		return "", endOfStringError
	}

	return string.String(), nil
}

func (p *Parser) parseArray() ([]string, error) {

	if p.ch != '*' {
		return nil, errors.New("invalid array")
	}

	p.readChar()

	arrayLength, error := p.parseInteger()

	if error != nil {
		return nil, error
	}

	arrayStartingError := p.skipNewLine()

	if arrayStartingError != nil {
		return nil, arrayStartingError
	}

	array := make([]string, arrayLength)

	for i := 0; i < arrayLength; i++ {
		bulkString, error := p.parseBulkString()

		if error != nil {
			return nil, error
		}

		array[i] = bulkString
	}

	return array, nil
}

func (p *Parser) parse() ([]string, error) {
	if p.ch == '*' {
		elems, err := p.parseArray()
		if err != nil {
			return nil, err
		}
		return elems, nil
	}
	return nil, errors.New("unsupported data type")
}
