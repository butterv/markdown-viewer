package lexer

import (
	"testing"

	"github.com/istsh/markdown-viewer/token"
)

func Test(t *testing.T) {
	input := `# Heading1
## Heading2
- List1
	- Nest List1_1
		- Nest List1_1_1
		- Nest List1_1_2
	- Nest List1_2
- List2
- List3

### Heading3
1. NumberedList1
	1. NumberedList1_1
	2. NumberedList1_2
2. NumberedList2
3. NumberedList3
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.HEADING1, "#"},
		{token.STRING, "Heading1"},
		{token.HEADING2, "##"},
		{token.STRING, "Heading2"},
		{token.MINUS, "-"},
		{token.STRING, "List1"},
		{token.TAB1, "\t"},
		{token.MINUS, "-"},
		{token.STRING, "Nest List1_1"},
		{token.TAB2, "\t\t"},
		{token.MINUS, "-"},
		{token.STRING, "Nest List1_1_1"},
		{token.TAB2, "\t\t"},
		{token.MINUS, "-"},
		{token.STRING, "Nest List1_1_2"},
		{token.TAB1, "\t"},
		{token.MINUS, "-"},
		{token.STRING, "Nest List1_2"},
		{token.MINUS, "-"},
		{token.STRING, "List2"},
		{token.MINUS, "-"},
		{token.STRING, "List3"},
		{token.HEADING3, "###"},
		{token.STRING, "Heading3"},
		{token.INT, "1"},
		{token.DOT, "."},
		{token.STRING, "NumberedList1"},
		{token.TAB1, "\t"},
		{token.INT, "1"},
		{token.DOT, "."},
		{token.STRING, "NumberedList1_1"},
		{token.TAB1, "\t"},
		{token.INT, "2"},
		{token.DOT, "."},
		{token.STRING, "NumberedList1_2"},
		{token.INT, "2"},
		{token.DOT, "."},
		{token.STRING, "NumberedList2"},
		{token.INT, "3"},
		{token.DOT, "."},
		{token.STRING, "NumberedList3"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer1(t *testing.T) {
	input := `# Heading1

## Heading2

### Heading3
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.HEADING1, "#"},
		{token.STRING, "Heading1"},
		{token.HEADING2, "##"},
		{token.STRING, "Heading2"},
		{token.HEADING3, "###"},
		{token.STRING, "Heading3"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer2(t *testing.T) {
	input := `# Heading1
## Heading1_1
- List1
- List2
# Heading2
## Heading2_1
- List1
	- List1_1
	- List1_2
# Heading3
## Heading3_1
- List1
	- List1_1
		- List1_1_1
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.HEADING1, "#"},
		{token.STRING, "Heading1"},
		{token.HEADING2, "##"},
		{token.STRING, "Heading1_1"},
		{token.MINUS, "-"},
		{token.STRING, "List1"},
		{token.MINUS, "-"},
		{token.STRING, "List2"},
		{token.HEADING1, "#"},
		{token.STRING, "Heading2"},
		{token.HEADING2, "##"},
		{token.STRING, "Heading2_1"},
		{token.MINUS, "-"},
		{token.STRING, "List1"},
		{token.TAB1, "\t"},
		{token.MINUS, "-"},
		{token.STRING, "List1_1"},
		{token.TAB1, "\t"},
		{token.MINUS, "-"},
		{token.STRING, "List1_2"},
		{token.HEADING1, "#"},
		{token.STRING, "Heading3"},
		{token.HEADING2, "##"},
		{token.STRING, "Heading3_1"},
		{token.MINUS, "-"},
		{token.STRING, "List1"},
		{token.TAB1, "\t"},
		{token.MINUS, "-"},
		{token.STRING, "List1_1"},
		{token.TAB2, "\t\t"},
		{token.MINUS, "-"},
		{token.STRING, "List1_1_1"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		t.Logf("tests[%d] - got=%q,%q, want=%q,%q", i, tok.Type, tok.Literal, tt.expectedType, tt.expectedLiteral)

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
