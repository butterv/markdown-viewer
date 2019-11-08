package lexer

import (
	"bytes"
	"regexp"

	"github.com/istsh/markdown-viewer/token"
)

const (
	SHARP             = '#'
	HYPHEN            = '-'
	TAB               = '\t'
	SPACE             = ' '
	LINE_BREAK_CODE_N = '\n'
	LINE_BREAK_CODE_R = '\r'
	GT                = '>'
	BACK_QUOTE        = '`'
	ASTERISK          = '*'
	UNDER_SCORE       = '_'
	LBRACKET          = "["
	RBRACKET          = "]"
	LPAREN            = "("
	RPAREN            = ")"
)

var (
	regexpUrl = regexp.MustCompile(`\[.*]\(http(s)?://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?\)`)
)

type Lexer struct {
	input []byte

	currentPosition int
	nextPosition    int

	currentCh byte
	beforeCh  byte

	startedBackQuoteArea   bool
	startedAsteriskToken   token.TokenType
	startedUnderScoreToken token.TokenType
	startedLinkText        token.TokenType
}

func New(input []byte) *Lexer {
	// 必ず最後は改行コードで終わらせたい
	if !bytes.HasSuffix(input, []byte{LINE_FEED_CODE_N}) {
		input = append(input, LINE_FEED_CODE_N)
	}

	l := &Lexer{
		input:                  input,
		beforeCh:               LINE_FEED_CODE_N, // 直前の文字の初期値は改行コード
		startedAsteriskToken:   token.NONE,
		startedUnderScoreToken: token.NONE,
	}

	return l
}

