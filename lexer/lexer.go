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

	startedBackQuoteArea   bool            // バッククォートエリアが開始されているか
	startedAsteriskToken   token.TokenType // アスタリスクエリアが開始されているか
	startedUnderScoreToken token.TokenType // アンダースコアエリアが開始されているか
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
		input:                  input,
		justBeforeCh:           '\n', // 直前の文字の初期値は改行コード
		startedAsteriskToken:   token.NONE,
		startedUnderScoreToken: token.NONE,
	}

	//l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	// 1文字進める
	l.readChar()

	// 空白もタブも改行も、全てスキップせずに解析していく

	var tok token.Token

	// fmt.Printf("%q\n", l.startedAsteriskToken)

	switch l.ch {
	case '#':
		if isLineFeedCode(l.justBeforeCh) {
			literal := l.readHeading()
			nextCh := l.peekNextChar()
			if isSpace(nextCh) {
				tok = newToken(token.GetHeadingToken(len(literal)))
				// 空白をスキップする
				l.readChar()
			} else {
				l.readChar()
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, l.readString()...)
				tok = newTokenWithLiteral(token.STRING, tmpChs)
			}
		} else {
			tok = newTokenWithLiteral(token.STRING, l.readString())
		}
	case '-':
		tok = newToken(token.HYPHEN)
	case '\t':
		literal := l.readTab()
		tok = newToken(token.GetTabToken(len(literal)))
	case ' ':
		//if l.lookBackChar() == '\n' {
		//	literal, cnt := l.readSpace()
		//	tok.Literal = literal
		//	tok.Type = token.GetSpaceToken(cnt)
		//} else {
		//	tok = newToken(token.SPACE1, l.ch)
		//}
		tok = newToken(token.SPACE)
	case '\n', '\r':
		l.startedBackQuoteArea = false
		l.startedAsteriskToken = token.NONE
		l.startedUnderScoreToken = token.NONE
		tok = newToken(token.LINE_FEED_CODE)
	case '>':
		if isLineFeedCode(l.justBeforeCh) {
			literal := l.readCitation()
			tok = newToken(token.GetCitationToken(len(literal)))
		} else {
			tok = newTokenWithLiteral(token.STRING, l.readString())
		}
	case '`':
		if l.startedBackQuoteArea {
			// バッククォートエリアはすでに始まっている
			nextCh := l.peekNextChar()
			if isSpace(nextCh) || isLineFeedCode(nextCh) {
				tok = newToken(token.BACK_QUOTE_FINISH)
			} else {
				tok.Literal = l.readString()
				tok.Type = token.STRING
			}
			l.startedBackQuoteArea = false
		} else {
			// バッククオートエリアを始めようとしている
			switch {
			case isLineFeedCode(l.justBeforeCh):
				if l.existsByEndOfLine([]byte("` ")) {
					l.startedBackQuoteArea = true
					tok = newToken(token.BACK_QUOTE, l.ch)
				} else {
					l.startedBackQuoteArea = false
					tok.Literal = l.readString()
					tok.Type = token.STRING
				}
			case isSpace(l.justBeforeCh):
				if l.existsByEndOfLine([]byte("` ")) {
					l.startedBackQuoteArea = true
					tok = newToken(token.BACK_QUOTE, l.ch)
				} else {
					l.startedBackQuoteArea = false
					tok.Literal = l.readString()
					tok.Type = token.STRING
				}
			default:
				l.startedBackQuoteArea = false
				tok.Literal = l.readString()
				tok.Type = token.STRING
			}
		}
	case '*':
		if l.startedAsteriskToken != token.NONE {
			// アスタリスクエリアはすでに始まっている
			switch l.startedAsteriskToken {
			case token.ASTERISK_ITALIC:
				nextCh := l.peekNextChar()
				if isSpace(nextCh) || isLineFeedCode(nextCh) {
					tok = newToken(token.ASTERISK_ITALIC, l.ch)
				} else {
					tok.Literal = l.readString()
					tok.Type = token.STRING
				}
				l.startedAsteriskToken = token.NONE
			case token.ASTERISK_BOLD:
				if isAsterisk(l.peekNextChar()) {
					peek2ndOrderChar := l.peek2ndOrderChar()
					if isSpace(peek2ndOrderChar) || isLineFeedCode(peek2ndOrderChar) {
						literal := l.readAsterisk()
						tok.Literal = literal
						tok.Type = token.ASTERISK_BOLD
					} else {
						tok.Literal = l.readString()
						tok.Type = token.STRING
					}
				} else {
					tok.Literal = l.readString()
					tok.Type = token.STRING
				}
				l.startedAsteriskToken = token.NONE
			case token.ASTERISK_ITALIC_BOLD:
				if isAsterisk(l.peekNextChar()) {
					if isAsterisk(l.peek2ndOrderChar()) {
						peek3ndOrderChar := l.peek3ndOrderChar()
						if isSpace(peek3ndOrderChar) || isLineFeedCode(peek3ndOrderChar) {
							literal := l.readAsterisk()
							tok.Literal = literal
							tok.Type = token.ASTERISK_ITALIC_BOLD
						} else {
							tok.Literal = l.readString()
							tok.Type = token.STRING
						}
					} else {
						tok.Literal = l.readString()
						tok.Type = token.STRING
					}
				} else {
					tok.Literal = l.readString()
					tok.Type = token.STRING
				}
				l.startedAsteriskToken = token.NONE
			}

		} else {
			// アスタリスクエリアを始めようとしている
			literal := l.readAsterisk()
			//fmt.Printf("asterisk: %s\n", literal)

			var beforeCh byte
			switch len(literal) {
			case 1:
				beforeCh = l.lookBackChar()
			case 2:
				beforeCh = l.twoBeforeChar()
			case 3:
				beforeCh = l.threeBeforeChar()
			}

			switch {
			case isLineFeedCode(beforeCh):
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, ' ')
				if l.existsByEndOfLine(tmpChs) {
					tokenType := token.GetAsteriskToken(len(literal))
					l.startedAsteriskToken = tokenType
					tok.Literal = literal
					tok.Type = tokenType
				} else {
					tmpChs = nil
					tmpChs = append(tmpChs, literal...)
					// TODO: \rも検証すべき
					tmpChs = append(tmpChs, '\n')
					if l.existsByEndOfLine(tmpChs) {
						tokenType := token.GetAsteriskToken(len(literal))
						l.startedAsteriskToken = tokenType
						tok.Literal = literal
						tok.Type = tokenType
					} else {
						l.startedAsteriskToken = token.NONE
						tmpChs := literal
						tmpChs = append(tmpChs, l.readString()...)
						tok.Literal = tmpChs
						tok.Type = token.STRING
					}
				}
			case isSpace(beforeCh):
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, ' ')
				if l.existsByEndOfLine(tmpChs) {
					tokenType := token.GetAsteriskToken(len(literal))
					l.startedAsteriskToken = tokenType
					tok.Literal = literal
					tok.Type = tokenType
				} else {
					tmpChs = nil
					tmpChs = append(tmpChs, literal...)
					// TODO: \rも検証すべき
					tmpChs = append(tmpChs, '\n')
					if l.existsByEndOfLine(tmpChs) {
						tokenType := token.GetAsteriskToken(len(literal))
						l.startedAsteriskToken = tokenType
						tok.Literal = literal
						tok.Type = tokenType
					} else {
						l.startedAsteriskToken = token.NONE
						var tmpChs []byte
						tmpChs = append(tmpChs, literal...)
						tmpChs = append(tmpChs, l.readString()...)
						tok.Literal = tmpChs
						tok.Type = token.STRING
					}
				}
			default:
				l.startedAsteriskToken = token.NONE
				tok.Literal = l.readString()
				tok.Type = token.STRING
			}
		}
	case '_':
		if l.startedUnderScoreToken != token.NONE {
			// アンダースコアはすでに始まっている
			switch l.startedUnderScoreToken {
			case token.UNDER_SCORE_ITALIC:
				nextCh := l.peekNextChar()
				if isSpace(nextCh) || isLineFeedCode(nextCh) {
					tok = newToken(token.UNDER_SCORE_ITALIC, l.ch)
				} else {
					tok.Literal = l.readString()
					tok.Type = token.STRING
				}
				l.startedUnderScoreToken = token.NONE
			case token.UNDER_SCORE_BOLD:
				if isUnderScore(l.peekNextChar()) {
					peek2ndOrderChar := l.peek2ndOrderChar()
					if isSpace(peek2ndOrderChar) || isLineFeedCode(peek2ndOrderChar) {
						literal := l.readUnderScore()
						tok.Literal = literal
						tok.Type = token.UNDER_SCORE_BOLD
					} else {
						tok.Literal = l.readString()
						tok.Type = token.STRING
					}
				} else {
					tok.Literal = l.readString()
					tok.Type = token.STRING
				}
				l.startedUnderScoreToken = token.NONE
			case token.UNDER_SCORE_ITALIC_BOLD:
				if isUnderScore(l.peekNextChar()) {
					if isUnderScore(l.peek2ndOrderChar()) {
						peek3ndOrderChar := l.peek3ndOrderChar()
						if isSpace(peek3ndOrderChar) || isLineFeedCode(peek3ndOrderChar) {
							literal := l.readUnderScore()
							tok.Literal = literal
							tok.Type = token.UNDER_SCORE_ITALIC_BOLD
						} else {
							tok.Literal = l.readString()
							tok.Type = token.STRING
						}
					} else {
						tok.Literal = l.readString()
						tok.Type = token.STRING
					}
				} else {
					tok.Literal = l.readString()
					tok.Type = token.STRING
				}
				l.startedUnderScoreToken = token.NONE
			}

		} else {
			// アンダースコアエリアを始めようとしている
			literal := l.readUnderScore()

			var beforeCh byte
			switch len(literal) {
			case 1:
				beforeCh = l.lookBackChar()
			case 2:
				beforeCh = l.twoBeforeChar()
			case 3:
				beforeCh = l.threeBeforeChar()
			}

			switch {
			case isLineFeedCode(beforeCh):
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, ' ')
				if l.existsByEndOfLine(tmpChs) {
					tokenType := token.GetUnderScoreToken(len(literal))
					l.startedUnderScoreToken = tokenType
					tok.Literal = literal
					tok.Type = tokenType
				} else {
					tmpChs = nil
					tmpChs = append(tmpChs, literal...)
					// TODO: \rも検証すべき
					tmpChs = append(tmpChs, '\n')
					if l.existsByEndOfLine(tmpChs) {
						tokenType := token.GetUnderScoreToken(len(literal))
						l.startedUnderScoreToken = tokenType
						tok.Literal = literal
						tok.Type = tokenType
					} else {
						l.startedUnderScoreToken = token.NONE
						var tmpChs []byte
						tmpChs = append(tmpChs, literal...)
						tmpChs = append(tmpChs, l.readString()...)
						tok.Literal = tmpChs
						tok.Type = token.STRING
					}
				}
			case isSpace(beforeCh):
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, ' ')
				if l.existsByEndOfLine(tmpChs) {
					tokenType := token.GetUnderScoreToken(len(literal))
					l.startedUnderScoreToken = tokenType
					tok.Literal = literal
					tok.Type = tokenType
				} else {
					tmpChs = nil
					tmpChs = append(tmpChs, literal...)
					// TODO: \rも検証すべき
					tmpChs = append(tmpChs, '\n')
					if l.existsByEndOfLine(tmpChs) {
						tokenType := token.GetUnderScoreToken(len(literal))
						l.startedUnderScoreToken = tokenType
						tok.Literal = literal
						tok.Type = tokenType
					} else {
						l.startedUnderScoreToken = token.NONE
						var tmpChs []byte
						tmpChs = append(tmpChs, literal...)
						tmpChs = append(tmpChs, l.readString()...)
						tok.Literal = tmpChs
						tok.Type = token.STRING
					}
				}
			default:
				l.startedUnderScoreToken = token.NONE
				tok.Literal = l.readString()
				tok.Type = token.STRING
			}
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
		tok.Literal = nil
		tok.Type = token.EOF
	default:
		tok.Literal = l.readString()
		tok.Type = token.STRING
	}

	return tok
}

