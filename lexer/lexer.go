package lexer

import (
	"bytes"
	"regexp"

	"github.com/istsh/markdown-viewer/token"
)

const (
	SHARP            = '#'
	HYPHEN           = '-'
	TAB              = '\t'
	SPACE            = ' '
	LINE_FEED_CODE_N = '\n'
	LINE_FEED_CODE_R = '\r'
	GT               = '>'
	BACK_QUOTE       = '`'
	ASTERISK         = '*'
	UNDER_SCORE      = '_'
	LBRACKET         = "["
	RBRACKET         = "]"
	LPAREN           = "("
	RPAREN           = ")"
)

type Lexer struct {
	input []byte // 入力

	incompleteChs []byte // まだ検証が完了していない
	completedChs  []byte // 検証が完了した

	currentPosition int // 入力における現在の位置(現在の文字を指し示す)
	nextPosition    int // これから読み込む位置(現在の文字の次)

	currentCh byte // 現在検査中の文字
	beforeCh  byte // 直前の文字

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
	// TODO:
	// 1行全ての検証が終わるまでは未確定状態として検証を進めて、
	// 改行コードがきたら確定する
	// といいかも。

	// 必ず最後は改行コードで終わらせたい
	if !bytes.HasSuffix(input, []byte("\n")) {
		input = append(input, '\n')
	}

	l := &Lexer{
		input:                  input,
		beforeCh:               '\n', // 直前の文字の初期値は改行コード
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

	switch l.currentCh {
	case '#':
		if isLineFeedCode(l.beforeCh) {
			literal := l.readHeading()
			nextCh := l.peekNextChar()
			if isSpace(nextCh) {
				tok = newToken(token.GetHeadingToken(len(literal)))
				// 空白をスキップする
				l.readChar()
			} else if isLineFeedCode(nextCh) && len(literal) == 3 {
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
		if isLineFeedCode(l.beforeCh) {
			literal := l.readHyphen()
			nextCh := l.peekNextChar()
			if isLineFeedCode(nextCh) && len(literal) == 3 {
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
		if isLineFeedCode(l.beforeCh) {
			literal := l.readCitation()
			tok = newToken(token.GetCitationToken(len(literal)))
		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case '`':
		if l.startedBackQuoteArea {
			// バッククォートエリアはすでに始まっている
			nextCh := l.peekNextChar()
			if isSpace(nextCh) || isLineFeedCode(nextCh) {
				tok = newToken(token.BACK_QUOTE_FINISH)
			} else {
				tok = newToken(token.STRING, l.readString()...)
			}
			l.startedBackQuoteArea = false
		} else {
			// バッククオートエリアを始めようとしている
			switch {
			case isLineFeedCode(l.beforeCh):
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
			// アスタリスクエリアはすでに始まっている
			switch l.startedAsteriskToken {
			case token.ASTERISK_ITALIC_BEGIN:
				nextCh := l.peekNextChar()
				if isSpace(nextCh) || isLineFeedCode(nextCh) {
					tok = newToken(token.ASTERISK_ITALIC_FINISH)
				} else {
					tok = newToken(token.STRING, l.readString()...)
				}
				l.startedAsteriskToken = token.NONE
			case token.ASTERISK_BOLD_BEGIN:
				if isAsterisk(l.peekNextChar()) {
					peek2ndOrderChar := l.peek2ndOrderChar()
					if isSpace(peek2ndOrderChar) || isLineFeedCode(peek2ndOrderChar) {
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
						if isSpace(peek3ndOrderChar) || isLineFeedCode(peek3ndOrderChar) {
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
			case isSpace(beforeCh):
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
			// アンダースコアはすでに始まっている
			switch l.startedUnderScoreToken {
			case token.UNDER_SCORE_ITALIC_BEGIN:
				nextCh := l.peekNextChar()
				if isSpace(nextCh) || isLineFeedCode(nextCh) {
					tok = newToken(token.UNDER_SCORE_ITALIC_FINISH)
				} else {
					tok = newToken(token.STRING, l.readString()...)
				}
				l.startedUnderScoreToken = token.NONE
			case token.UNDER_SCORE_BOLD_BEGIN:
				if isUnderScore(l.peekNextChar()) {
					peek2ndOrderChar := l.peek2ndOrderChar()
					if isSpace(peek2ndOrderChar) || isLineFeedCode(peek2ndOrderChar) {
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
						if isSpace(peek3ndOrderChar) || isLineFeedCode(peek3ndOrderChar) {
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
			// アンダースコアエリアを始めようとしている
			literal := l.readUnderScore()
			if isLineFeedCode(l.beforeCh) && isLineFeedCode(l.peekNextChar()) && len(literal) == 3 {
				l.readChar()
				tok = newToken(token.HORIZON)
			} else {
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
				case isSpace(beforeCh):
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
		}

	case '[':
		var chs []byte
		po := l.currentPosition
		chs = append(chs, l.currentCh)
		var cnt int
		for {
			cnt++
			ch := l.input[po+cnt]
			if ch == '\n' {
				break
			}
			chs = append(chs, ch)
		}
		r := regexp.MustCompile(`[*](https?://[\w/:%#\$&\?\(\)~\.=\+\-]+)`)
		matchedChs := r.Find(chs)
		if matchedChs != nil {

		} else {
			tok = newToken(token.STRING, l.readString()...)
		}
	case ']':
	case '(':
	case ')':

	case 0:
		tok = newToken(token.EOF)
	default:
		tok = newToken(token.STRING, l.readString()...)
	}

	return tok
}

func newToken(tokenType token.TokenType, chs ...byte) token.Token {
	// Tokenオブジェクトを初期化する
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
	// 次の次の文字を覗き見る
	if l.nextPosition+1 >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition+1]
	}
}

func (l *Lexer) peek3ndOrderChar() byte {
	// 次の次の次の文字を覗き見る
	if l.nextPosition+2 >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition+2]
	}
}

func (l *Lexer) peek4ndOrderChar() byte {
	// 次の次の次の次の文字を覗き見る
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
		if isLineFeedCode(ch) {
			break
		}
	}

	return false
}

func (l *Lexer) lookBackChar() byte {
	// 直前の文字を見る
	if l.nextPosition < 2 {
		return 0
	}
	return l.input[l.nextPosition-2]
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
		case isSpace(nextCh), isLineFeedCode(nextCh):
			breakFlg = true
		case isBackQuote(nextCh):
			peeked2ndOrderCh := l.peek2ndOrderChar()
			if isSpace(peeked2ndOrderCh) || isLineFeedCode(peeked2ndOrderCh) {
				breakFlg = true
			}
		case isAsterisk(nextCh):
			switch l.startedAsteriskToken {
			case token.ASTERISK_ITALIC_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isSpace(peeked2ndOrderCh) || isLineFeedCode(peeked2ndOrderCh) {
					breakFlg = true
				}
			case token.ASTERISK_BOLD_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isAsterisk(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isSpace(peeked3ndOrderCh) || isLineFeedCode(peeked3ndOrderCh) {
						breakFlg = true
					}
				}
			case token.ASTERISK_ITALIC_BOLD_BEGIN:
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
			case token.UNDER_SCORE_ITALIC_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isSpace(peeked2ndOrderCh) || isLineFeedCode(peeked2ndOrderCh) {
					breakFlg = true
				}
			case token.UNDER_SCORE_BOLD_BEGIN:
				peeked2ndOrderCh := l.peek2ndOrderChar()
				if isUnderScore(peeked2ndOrderCh) {
					peeked3ndOrderCh := l.peek3ndOrderChar()
					if isSpace(peeked3ndOrderCh) || isLineFeedCode(peeked3ndOrderCh) {
						breakFlg = true
					}
				}
			case token.UNDER_SCORE_ITALIC_BOLD_BEGIN:
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
		case isLeftBracket(nextCh):
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

func isLineFeedCode(ch byte) bool {
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

func isLeftBracket(ch byte) bool {
	return ch == '['
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
