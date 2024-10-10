package lexer

import (
	"github.com/light-speak/lighthouse/errors"
)

type TokenType string

const (
	EOF         TokenType = "EOF"
	Letter      TokenType = "Letter"
	Boolean     TokenType = "Boolean"
	IntNumber   TokenType = "IntNumber"
	FloatNumber TokenType = "FloatNumber"
	Comment     TokenType = "Comment"
	Message     TokenType = "Message"

	Schema       TokenType = "Schema"
	Type         TokenType = "Type"
	Interface    TokenType = "Interface"
	Enum         TokenType = "Enum"
	Input        TokenType = "Input"
	Query        TokenType = "Query"
	Mutation     TokenType = "Mutation"
	Subscription TokenType = "Subscription"
	Extend       TokenType = "Extend"
	Implements   TokenType = "Implements"
	Scalar       TokenType = "Scalar"
	Union        TokenType = "Union"
	Directive    TokenType = "Directive"
	Fragment     TokenType = "Fragment"
	On           TokenType = "On"

	LeftBrace    TokenType = "{"
	RightBrace   TokenType = "}"
	LeftParent   TokenType = "("
	RightParent  TokenType = ")"
	LeftBracket  TokenType = "["
	RightBracket TokenType = "]"
	Colon        TokenType = ":"
	Comma        TokenType = ","
	Semicolon    TokenType = ";"
	Dot          TokenType = "."
	At           TokenType = "@"
	Hash         TokenType = "#"
	Pipe         TokenType = "|"
	DoubleQuote  TokenType = "\""
	SingleQuote  TokenType = "'"
	Backslash    TokenType = "\\"
	Exclamation  TokenType = "!"
	Equal        TokenType = "="
	And          TokenType = "&"
	Repeatable   TokenType = "repeatable"
	TripleDot    TokenType = "..."
)

var keywords = map[string]TokenType{
	"schema":       Schema,
	"type":         Type,
	"interface":    Interface,
	"enum":         Enum,
	"input":        Input,
	"Query":        Query,
	"Mutation":     Mutation,
	"Subscription": Subscription,
	"extend":       Extend,
	"implements":   Implements,
	"scalar":       Scalar,
	"union":        Union,
	"directive":    Directive,
	"on":           On,
	"fragment":     Fragment,
	"true":         Boolean,
	"false":        Boolean,
	"repeatable":   Repeatable,
	"...":          TripleDot,
}

type Token struct {
	Type         TokenType
	Value        string
	Line         int
	LinePosition int
}

type Content struct {
	Path    *string
	Content string
}

type Lexer struct {
	contents       []*Content
	currentContent *Content
	contentIndex   int
	position       int
	readPosition   int
	line           int
	linePosition   int
	// current character
	ch byte
	// whitespaceSet is a set of whitespace characters
	whitespaceSet map[byte]struct{}
	// specialSet is a set of special characters
	specialSet map[byte]struct{}

	currentToken  *Token
	previousToken *Token
}

// IsKeyword check if the word is a keyword
func (l *Lexer) IsKeyword(word string) bool {
	_, ok := keywords[word]
	return ok
}

// NewLexer create a new lexer
// and init specialSet and whitespaceSet
func NewLexer(contents []*Content) *Lexer {
	l := &Lexer{
		contents:     contents,
		contentIndex: 0,
		line:         1,
		whitespaceSet: map[byte]struct{}{
			' ': {}, '\t': {}, '\n': {}, '\r': {},
		},
		specialSet: map[byte]struct{}{
			'{': {}, '}': {}, '(': {}, ')': {}, '[': {}, ']': {},
			':': {}, ',': {}, ';': {}, '.': {}, '@': {}, '#': {},
			'|': {}, '"': {}, '\'': {}, '!': {}, '=': {}, '&': {},
		},
	}
	l.switchToNextContent()
	l.readChar()
	return l
}

// switchToNextContent switches to the next content file
func (l *Lexer) switchToNextContent() bool {
	if l.contentIndex < len(l.contents) {
		l.currentContent = l.contents[l.contentIndex]
		l.contentIndex++
		l.position = 0
		l.readPosition = 0
		l.line = 1
		l.linePosition = 0
		return true
	}
	return false
}