func (l *Lexer) NextToken() token.Token {
	// 1文字進める
	l.readChar()

	// 空白もタブも改行も、全てスキップせずに解析していく

	var tok token.Token

	switch l.currentCh {
	case '#':
		if isLineBreakCode(l.beforeCh) {
			literal := l.readHeading()
			nextCh := l.peekNextChar()
			if isSpace(nextCh) {
				tok = newToken(token.GetHeadingToken(len(literal)))
				// 空白をスキップする
				l.readChar()
			} else if isLineBreakCode(nextCh) && len(literal) == 3 {
				l.readChar()
				tok = newToken(token.HORIZON)
			} else {
				l.readChar()
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, l.readString()...)
				tok = newToken(token.STRING, tmpChs...)
			}
		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case '-':
		if isLineBreakCode(l.beforeCh) {
			literal := l.readHyphen()
			nextCh := l.peekNextChar()
			if isLineBreakCode(nextCh) && len(literal) == 3 {
				tok = newToken(token.HORIZON)
			} else if isSpace(nextCh) && len(literal) == 1 {
				tok = newToken(token.HYPHEN, literal...)
			} else {
				l.readChar()
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, l.readString()...)
				tok = newToken(token.STRING, tmpChs...)
			}
		} else if isTab(l.beforeCh) {
			if isLineBreakCode(l.twoBeforeChar()) {
				tok = newToken(token.HYPHEN, l.currentCh)
			} else if isTab(l.twoBeforeChar()) && isLineBreakCode(l.threeBeforeChar()) {
				tok = newToken(token.HYPHEN, l.currentCh)
			} else {
				tok = newToken(token.STRING, l.currentCh)
			}
		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case '\t':
		// TODO: Tabはいらないかも
		literal := l.readTab()
		tok = newToken(token.GetTabToken(len(literal)))
	case ' ':
		tok = newToken(token.SPACE)
	case '\n', '\r':
		l.startedBackQuoteArea = false
		l.startedAsteriskToken = token.NONE
		l.startedUnderScoreToken = token.NONE
		tok = newToken(token.LINE_FEED_CODE)
	case '>':
		if isLineBreakCode(l.beforeCh) {
			literal := l.readCitation()
			tok = newToken(token.GetCitationToken(len(literal)))
		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case '`':
		if l.startedBackQuoteArea {
			nextCh := l.peekNextChar()
			if isSpace(nextCh) || isLineBreakCode(nextCh) {
				tok = newToken(token.BACK_QUOTE_FINISH)
			} else {
				tok = newToken(token.STRING, l.readString()...)
			}
			l.startedBackQuoteArea = false
		} else {
			switch {
			case isLineBreakCode(l.beforeCh):
				if l.existsByEndOfLine([]byte("` ")) {
					l.startedBackQuoteArea = true
					tok = newToken(token.BACK_QUOTE_BEGIN)
				} else {
					l.startedBackQuoteArea = false
					tok = newToken(token.STRING, l.readString()...)
				}
			case isSpace(l.beforeCh):
				if l.existsByEndOfLine([]byte("` ")) {
					l.startedBackQuoteArea = true
					tok = newToken(token.BACK_QUOTE_BEGIN)
				} else {
					l.startedBackQuoteArea = false
					tok = newToken(token.STRING, l.readString()...)
				}
			default:
				l.startedBackQuoteArea = false
				tok = newToken(token.STRING, l.readString()...)
			}
		}
	case '*':
		if l.startedAsteriskToken != token.NONE {
			switch l.startedAsteriskToken {
			case token.ASTERISK_ITALIC_BEGIN:
				nextCh := l.peekNextChar()
				if isSpace(nextCh) || isLineBreakCode(nextCh) {
					tok = newToken(token.ASTERISK_ITALIC_FINISH)
				} else {
					tok = newToken(token.STRING, l.readString()...)
				}
				l.startedAsteriskToken = token.NONE
			case token.ASTERISK_BOLD_BEGIN:
				if isAsterisk(l.peekNextChar()) {
					peek2ndOrderChar := l.peek2ndOrderChar()
					if isSpace(peek2ndOrderChar) || isLineBreakCode(peek2ndOrderChar) {
						l.readAsterisk()
						tok = newToken(token.ASTERISK_BOLD_FINISH)
					} else {
						tok = newToken(token.STRING, l.readString()...)
					}
				} else {
					tok = newToken(token.STRING, l.readString()...)
				}
				l.startedAsteriskToken = token.NONE
			case token.ASTERISK_ITALIC_BOLD_BEGIN:
				if isAsterisk(l.peekNextChar()) {
					if isAsterisk(l.peek2ndOrderChar()) {
						peek3ndOrderChar := l.peek3ndOrderChar()
						if isSpace(peek3ndOrderChar) || isLineBreakCode(peek3ndOrderChar) {
							l.readAsterisk()
							tok = newToken(token.ASTERISK_ITALIC_BOLD_FINISH)
						} else {
							tok = newToken(token.STRING, l.readString()...)
						}
					} else {
						tok = newToken(token.STRING, l.readString()...)
					}
				} else {
					tok = newToken(token.STRING, l.readString()...)
				}
				l.startedAsteriskToken = token.NONE
			}

		} else {
			switch {
			case isLineBreakCode(l.beforeCh):
				literal := l.readAsterisk()
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, ' ')
				if l.existsByEndOfLine(tmpChs) {
					tokenType := token.GetAsteriskToken(len(literal))
					l.startedAsteriskToken = tokenType
					tok = newToken(tokenType)
				} else if isLineBreakCode(l.peekNextChar()) && len(literal) == 3 {
					tok = newToken(token.HORIZON)
				} else {
					tmpChs = nil
					tmpChs = append(tmpChs, literal...)
					// TODO: \rも検証すべき
					tmpChs = append(tmpChs, '\n')
					if l.existsByEndOfLine(tmpChs) {
						tokenType := token.GetAsteriskToken(len(literal))
						l.startedAsteriskToken = tokenType
						tok = newToken(tokenType)
					} else {
						l.startedAsteriskToken = token.NONE
						var tmpChs []byte
						tmpChs = append(tmpChs, literal...)
						tmpChs = append(tmpChs, l.readString()...)
						tok = newToken(token.STRING, tmpChs...)
					}
				}
			case isSpace(l.beforeCh):
				literal := l.readAsterisk()
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, ' ')
				if l.existsByEndOfLine(tmpChs) {
					tokenType := token.GetAsteriskToken(len(literal))
					l.startedAsteriskToken = tokenType
					tok = newToken(tokenType)
				} else {
					tmpChs = nil
					tmpChs = append(tmpChs, literal...)
					// TODO: \rも検証すべき
					tmpChs = append(tmpChs, '\n')
					if l.existsByEndOfLine(tmpChs) {
						tokenType := token.GetAsteriskToken(len(literal))
						l.startedAsteriskToken = tokenType
						tok = newToken(tokenType)
					} else {
						l.startedAsteriskToken = token.NONE
						var tmpChs []byte
						tmpChs = append(tmpChs, literal...)
						tmpChs = append(tmpChs, l.readString()...)
						tok = newToken(token.STRING, tmpChs...)
					}
				}
			default:
				l.startedAsteriskToken = token.NONE
				tok = newToken(token.STRING, l.readString()...)
			}
		}
	case '_':
		if l.startedUnderScoreToken != token.NONE {
			switch l.startedUnderScoreToken {
			case token.UNDER_SCORE_ITALIC_BEGIN:
				nextCh := l.peekNextChar()
				if isSpace(nextCh) || isLineBreakCode(nextCh) {
					tok = newToken(token.UNDER_SCORE_ITALIC_FINISH)
				} else {
					tok = newToken(token.STRING, l.readString()...)
				}
				l.startedUnderScoreToken = token.NONE
			case token.UNDER_SCORE_BOLD_BEGIN:
				if isUnderScore(l.peekNextChar()) {
					peek2ndOrderChar := l.peek2ndOrderChar()
					if isSpace(peek2ndOrderChar) || isLineBreakCode(peek2ndOrderChar) {
						l.readUnderScore()
						tok = newToken(token.UNDER_SCORE_BOLD_FINISH)
					} else {
						tok = newToken(token.STRING, l.readString()...)
					}
				} else {
					tok = newToken(token.STRING, l.readString()...)
				}
				l.startedUnderScoreToken = token.NONE
			case token.UNDER_SCORE_ITALIC_BOLD_BEGIN:
				if isUnderScore(l.peekNextChar()) {
					if isUnderScore(l.peek2ndOrderChar()) {
						peek3ndOrderChar := l.peek3ndOrderChar()
						if isSpace(peek3ndOrderChar) || isLineBreakCode(peek3ndOrderChar) {
							l.readUnderScore()
							tok = newToken(token.UNDER_SCORE_ITALIC_BOLD_FINISH)
						} else {
							tok = newToken(token.STRING, l.readString()...)
						}
					} else {
						tok = newToken(token.STRING, l.readString()...)
					}
				} else {
					tok = newToken(token.STRING, l.readString()...)
				}
				l.startedUnderScoreToken = token.NONE
			}

		} else {
			switch {
			case isLineBreakCode(l.beforeCh):
				literal := l.readUnderScore()
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, ' ')
				if l.existsByEndOfLine(tmpChs) {
					tokenType := token.GetUnderScoreToken(len(literal))
					l.startedUnderScoreToken = tokenType
					tok = newToken(tokenType)
				} else if isLineBreakCode(l.peekNextChar()) && len(literal) == 3 {
					tok = newToken(token.HORIZON)
				} else {
					tmpChs = nil
					tmpChs = append(tmpChs, literal...)
					// TODO: \rも検証すべき
					tmpChs = append(tmpChs, '\n')
					if l.existsByEndOfLine(tmpChs) {
						tokenType := token.GetUnderScoreToken(len(literal))
						l.startedUnderScoreToken = tokenType
						tok = newToken(tokenType)
					} else {
						l.startedUnderScoreToken = token.NONE
						var tmpChs []byte
						tmpChs = append(tmpChs, literal...)
						tmpChs = append(tmpChs, l.readString()...)
						tok = newToken(token.STRING, tmpChs...)
					}
				}
			case isSpace(l.beforeCh):
				literal := l.readUnderScore()
				var tmpChs []byte
				tmpChs = append(tmpChs, literal...)
				tmpChs = append(tmpChs, ' ')
				if l.existsByEndOfLine(tmpChs) {
					tokenType := token.GetUnderScoreToken(len(literal))
					l.startedUnderScoreToken = tokenType
					tok = newToken(tokenType)
				} else {
					tmpChs = nil
					tmpChs = append(tmpChs, literal...)
					// TODO: \rも検証すべき
					tmpChs = append(tmpChs, '\n')
					if l.existsByEndOfLine(tmpChs) {
						tokenType := token.GetUnderScoreToken(len(literal))
						l.startedUnderScoreToken = tokenType
						tok = newToken(tokenType)
					} else {
						l.startedUnderScoreToken = token.NONE
						var tmpChs []byte
						tmpChs = append(tmpChs, literal...)
						tmpChs = append(tmpChs, l.readString()...)
						tok = newToken(token.STRING, tmpChs...)
					}
				}
			default:
				l.startedUnderScoreToken = token.NONE
				tok = newToken(token.STRING, l.readString()...)
			}
		}
	case '[':
		chs := l.untilLineFeedCode()
		matchedChs := regexpUrl.Find(chs)
		if matchedChs != nil {
			l.startedLinkText = token.LINK_TEXT_BEGIN
			tok = newToken(token.LINK_TEXT_BEGIN, l.currentCh)
		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case ']':
		if l.startedLinkText == token.LINK_TEXT_BEGIN {
			tok = newToken(token.LINK_TEXT_FINISH, l.currentCh)
		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case '(':
		if l.startedLinkText == token.LINK_TEXT_BEGIN {
			tok = newToken(token.LINK_BEGIN, l.currentCh)
		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case ')':
		if l.startedLinkText == token.LINK_TEXT_BEGIN {
			l.startedLinkText = token.NONE
			tok = newToken(token.LINK_FINISH, l.currentCh)
		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case 0:
		tok = newToken(token.EOF)
	default:
		tok = newToken(token.STRING, l.readString()...)
	}

	return tok
}

func newToken(tokenType token.TokenType, chs ...byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: chs,
	}
}

func (l *Lexer) readChar() {
	if l.currentPosition > 0 {
		// 直前の文字をセット
		l.beforeCh = l.currentCh
	}

	// 次の文字が存在するか
	if l.nextPosition >= len(l.input) {
		// 次の文字は存在しない(まだ何も読み込んでいない or ファイルの終わり)
		l.currentCh = 0
	} else {
		// 次の文字をセット
		l.currentCh = l.input[l.nextPosition]
	}
	// 数値を1つ進める
	l.currentPosition = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) peekNextChar() byte {
	// 次の文字を覗き見る
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

func (l *Lexer) peek2ndOrderChar() byte {
	// 2つ次の文字を覗き見る
	if l.nextPosition+1 >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition+1]
	}
}

