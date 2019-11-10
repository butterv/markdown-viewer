package parser

import (
	"fmt"

	"github.com/istsh/markdown-viewer/lexer"
	"github.com/istsh/markdown-viewer/token"
)

// Parser has a lexer pointer.
type Parser struct {
	l *lexer.Lexer
	//errors []string
}

// New initializes Parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		//errors: []string{},
	}

	return p
}

// Parse parses markdown text to html text.
func (p *Parser) Parse() []byte {
	var result []byte

	var pendingToken token.TokenType
	for {
		currentToken := p.l.NextToken()
		if currentToken.Type == token.EOF {
			break
		}

		switch currentToken.Type {
		case token.HEADING1:
			pendingToken = token.HEADING1
			result = appendStr(result, "<h1>")
		case token.HEADING2:
			pendingToken = token.HEADING2
			result = appendStr(result, "<h2>")
		case token.HEADING3:
			pendingToken = token.HEADING3
			result = appendStr(result, "<h3>")
		case token.HEADING4:
			pendingToken = token.HEADING4
			result = appendStr(result, "<h4>")
		case token.HEADING5:
			pendingToken = token.HEADING5
			result = appendStr(result, "<h5>")
		case token.HEADING6:
			pendingToken = token.HEADING6
			result = appendStr(result, "<h6>")
		case token.BACK_QUOTE_BEGIN:
			result = appendStr(result, "<back_quote>")
		case token.BACK_QUOTE_FINISH:
			result = appendStr(result, "</back_quote>")
		case token.ASTERISK_ITALIC_BEGIN, token.UNDER_SCORE_ITALIC_BEGIN:
			result = appendStr(result, "<italic>")
		case token.ASTERISK_ITALIC_FINISH, token.UNDER_SCORE_ITALIC_FINISH:
			result = appendStr(result, "</italic>")
		case token.ASTERISK_BOLD_BEGIN, token.UNDER_SCORE_BOLD_BEGIN:
			result = appendStr(result, "<bold>")
		case token.ASTERISK_BOLD_FINISH, token.UNDER_SCORE_BOLD_FINISH:
			result = appendStr(result, "</bold>")
		case token.ASTERISK_ITALIC_BOLD_BEGIN, token.UNDER_SCORE_ITALIC_BOLD_BEGIN:
			result = appendStr(result, "<italic_bold>")
		case token.ASTERISK_ITALIC_BOLD_FINISH, token.UNDER_SCORE_ITALIC_BOLD_FINISH:
			result = appendStr(result, "</italic_bold>")
		case token.CITATION1:
			result = appendStr(result, "<citation1>")
		case token.CITATION2:
			result = appendStr(result, "<citation2>")
		case token.SPACE:
			result = appendStr(result, "<space>")
		case token.HYPHEN:
			result = appendStr(result, "<hyphen>")
		case token.STRING:
			result = append(result, currentToken.Literal...)
		case token.TAB1:
			result = appendStr(result, "\t")
		case token.TAB2:
			result = appendStr(result, "\t\t")
		case token.TAB3:
			result = appendStr(result, "\t\t\t")
		case token.HORIZON:
			result = appendStr(result, "<hr>")
		case token.LINK_TEXT_BEGIN:
			result = appendStr(result, "<a link=\"\">")
		case token.LINK_TEXT_FINISH:
			result = appendStr(result, "</a>")
		case token.LINK_BEGIN:
			result = append(result, currentToken.Literal...)
		case token.LINK_FINISH:
			result = append(result, currentToken.Literal...)
		case token.LINE_FEED_CODE:
			switch pendingToken {
			case token.HEADING1:
				result = appendStr(result, "</h1>")
			case token.HEADING2:
				result = appendStr(result, "</h2>")
			case token.HEADING3:
				result = appendStr(result, "</h3>")
			case token.HEADING4:
				result = appendStr(result, "</h4>")
			case token.HEADING5:
				result = appendStr(result, "</h5>")
			case token.HEADING6:
				result = appendStr(result, "</h6>")
			}
			pendingToken = token.NONE
			result = appendStr(result, "\n")
		default:
			panic(fmt.Sprintf("unsupported token type: %q", currentToken.Type))
		}
	}

	return result
}

func appendStr(slice []byte, str string) []byte {
	return append(slice, []byte(str)...)
}
