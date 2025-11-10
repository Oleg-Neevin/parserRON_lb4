package parserron

import (
	"fmt"
	"strings"
	"unicode"

	s "gitlab.se.ifmo.ru/s503298/inf_lab_4/pkg"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdent
	TokenString
	TokenLParen
	TokenRParen
	TokenLBrack
	TokenRBrack
	TokenColon
	TokenComma
)

type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenIdent:
		return "IDENT"
	case TokenString:
		return "STRING"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenLBrack:
		return "["
	case TokenRBrack:
		return "]"
	case TokenColon:
		return ":"
	case TokenComma:
		return ","
	default:
		return "UNKNOWN"
	}
}

type Tokenizer struct {
	input string
	pos   int
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input: input,
		pos:   0,
	}
}

func (t *Tokenizer) peek() byte {
	if t.pos >= len(t.input) {
		return 0
	}
	return t.input[t.pos]
}

func (t *Tokenizer) advance() byte {
	if t.pos >= len(t.input) {
		return 0
	}
	ch := t.input[t.pos]
	t.pos++
	return ch
}

func (t *Tokenizer) skipWhitespace() {
	for t.pos < len(t.input) {
		ch := t.input[t.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			t.pos++
		} else {
			break
		}
	}
}

func (t *Tokenizer) skipComments() {
	for {
		t.skipWhitespace()

		if t.pos >= len(t.input) {
			break
		}

		if t.pos+1 < len(t.input) && t.input[t.pos:t.pos+2] == "//" {
			for t.pos < len(t.input) && t.input[t.pos] != '\n' {
				t.pos++
			}
			continue
		}

		if t.pos+1 < len(t.input) && t.input[t.pos:t.pos+2] == "/*" {
			t.pos += 2
			for t.pos+1 < len(t.input) {
				if t.input[t.pos:t.pos+2] == "*/" {
					t.pos += 2
					break
				}
				t.pos++
			}
			continue
		}

		break
	}
}

func (t *Tokenizer) readString() (string, error) {
	if t.peek() != '"' {
		return "", fmt.Errorf("expected '\"' at position %d", t.pos)
	}

	t.advance()
	var result strings.Builder

	for t.pos < len(t.input) {
		ch := t.peek()

		if ch == '\\' {
			t.advance()
			if t.pos >= len(t.input) {
				return "", fmt.Errorf("unexpected end after backslash")
			}
			next := t.advance()
			switch next {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			default:
				result.WriteByte(next)
			}
			continue
		}

		if ch == '"' {
			t.advance()
			return result.String(), nil
		}
		result.WriteByte(ch)
		t.advance()
	}

	return "", fmt.Errorf("unclosed string starting at position %d", t.pos)
}

func (t *Tokenizer) readIdent() string {
	start := t.pos
	for t.pos < len(t.input) {
		ch := rune(t.input[t.pos])
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			t.pos++
		} else {
			break
		}
	}
	return t.input[start:t.pos]
}

func (t *Tokenizer) NextToken() (Token, error) {
	t.skipComments()

	if t.pos >= len(t.input) {
		return Token{Type: TokenEOF, Pos: t.pos}, nil
	}

	ch := t.peek()
	startPos := t.pos

	if ch == '"' {
		str, err := t.readString()
		if err != nil {
			return Token{}, err
		}
		return Token{Type: TokenString, Value: str, Pos: startPos}, nil
	}

	if unicode.IsLetter(rune(ch)) || ch == '_' {
		ident := t.readIdent()
		return Token{Type: TokenIdent, Value: ident, Pos: startPos}, nil
	}

	t.advance()
	switch ch {
	case '(':
		return Token{Type: TokenLParen, Value: "(", Pos: startPos}, nil
	case ')':
		return Token{Type: TokenRParen, Value: ")", Pos: startPos}, nil
	case '[':
		return Token{Type: TokenLBrack, Value: "[", Pos: startPos}, nil
	case ']':
		return Token{Type: TokenRBrack, Value: "]", Pos: startPos}, nil
	case ':':
		return Token{Type: TokenColon, Value: ":", Pos: startPos}, nil
	case ',':
		return Token{Type: TokenComma, Value: ",", Pos: startPos}, nil
	default:
		return Token{}, fmt.Errorf("unexpected character '%c' at position %d", ch, startPos)
	}
}