func (l *Lexer) peek3ndOrderChar() byte {
	// 3つ次の文字を覗き見る
	if l.nextPosition+2 >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition+2]
	}
}

func (l *Lexer) peek4ndOrderChar() byte {
	// 4つ次の文字を覗き見る
	if l.nextPosition+3 >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition+3]
	}
}

func (l *Lexer) existsByEndOfLine(chs []byte) bool {
	// 現在の位置から次の改行コードでの文字を確認するだけなので、readCharは実行しない
	position := l.currentPosition

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
		if isLineBreakCode(ch) {
			break
		}
	}

	return false
}

func (l *Lexer) twoBeforeChar() byte {
	// 2つ前の文字を見る
	if l.nextPosition < 3 {
		return 0
	}
	return l.input[l.nextPosition-3]
}

func (l *Lexer) threeBeforeChar() byte {
	// 3つ前の文字を見る
	if l.nextPosition < 4 {
		return 0
	}
	return l.input[l.nextPosition-4]
}

func isSharp(ch byte) bool {
	return ch == '#'
}

func (l *Lexer) readHeading() []byte {
	position := l.currentPosition

	for {
		nextCh := l.peekNextChar()
		if !isSharp(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.currentPosition+1]
}

func (l *Lexer) readString() []byte {
	position := l.currentPosition

	// 次の両方を満たす場合、文字列と判断して読み進める
	// 1. 次の文字が改行コードではない
	// 2. 現在の文字が空白でない
	for {
		nextCh := l.peekNextChar()
		var breakFlg bool
		switch {
		case isSpace(nextCh), isLineBreakCode(nextCh):
			breakFlg = true
		case isBackQuote(nextCh):
			peeked2ndOrderCh := l.peek2ndOrderChar()
			if isSpace(peeked2ndOrderCh) || isLineBreakCode(peeked2ndOrderCh) {
				breakFlg = true
			}
		case isAsterisk(nextCh):
			switch l.startedAsteriskToken {
			case token.ASTERISK_ITALIC_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isSpace(peeked2ndOrderCh) || isLineBreakCode(peeked2ndOrderCh) {
					breakFlg = true
				}
			case token.ASTERISK_BOLD_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isAsterisk(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isSpace(peeked3ndOrderCh) || isLineBreakCode(peeked3ndOrderCh) {
						breakFlg = true
					}
				}
			case token.ASTERISK_ITALIC_BOLD_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isAsterisk(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isAsterisk(peeked3ndOrderCh) {
						peeked4ndOrderCh := l.peek4ndOrderChar()
						if isSpace(peeked4ndOrderCh) || isLineBreakCode(peeked4ndOrderCh) {
							breakFlg = true
						}
					}
				}
			}
		case isUnderScore(nextCh):
			switch l.startedUnderScoreToken {
			case token.UNDER_SCORE_ITALIC_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isSpace(peeked2ndOrderCh) || isLineBreakCode(peeked2ndOrderCh) {
					breakFlg = true
				}
			case token.UNDER_SCORE_BOLD_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isUnderScore(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isSpace(peeked3ndOrderCh) || isLineBreakCode(peeked3ndOrderCh) {
						breakFlg = true
					}
				}
			case token.UNDER_SCORE_ITALIC_BOLD_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isUnderScore(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isUnderScore(peeked3ndOrderCh) {
						peeked4ndOrderCh := l.peek4ndOrderChar()
						if isSpace(peeked4ndOrderCh) || isLineBreakCode(peeked4ndOrderCh) {
							breakFlg = true
						}
					}
				}
			}
		case isRightBracket(nextCh), isRightParen(nextCh):
			if l.startedLinkText == token.LINK_TEXT_BEGIN {
				breakFlg = true
			}
		}

		if breakFlg {
			break
		}

		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.currentPosition+1]
}

func isLineBreakCode(ch byte) bool {
	return ch == '\n' || ch == '\r'
}

func isTab(ch byte) bool {
	return ch == '\t'
}

func (l *Lexer) readTab() []byte {
	position := l.currentPosition

	for {
		nextCh := l.peekNextChar()
		if !isTab(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.currentPosition+1]
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
	position := l.currentPosition

	for {
		nextCh := l.peekNextChar()
		if !isAsterisk(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.currentPosition+1]
}

func isUnderScore(ch byte) bool {
	return ch == '_'
}

func (l *Lexer) readUnderScore() []byte {
	position := l.currentPosition

	for {
		nextCh := l.peekNextChar()
		if !isUnderScore(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.currentPosition+1]
}

func isCitation(ch byte) bool {
	return ch == '>'
}

func (l *Lexer) readCitation() []byte {
	position := l.currentPosition

	for {
		nextCh := l.peekNextChar()
		if !isCitation(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.currentPosition+1]
}

func isHyphen(ch byte) bool {
	return ch == '-'
}

func (l *Lexer) readHyphen() []byte {
	position := l.currentPosition

	for {
		nextCh := l.peekNextChar()
		if !isHyphen(nextCh) {
			break
		}
		// 文字が途切れるまで読み込む
		l.readChar()
	}

	// positionから、readCharで進んだところまで抽出
	return l.input[position : l.currentPosition+1]
}

func isRightBracket(ch byte) bool {
	return ch == ']'
}

func isRightParen(ch byte) bool {
	return ch == ')'
}

func (l *Lexer) untilLineFeedCode() []byte {
	position := l.currentPosition

	var chs []byte
	for {
		ch := l.input[position]
		if ch == LINE_FEED_CODE_N || ch == LINE_FEED_CODE_R {
			break
		}
		chs = append(chs, ch)
		position++
	}

	return chs
}
