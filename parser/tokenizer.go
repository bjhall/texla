package parser

import (
	"fmt"
	"strconv"
	"unicode"
)

type Token struct {
	kind   TokenKind
	str    string
	line   int
	column int
}

type Tokenizer struct {
	source        string
	pos           int
	currentLine   int
	currentColumn int
	state         string
}

func newTokenizer(source_string string) Tokenizer {
	parser := Tokenizer{
		source:        source_string,
		pos:           0,
		currentLine:   0,
		currentColumn: 0,
		state:         "top",
	}
	return parser
}


func (p *Tokenizer) currentRune() rune {
	return rune(p.source[p.pos])
}

func (p *Tokenizer) skip(n int) {
	end := p.pos + n
	for p.pos < end {
		if p.EOF() {
			break
		}
		if p.currentRune() == '\n' {
			p.currentColumn = 0
			p.currentLine += 1
		} else {
			p.currentColumn += 1
		}
		p.pos++
	}
}

func (p *Tokenizer) EOF() bool {
	return p.pos >= len(p.source)-1
}

func (p *Tokenizer) atLetter() bool {
	return unicode.IsLetter(p.currentRune())
}

func (p *Tokenizer) atNumber() bool {
	return unicode.IsNumber(p.currentRune())
}

func (p *Tokenizer) atWhitespace() bool {
	return unicode.IsSpace(p.currentRune())
}

func (p *Tokenizer) skipWhitespace() {
	for unicode.IsSpace(p.currentRune()) && !p.EOF() {
		p.skip(1)
	}
}

func (p *Tokenizer) consume() rune {
	r := p.currentRune()
	p.skip(1)
	return r
}

func (p *Tokenizer) consumeMany(n int) string {
	consumedString := ""
	start := p.pos
	for p.pos < start+n {
		consumedString += string(p.consume())
	}
	return consumedString
}

func (p *Tokenizer) consumeUntil(lastRune rune) string {
	consumedString := ""
	for p.currentRune() != lastRune && !p.EOF() {
		consumedString += string(p.consume())
	}
	return consumedString
}

func (p *Tokenizer) peek(ahead int) rune {
	return rune(p.source[p.pos+ahead])
}

func isAlNum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_'
}

func (t *Tokenizer) consumeIdentifier() Token {
	var identifierString string
	for isAlNum(t.currentRune()) && !t.EOF() {
		identifierString += string(t.consume())
	}
	switch identifierString {
	case "fn", "if", "for", "in", "read", "print", "return":
		return Token{kind: Keyword, str: identifierString}
	default:
		return Token{kind: Identifier, str: identifierString}
	}
}

func (t *Tokenizer) consumeNumber() (Token, error) {
	var numberString string
	dotCount := 0
	for (t.atNumber() || t.currentRune() == '.') && !t.EOF() {
		numberString += string(t.consume())
		if t.currentRune() == '.' {
			dotCount += 1
		}
	}

	switch dotCount {
	case 0:
		return t.createTokenFromString(Integer, numberString), nil
	case 1:
		return t.createTokenFromString(Float, numberString), nil
	default:
		return Token{}, fmt.Errorf("Invalid number: %s", numberString)
	}
}

func (t *Tokenizer) createTokenConsume(kind TokenKind, nchar int) Token {
	return Token{
		kind: kind,
		str: string(t.consumeMany(nchar)),
		line: t.currentLine,
		column: t.currentColumn,
	}
}

func (t *Tokenizer) createTokenFromString(kind TokenKind, str string) Token {
	return Token{
		kind: kind,
		str: str,
		line: t.currentLine,
		column: t.currentColumn,
	}
}



func (t *Tokenizer) nextToken() (Token, error) {

	if t.EOF() {
		return t.createTokenFromString(Eof, ""), nil
	}

	if t.atWhitespace() {
		t.skipWhitespace()
		return t.createTokenFromString(Whitespace, ""), nil
	}

	if t.atLetter() {
		token := t.consumeIdentifier()
		return token, nil
	}

	if t.atNumber() {
		token, err := t.consumeNumber()
		if err != nil {
			return Token{}, err
		}
		return token, nil
	}

	switch t.currentRune() {
	case '{':
		return t.createTokenConsume(OpenCurly, 1), nil
	case '}':
		return t.createTokenConsume(CloseCurly, 1), nil
	case '(':
		return t.createTokenConsume(OpenParen, 1), nil
	case ')':
		return t.createTokenConsume(CloseParen, 1), nil
	case '[':
		return t.createTokenConsume(OpenBracket, 1), nil
	case ']':
		return t.createTokenConsume(CloseBracket, 1), nil
	case '>':
		if t.peek(1) == '=' {
			return t.createTokenConsume(GreaterEqual, 2), nil
		}
		return t.createTokenConsume(Greater, 1), nil
	case '<':
		if t.peek(1) == '=' {
			return t.createTokenConsume(LessEqual, 2), nil
		}
		return t.createTokenConsume(Less, 1), nil
	case '-':
		if t.peek(1) == '>' {
			return t.createTokenConsume(RightArrow, 2), nil
		}
		return t.createTokenConsume(Minus, 1), nil
	case '+':
		return t.createTokenConsume(Plus, 1), nil
	case '*':
		return t.createTokenConsume(Mult, 1), nil
	case '/':
		if t.peek(1) == '/' {
			return t.createTokenFromString(Comment, t.consumeUntil('\n')), nil
		}
		return t.createTokenConsume(Div, 1), nil
	case ',':
		return t.createTokenConsume(Comma, 1), nil
	case '"':
		t.skip(1)
		stringLiteral := t.consumeUntil('"')
		t.skip(1)
		return t.createTokenFromString(StringLiteral, stringLiteral), nil
	case '=':
		if t.peek(1) == '=' {
			return t.createTokenConsume(Equal, 2), nil
		}
		return t.createTokenConsume(Assign, 1), nil
	case '!':
		if t.peek(1) == '=' {
			return t.createTokenConsume(NotEqual, 2), nil
		}
		return t.createTokenConsume(Not, 1), nil
	default:
		return Token{}, fmt.Errorf("Unknown token: %s", strconv.QuoteRune(t.currentRune()))
	}

	panic("Unsupported token")
}

func Tokenize(code_string string) ([]Token, error) {
	tokenizer := newTokenizer(code_string)

	var tokens []Token
	for {
		token, err := tokenizer.nextToken()
		if err != nil {
			return tokens, err
		}
		if token.kind != Whitespace && token.kind != Comment {
			tokens = append(tokens, token)
		}
		if token.kind == Eof {
			break
		}

	}

	return tokens, nil
}
