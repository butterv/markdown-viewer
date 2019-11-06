package token

type TokenType string

// 各文字列が何を意味しているのかを対応づける為に、定数を設けている。
const (
	NONE = "NONE"

	// ILLEGAL = "ILLEGAL" // 解析に失敗した場合に設定するTokenType
	EOF = "EOF" // コードの終了

	//INT    = "INT"
	STRING = "STRING" // 文字列

	// 識別子
	HEADING1 = "HEADING1"
	HEADING2 = "HEADING2"
	HEADING3 = "HEADING3"
	HEADING4 = "HEADING4"
	HEADING5 = "HEADING5"
	HEADING6 = "HEADING6"

	UNORDERED_LIST_BEGIN  = "UNORDERED_LIST_BEGIN"
	UNORDERED_LIST_FINISH = "UNORDERED_LIST_FINISH"
	ORDERED_LIST_BEGIN    = "ORDERED_LIST_BEGIN"
	ORDERED_LIST_FINISH   = "ORDERED_LIST_FINISH"
	LIST_BEGIN            = "LIST_BEGIN"
	LIST_FINISH           = "LIST_FINISH"

	BACK_QUOTE_BEGIN  = "BACK_QUOTE_BEGIN"
	BACK_QUOTE_FINISH = "BACK_QUOTE_FINISH"

	ASTERISK_ITALIC_BEGIN       = "ASTERISK_ITALIC_BEGIN"
	ASTERISK_ITALIC_FINISH      = "ASTERISK_ITALIC_FINISH"
	ASTERISK_BOLD_BEGIN         = "ASTERISK_BOLD_BEGIN"
	ASTERISK_BOLD_FINISH        = "ASTERISK_BOLD_FINISH"
	ASTERISK_ITALIC_BOLD_BEGIN  = "ASTERISK_ITALIC_BOLD_BEGIN"
	ASTERISK_ITALIC_BOLD_FINISH = "ASTERISK_ITALIC_BOLD_FINISH"

	UNDER_SCORE_ITALIC_BEGIN       = "UNDER_SCORE_ITALIC_BEGIN"
	UNDER_SCORE_ITALIC_FINISH      = "UNDER_SCORE_ITALIC_FINISH"
	UNDER_SCORE_BOLD_BEGIN         = "UNDER_SCORE_BOLD_BEGIN"
	UNDER_SCORE_BOLD_FINISH        = "UNDER_SCORE_BOLD_FINISH"
	UNDER_SCORE_ITALIC_BOLD_BEGIN  = "UNDER_SCORE_ITALIC_BOLD_BEGIN"
	UNDER_SCORE_ITALIC_BOLD_FINISH = "UNDER_SCORE_ITALIC_BOLD_FINISH"

	CITATION1 = "CITATION1"
	CITATION2 = "CITATION2"

	HORIZON = "HORIZON"

	LINK_TEXT_BEGIN  = "LINK_TEXT_BEGIN"
	LINK_TEXT_FINISH = "LINK_TEXT_FINISH"
	LINK_BEGIN       = "LINK_BEGIN"
	LINK_FINISH      = "LINK_FINISH"

	//

	TAB1 = "TAB1"
	TAB2 = "TAB2"
	TAB3 = "TAB3"

	SPACE = "SPACE"
	//SPACE4 = "    "

	LINE_FEED_CODE = "LINE_FEED_CODE"

	HYPHEN = "HYPHEN"

	//ASTERISK_ITALIC         = "*"
	//ASTERISK_BOLD           = "**"
	//ASTERISK_ITALIC_BOLD    = "***"
	//UNDER_SCORE_ITALIC      = "_"
	//UNDER_SCORE_BOLD        = "__"
	//UNDER_SCORE_ITALIC_BOLD = "___"

	//DOT         = "."
	//ASTERISK    = "*"
	//UNDER_SCORE = "_"
	//
	//PLUS = "+"
	//
	//LPAREN   = "LPAREN"
	//RPAREN   = "RPAREN"
	//LBRACKET = "LBRACKET"
	//RBRACKET = "RBRACKET"
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
		//return ILLEGAL
		return STRING
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
		//return ILLEGAL
		return STRING
	}
}

func GetAsteriskToken(cnt int) TokenType {
	switch cnt {
	case 1:
		return ASTERISK_ITALIC_BEGIN
	case 2:
		return ASTERISK_BOLD_BEGIN
	case 3:
		return ASTERISK_ITALIC_BOLD_BEGIN
	default:
		//return ILLEGAL
		return STRING
	}
}

func GetUnderScoreToken(cnt int) TokenType {
	switch cnt {
	case 1:
		return UNDER_SCORE_ITALIC_BEGIN
	case 2:
		return UNDER_SCORE_BOLD_BEGIN
	case 3:
		return UNDER_SCORE_ITALIC_BOLD_BEGIN
	default:
		//return ILLEGAL
		return STRING
	}
}

func GetCitationToken(cnt int) TokenType {
	switch cnt {
	case 1:
		return CITATION1
	case 2:
		return CITATION2
	default:
		//return ILLEGAL
		return STRING
	}
}