func (t *Tokenizer) Tokenize() ([]Token, error) {
	var tokens []Token
	for {
		token, err := t.NextToken()
		if err != nil {
			return nil, err
		}
		if token.Type == TokenEOF {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *Parser) expect(typ TokenType) (Token, error) {
	tok := p.advance()
	if tok.Type != typ {
		return tok, fmt.Errorf("expected %s, got %s at position %d", typ, tok.Type, tok.Pos)
	}
	return tok, nil
}

func (p *Parser) parseValue() (interface{}, error) {
	tok := p.peek()

	switch tok.Type {
	case TokenString:
		p.advance()
		return tok.Value, nil

	case TokenIdent:
		typeName := p.advance().Value
		_, err := p.expect(TokenLParen)
		if err != nil {
			return nil, err
		}
		fields, err := p.parseFields()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(TokenRParen)
		if err != nil {
			return nil, err
		}

		return p.convertToStruct(typeName, fields)

	case TokenLBrack:
		return p.parseArray()

	default:
		return nil, fmt.Errorf("unexpected token %s at position %d", tok.Type, tok.Pos)
	}
}

func (p *Parser) parseArray() ([]interface{}, error) {
	_, err := p.expect(TokenLBrack)
	if err != nil {
		return nil, err
	}

	var elements []interface{}

	if p.peek().Type == TokenRBrack {
		p.advance()
		return elements, nil
	}

	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		elements = append(elements, value)

		tok := p.peek()
		if tok.Type == TokenComma {
			p.advance()
			if p.peek().Type == TokenRBrack {
				break
			}
		} else if tok.Type == TokenRBrack {
			break
		} else {
			return nil, fmt.Errorf("expected ',' or ']', got %s at position %d", tok.Type, tok.Pos)
		}
	}

	_, err = p.expect(TokenRBrack)
	if err != nil {
		return nil, err
	}

	return elements, nil
}

func (p *Parser) parseFields() (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	if p.peek().Type == TokenRParen {
		return fields, nil
	}

	for {
		tok := p.peek()
		if tok.Type != TokenIdent {
			return nil, fmt.Errorf("expected field name, got %s at position %d", tok.Type, tok.Pos)
		}
		fieldName := p.advance().Value

		_, err := p.expect(TokenColon)
		if err != nil {
			return nil, err
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		fields[fieldName] = value

		tok = p.peek()
		if tok.Type == TokenComma {
			p.advance()
			if p.peek().Type == TokenRParen {
				break
			}
		} else if tok.Type == TokenRParen {
			break
		} else {
			return nil, fmt.Errorf("expected ',' or ')', got %s at position %d", tok.Type, tok.Pos)
		}
	}

	return fields, nil
}

func (p *Parser) convertToStruct(typeName string, fields map[string]interface{}) (interface{}, error) {
	switch typeName {
	case "Schedule":
		schedule := &s.Schedule{}
		if daysRaw, ok := fields["days"].([]interface{}); ok {
			for _, dayRaw := range daysRaw {
				if day, ok := dayRaw.(*s.Day); ok {
					schedule.Days = append(schedule.Days, *day)
				} else {
					return nil, fmt.Errorf("expected Day in days array, got %T", dayRaw)
				}
			}
		}
		return schedule, nil

	case "Day":
		day := &s.Day{}
		if name, ok := fields["name"].(string); ok {
			day.Name = name
		}
		if lessonsRaw, ok := fields["lessons"].([]interface{}); ok {
			for _, lessonRaw := range lessonsRaw {
				if lesson, ok := lessonRaw.(*s.Lesson); ok {
					day.Lessons = append(day.Lessons, *lesson)
				} else {
					return nil, fmt.Errorf("expected Lesson in lessons array, got %T", lessonRaw)
				}
			}
		}
		return day, nil

	case "Lesson":
		lesson := &s.Lesson{}
		if v, ok := fields["time"].(string); ok {
			lesson.Time = v
		}
		if v, ok := fields["subject"].(string); ok {
			lesson.Subject = v
		}
		if v, ok := fields["teacher"].(string); ok {
			lesson.Teacher = v
		}
		if v, ok := fields["room"].(string); ok {
			lesson.Room = v
		}
		if v, ok := fields["building"].(string); ok {
			lesson.Building = v
		}
		if v, ok := fields["type"].(string); ok {
			lesson.Type = v
		}
		return lesson, nil

	default:
		return nil, fmt.Errorf("unknown type: %s", typeName)
	}
}

func (p *Parser) Parse() (*s.Schedule, error) {
	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	schedule, ok := value.(*s.Schedule)
	if !ok {
		return nil, fmt.Errorf("expected Schedule at root, got %T", value)
	}

	return schedule, nil
}

func ParseRON(input string) (*s.Schedule, error) {
	tokenizer := NewTokenizer(input)
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("tokenization error: %w", err)
	}

	parser := NewParser(tokens)
	schedule, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return schedule, nil
}