func newToken(tokenType token.TokenType) token.Token {
	// Tokenオブジェクトを初期化する
	return token.Token{
		Type: tokenType,
	}
}

func newTokenWithLiteral(tokenType token.TokenType, chs []byte) token.Token {
	// Tokenオブジェクトを初期化する
	return token.Token{
		Type:    tokenType,
		Literal: chs,
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

func (l *Lexer) peek2ndOrderChar() byte {
	// 次の次の文字を覗き見る
	if l.readPosition+1 >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition+1]
	}
}

func (l *Lexer) peek3ndOrderChar() byte {
	// 次の次の次の文字を覗き見る
	if l.readPosition+2 >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition+2]
	}
}

func (l *Lexer) peek4ndOrderChar() byte {
	// 次の次の次の次の文字を覗き見る
	if l.readPosition+3 >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition+3]
	}
}

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

func (l *Lexer) lookBackChar() byte {
	// 直前の文字を見る
	if l.readPosition < 2 {
		return 0
	}
	return l.input[l.readPosition-2]
}

func (l *Lexer) twoBeforeChar() byte {
	// 2つ前の文字を見る
	if l.readPosition < 3 {
		return 0
	}
	return l.input[l.readPosition-3]
}

func (l *Lexer) threeBeforeChar() byte {
	// 3つ前の文字を見る
	if l.readPosition < 4 {
		return 0
	}
	return l.input[l.readPosition-4]
}

