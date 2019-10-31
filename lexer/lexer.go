package lexer

import (
	"bytes"

	"github.com/istsh/markdown-viewer/token"
)

type Lexer struct {
	input []byte // 入力

	position     int // 入力における現在の位置(現在の文字を指し示す)
	readPosition int // これから読み込む位置(現在の文字の次)

	ch           byte // 現在検査中の文字
	justBeforeCh byte // 直前の文字

	// 判定待ちのpositionが合ってもいいかも
	// italicやboldの判定は、間にいくつかの文字があってから閉じる文字がくるので、閉じる文字があるかないかによって、
	// 開始の文字のtypeを変えないといけない。
	// これは他の文字の判定でも使えるかもしれない。
	// それか、ある開始文字があったら、それがただの文字列なのか何かの開始文字なのかを判定するために、
	// その行を必要な分peekしてもいいかもしれない。この方が都合よさそう。

	startedBackQuoteArea bool // バッククォートエリアが開始されているか
	startedItalic        bool // 斜体が開始されているか
	startedBold          bool // 強調が開始されているか
}

//func (l *Lexer) GetInput() []byte {
//	return l.input
//}

func New(input []byte) *Lexer {
	// 必ず最後は改行コードで終わらせたい
	if !bytes.HasSuffix(input, []byte("\n")) {
		input = append(input, '\n')
	}

	l := &Lexer{
		input:        input,
		justBeforeCh: '\n', // 直前の文字の初期値は改行コード
	}

	//l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	// 1文字進める
	l.readChar()

	// 空白もタブも改行も、全てスキップせずに解析していく

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
		l.startedBackQuoteArea = false
		l.startedItalic = false
		l.startedBold = false
		tok = newToken(token.LINE_FEED_CODE_N, l.ch)
	case '\r':
		l.startedBackQuoteArea = false
		l.startedItalic = false
		l.startedBold = false
		tok = newToken(token.LINE_FEED_CODE_R, l.ch)
	case '>':
		if isLineFeedCode(l.justBeforeCh) {
			tok = newToken(token.CITATION, l.ch)
		}
	case '`':
		switch {
		case isLineFeedCode(l.justBeforeCh):
			if l.existsByEndOfLine([]byte("` ")) {
				l.startedBackQuoteArea = true
				tok = newToken(token.BACK_QUOTE, l.ch)
			} else {
				l.startedBackQuoteArea = false
				tok = newToken(token.STRING, l.ch)
			}
		case isSpace(l.justBeforeCh):
		}

	case '*':

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
		tok.Literal = nil
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
		Literal: []byte{ch},
	}
}

func (l *Lexer) readChar() {
	if l.position > 0 {
		// 直前の文字をセット
		l.justBeforeCh = l.ch
	}

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

//func (l *Lexer) peek2ndOrderChar() byte {
//	// 次の次の文字を覗き見る
//	if l.readPosition+1 >= len(l.input) {
//		return 0
//	} else {
//		return l.input[l.readPosition+1]
//	}
//}

func (l *Lexer) existsByEndOfLine(chs []byte) bool {
	// 現在の位置から次の改行コードでの文字を確認するだけなので、readCharは実行しない
	position := l.position

	var tmp []byte
	for {
		// ポジションを進める
		position++
		ch := l.input[position]
		tmp = append(tmp, ch)
		if bytes.Contains(tmp, chs) {
			// 1つでも見つかればOK
			return true
		}
		if isLineFeedCode(ch) {
			break
		}
	}

	return false
}

//func (l *Lexer) lookBackChar() byte {
//	// 直前の文字を見る
//	if l.readPosition < 2 {
//		return 0
//	}
//	return l.input[l.readPosition-2]
//}

func isHeading(ch byte) bool {
	return ch == '#'
}

func (l *Lexer) readHeading() ([]byte, int) {
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

func (l *Lexer) readString() []byte {
	position := l.position

	// 次の両方を満たす場合、文字列と判断して読み進める
	// 1. 次の文字が改行コードではない
	// 2. 現在の文字が空白でない
	for !isSpace(l.ch) && !isLineFeedCode(l.peekNextChar()) {
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

func (l *Lexer) readTab() ([]byte, int) {
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

func (l *Lexer) readSpace() ([]byte, int) {
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