// readChar read next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.currentContent.Content) {
		if l.switchToNextContent() {
			l.ch = l.currentContent.Content[0]
			l.position = 0
			l.readPosition = 1
		} else {
			l.ch = 0
		}
	} else {
		l.ch = l.currentContent.Content[l.readPosition]
		l.position = l.readPosition
		l.readPosition++
	}
	l.linePosition++
}

// skipWhitespace skip whitespace
func (l *Lexer) skipWhitespace() {
	for _, ok := l.whitespaceSet[l.ch]; ok; _, ok = l.whitespaceSet[l.ch] {
		if l.ch == '\n' {
			l.line++
			l.linePosition = 0
		}
		l.readChar()
	}
}

// isSpecialChar check if the character is a special character
// for example: {, }, (, ), [, ], :, ,
func (l *Lexer) isSpecialChar(ch byte) bool {
	_, ok := l.specialSet[ch]
	return ok
}

// NextToken get next token
func (l *Lexer) NextToken() (token *Token, err error) {
	l.skipWhitespace()

	switch {
	case l.isSpecialChar(l.ch):
		token = l.handleSpecialChar()
	case isLetter(l.ch):
		token = l.handleLetter()
	case isDigit(l.ch):
		token = l.handleNumber()
	case l.ch == 0:
		token = &Token{Type: EOF, Line: l.line, LinePosition: l.linePosition}
	default:
		token, err = l.handleUnrecognized()
		if err != nil {
			return nil, err
		}
	}

	l.previousToken = l.currentToken
	l.currentToken = token
	return token, nil
}

// handleSpecialChar handle special character
// if current character is # , read comment
// if current character is " , read message
// else return special character
func (l *Lexer) handleSpecialChar() *Token {
	switch l.ch {
	case '#':
		return l.readComment()
	case '"':
		return l.readMessage()
	default:
		tok := &Token{
			Type:         TokenType(l.ch),
			Value:        string(l.ch),
			Line:         l.line,
			LinePosition: l.linePosition,
		}
		l.readChar()
		return tok
	}
}

// readComment read comment
// for example: # this is a comment
func (l *Lexer) readComment() *Token {
	start := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return &Token{
		Type:         Comment,
		Value:        l.currentContent.Content[start:l.position],
		Line:         l.line,
		LinePosition: l.linePosition,
	}
}

// readMessage read message
// for example: "hello", "world", "123", "1.23", "1.23e-10"
func (l *Lexer) readMessage() *Token {
	start := l.position
	l.readChar()
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar()
	}

	tokenType := Message
	if l.currentToken != nil && l.currentToken.Type == Colon {
		tokenType = Letter
	}

	return &Token{
		Type:         tokenType,
		Value:        l.currentContent.Content[start:l.position],
		Line:         l.line,
		LinePosition: l.linePosition,
	}
}

// handleLetter handle letter and bool
// for example: id, name, age, email, createdAt, role
func (l *Lexer) handleLetter() *Token {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	word := l.currentContent.Content[start:l.position]
	tokType, ok := keywords[word]
	if !ok {
		tokType = Letter
	}
	return &Token{
		Type:         tokType,
		Value:        word,
		Line:         l.line,
		LinePosition: l.linePosition,
	}
}

// handleNumber handle number
// for example: 123, 1.23, 1.23e-10
func (l *Lexer) handleNumber() *Token {
	start := l.position
	isFloat := false

	for isDigit(l.ch) || l.ch == '.' || l.ch == 'e' || l.ch == 'E' || l.ch == '-' || l.ch == '+' {
		if l.ch == '.' || l.ch == 'e' || l.ch == 'E' {
			isFloat = true
		}
		l.readChar()
	}

	tokenType := IntNumber
	if isFloat {
		tokenType = FloatNumber
	}

	return &Token{
		Type:         tokenType,
		Value:        l.currentContent.Content[start:l.position],
		Line:         l.line,
		LinePosition: l.linePosition,
	}
}

// handleUnrecognized handle unrecognized character
// for example: %, ^, &, *
func (l *Lexer) handleUnrecognized() (*Token, error) {
	return nil, &errors.LexerError{
		Path:         l.currentContent.Path,
		Line:         l.line,
		LinePosition: l.linePosition,
		Message:      "unrecognized character: " + string(l.ch),
	}
}

// isLetter check if the character is a letter
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// isDigit check if the character is a digit
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// PreviousToken returns the previous token without moving the position
func (l *Lexer) PreviousToken() *Token {
	if l.previousToken != nil {
		return l.previousToken
	}
	return &Token{Type: EOF}
}