func isSharp(ch byte) bool {
	return ch == '#'
}

func (l *Lexer) readHeading() []byte {
	position := l.position

	for {
		nextCh := l.peekNextChar()
		if !isSharp(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1]
}

func (l *Lexer) readString() []byte {
	position := l.position

	// 次の両方を満たす場合、文字列と判断して読み進める
	// 1. 次の文字が改行コードではない
	// 2. 現在の文字が空白でない
	for {
		nextCh := l.peekNextChar()
		var breakFlg bool
		switch {
		case isSpace(nextCh), isLineFeedCode(nextCh):
			breakFlg = true
		case isBackQuote(nextCh):
			peeked2ndOrderCh := l.peek2ndOrderChar()
			if isSpace(peeked2ndOrderCh) || isLineFeedCode(peeked2ndOrderCh) {
				breakFlg = true
			}
		case isAsterisk(nextCh):
			switch l.startedAsteriskToken {
			case token.ASTERISK_ITALIC:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isSpace(peeked2ndOrderCh) || isLineFeedCode(peeked2ndOrderCh) {
					breakFlg = true
				}
			case token.ASTERISK_BOLD:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isAsterisk(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isSpace(peeked3ndOrderCh) || isLineFeedCode(peeked3ndOrderCh) {
						breakFlg = true
					}
				}
			case token.ASTERISK_ITALIC_BOLD:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isAsterisk(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isAsterisk(peeked3ndOrderCh) {
						peeked4ndOrderCh := l.peek4ndOrderChar()
						if isSpace(peeked4ndOrderCh) || isLineFeedCode(peeked4ndOrderCh) {
							breakFlg = true
						}
					}
				}
			}
		case isUnderScore(nextCh):
			switch l.startedUnderScoreToken {
			case token.UNDER_SCORE_ITALIC:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isSpace(peeked2ndOrderCh) || isLineFeedCode(peeked2ndOrderCh) {
					breakFlg = true
				}
			case token.UNDER_SCORE_BOLD:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isUnderScore(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isSpace(peeked3ndOrderCh) || isLineFeedCode(peeked3ndOrderCh) {
						breakFlg = true
					}
				}
			case token.UNDER_SCORE_ITALIC_BOLD:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isUnderScore(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isUnderScore(peeked3ndOrderCh) {
						peeked4ndOrderCh := l.peek4ndOrderChar()
						if isSpace(peeked4ndOrderCh) || isLineFeedCode(peeked4ndOrderCh) {
							breakFlg = true
						}
					}
				}
			}
		}

		if breakFlg {
			break
		}

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

func (l *Lexer) readTab() []byte {
	position := l.position

	for {
		nextCh := l.peekNextChar()
		if !isTab(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1]
}

func isSpace(ch byte) bool {
	return ch == ' '
}

func isBackQuote(ch byte) bool {
	return ch == '`'
}

func isAsterisk(ch byte) bool {
	return ch == '*'
}

func (l *Lexer) readAsterisk() []byte {
	position := l.position

	for {
		nextCh := l.peekNextChar()
		if !isAsterisk(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1]
}

func isUnderScore(ch byte) bool {
	return ch == '_'
}

func (l *Lexer) readUnderScore() []byte {
	position := l.position

	for {
		nextCh := l.peekNextChar()
		if !isUnderScore(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1]
}

func isCitation(ch byte) bool {
	return ch == '>'
}

func (l *Lexer) readCitation() []byte {
	position := l.position

	for {
		nextCh := l.peekNextChar()
		if !isCitation(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.position+1]
}

//func (l *Lexer) readSpace() ([]byte, int) {
//	position := l.position
//
//	cnt := 1
//	for isSpace(l.ch) && isSpace(l.peekNextChar()) && cnt <= 4 {
//		// 文字が途切れるまで読み込む
//		l.readChar()
//		cnt++
//	}
//	// positionから、readCharで進んだところまで抽出
//	return l.input[position : l.position+1], cnt
//}

//func isDigit(ch byte) bool {
//	// 数値
//	return '0' <= ch && ch <= '9'
//}
