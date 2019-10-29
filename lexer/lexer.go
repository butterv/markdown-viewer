package lexer

import (
	"github.com/istsh/markdown-viewer/token"
)

type Lexer struct {
	input        string // 入力
	position     int    // 入力における現在の位置(現在の文字を指し示す)
	readPosition int    // これから読み込む位置(現在の文字の次)
	ch           byte   // 現在検査中の文字
}

func (l *Lexer) GetInput() string {
	return l.input
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}

	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// 空白を飛ばす
	l.skipWhitespace()

	// TODO: 行の先頭がタブではなくスペースになっている場合も考慮する？

	switch l.ch {
	case '#':
		literal, cnt := l.readHeading()
		tok.Literal = literal
		tok.Type = token.GetHeadingToken(cnt)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '\t':
		literal, cnt := l.readTab()
		tok.Literal = literal
		tok.Type = token.GetTabToken(cnt)
		// タブの直後はスペースが入らないので、最後のreadCharを実行させたくない
		return tok
	case '':
		// TODO: 数値+ドット+空白のセットで判定
	//	// TODO: to 3chars
	//	tok = newToken(token.MINUS, l.ch)
	//case '>':
	//	tok = newToken(token.GT, l.ch)
	//case '.':
	//	tok = newToken(token.DOT, l.ch)
	//case '*':
	//	// TODO: to 3chars
	//	tok = newToken(token.ASTERISK, l.ch)
	//case '_':
	//	// TODO: to 3chars
	//	tok = newToken(token.UNDER_SCORE, l.ch)
	//case '+':
	//	tok = newToken(token.PLUS, l.ch)
	//
	//case '(':
	//	tok = newToken(token.LPAREN, l.ch)
	//case ')':
	//	tok = newToken(token.RPAREN, l.ch)
	//case '[':
	//	tok = newToken(token.LBRACKET, l.ch)
	//case ']':
	//	tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		tok.Literal = l.readString()
		tok.Type = token.STRING
	}

	// 1文字進める
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	// Tokenオブジェクトを初期化する
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) readChar() {
	// 次の文字が存在するか
	if l.readPosition >= len(l.input) {
		// 次の文字は存在しない(まだ何も読み込んでいない or ファイルの終わり)
		l.ch = 0
	} else {
		// 次の文字をセット
		l.ch = l.input[l.readPosition]
	}
	// 数値を1つ進める
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	// 次の文字を覗き見る
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	// TODO: タブはインデントを判定する為に使うのでスキップしない
	//for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
	//	// 空白、タブ、改行のときに飛ばして1文字進める
	//	l.readChar()
	//}
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\r' {
		// 空白や改行を飛ばして1文字進める
		l.readChar()
	}
}

func (l *Lexer) isHeading() bool {
	return l.ch == '#'
}

func (l *Lexer) readHeading() (string, int) {
	position := l.position

	var cnt int
	for l.isHeading() {
		// 文字が途切れるまで読み込む
		l.readChar()
		cnt++
	}
	// positionから、readCharで進んだところまで抽出
	return l.input[position:l.position], cnt
}

func (l *Lexer) readString() string {
	position := l.position
	for !l.isLineFeedCode() {
		// 文字が途切れるまで読み込む
		l.readChar()
	}
	// positionから、readCharで進んだところまで抽出
	return l.input[position:l.position]
}

func (l *Lexer) isLineFeedCode() bool {
	return l.ch == '\n' || l.ch == '\r'
}

func (l *Lexer) isTab() bool {
	return l.ch == '\t'
}

func (l *Lexer) readTab() (string, int) {
	position := l.position

	var cnt int
	for l.isTab() {
		// 文字が途切れるまで読み込む
		l.readChar()
		cnt++
	}
	// positionから、readCharで進んだところまで抽出
	return l.input[position:l.position], cnt
}
