package token

type TokenType string

// 各文字列が何を意味しているのかを対応づける為に、定数を設けている。
const (
	ILLEGAL = "ILLEGAL" // 解析に失敗した場合に設定するTokenType
	EOF     = "EOF"     // コードの終了

	//INT    = "INT"
	STRING = "STRING" // 文字列

	// 識別子
	HEADING1 = "#"
	HEADING2 = "##"
	HEADING3 = "###"
	HEADING4 = "####"
	HEADING5 = "#####"
	HEADING6 = "######"

	TAB1 = "\t"
	TAB2 = "\t\t"
	TAB3 = "\t\t\t"

	SPACE = " "
	//SPACE4 = "    "

	LINE_FEED_CODE_N = "\n"
	LINE_FEED_CODE_R = "\r"

	CITATION = ">"
	HYPHEN   = "-"

	BACK_QUOTE = "`"

	//DOT         = "."
	//ASTERISK    = "*"
	//UNDER_SCORE = "_"
	//
	//PLUS = "+"
	//
	//LPAREN   = "("
	//RPAREN   = ")"
	//LBRACKET = "["
	//RBRACKET = "]"
	//
	//// ここから下は不要かも
	//ASSIGN = "="
	//
	//BANG = "!"
	//
	//SLASH = "/"
	//
	//GT = ">"
	//LT = "<"
	//
	//EQ     = "=="
	//NOT_EQ = "!="
	//
	//// デリミタ
	//COMMA     = ","
	//SEMICOLON = ";"
	//COLON     = ":"
	//
	//LBRACE = "{"
	//RBRACE = "}"
	//
	//// キーワード
	//FUNCTION = "FUNCTION"
	//LET      = "LET"
	//TRUE     = "TRUE"
	//FALSE    = "FALSE"
	//IF       = "IF"
	//ELSE     = "ELSE"
	//RETURN   = "RETURN"
	//MACRO    = "MACRO"
)

type Token struct {
	Type    TokenType // トークンタイプ
	Literal []byte    // トークンリテラル、Typeに応じたbyte配列
}

func GetHeadingToken(cnt int) TokenType {
	switch cnt {
	case 1:
		return HEADING1
	case 2:
		return HEADING2
	case 3:
		return HEADING3
	case 4:
		return HEADING4
	case 5:
		return HEADING5
	case 6:
		return HEADING6
	default:
		return ILLEGAL
	}
}

func GetTabToken(cnt int) TokenType {
	switch cnt {
	case 1:
		return TAB1
	case 2:
		return TAB2
	case 3:
		return TAB3
	default:
		return ILLEGAL
	}
}
