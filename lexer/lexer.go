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

	//l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	// 1文字進める
	l.readChar()

	var tok token.Token

	switch l.ch {
	case '#':
		literal, cnt := l.readHeading()
		tok.Literal = literal
		tok.Type = token.GetHeadingToken(cnt)
	case '-':
		tok = newToken(token.HYPHEN, l.ch)
	case '\t':
		literal, cnt := l.readTab()
		tok.Literal = literal
		tok.Type = token.GetTabToken(cnt)
	case ' ':
		//if l.lookBackChar() == '\n' {
		//	literal, cnt := l.readSpace()
		//	tok.Literal = literal
		//	tok.Type = token.GetSpaceToken(cnt)
		//} else {
		//	tok = newToken(token.SPACE1, l.ch)
		//}
		tok = newToken(token.SPACE, l.ch)
	case '\n':
		tok = newToken(token.LINE_FEED_CODE_N, l.ch)
	case '\r':
		tok = newToken(token.LINE_FEED_CODE_R, l.ch)
	case '>':
		if l.lookBackChar() == '\n' {
			tok = newToken(token.CITATION, l.ch)
		}
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

func (l *Lexer) peekNextChar() byte {
	// 次の文字を覗き見る
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) peek2ndOrderChar() byte {
	// 次の次の文字を覗き見る
	if l.readPosition+1 >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition+1]
	}
}

func (l *Lexer) lookBackChar() byte {
	// 直前の文字を見る
	if l.readPosition < 2 {
		return 0
	}
	return l.input[l.readPosition-2]
}

func isHeading(ch byte) bool {
	return ch == '#'
}

func (l *Lexer) readHeading() (string, int) {
	position := l.position

	cnt := 1
	for isHeading(l.ch) && isHeading(l.peekNextChar()) {
		// 文字が途切れるまで読み込む
		l.readChar()
		cnt++
	}
	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1], cnt
}

func (l *Lexer) readString() string {
	position := l.position
	for !isLineFeedCode(l.ch) && !isLineFeedCode(l.peekNextChar()) {
		// 文字が途切れるまで読み込む
		l.readChar()
	}
	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1]
}

func isLineFeedCode(ch byte) bool {
	return ch == '\n' || ch == '\r'
}

func isTab(ch byte) bool {
	return ch == '\t'
}

func (l *Lexer) readTab() (string, int) {
	position := l.position

	cnt := 1
	for isTab(l.ch) && isTab(l.peekNextChar()) {
		// 文字が途切れるまで読み込む
		l.readChar()
		cnt++
	}
	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1], cnt
}

func isSpace(ch byte) bool {
	return ch == ' '
}

func (l *Lexer) readSpace() (string, int) {
	position := l.position

	cnt := 1
	for isSpace(l.ch) && isSpace(l.peekNextChar()) && cnt <= 4 {
		// 文字が途切れるまで読み込む
		l.readChar()
		cnt++
	}
	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1], cnt
}

func isDigit(ch byte) bool {
	// 数値
	return '0' <= ch && ch <= '9'
}
