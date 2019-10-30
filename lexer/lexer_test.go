package lexer

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/istsh/markdown-viewer/token"
)

func TestLexer1(t *testing.T) {
	input, err := ioutil.ReadFile("../testdata/1.md.golden")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.HEADING1, "#"},
		{token.SPACE, " "},
		{token.STRING, "Heading1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING2, "##"},
		{token.SPACE, " "},
		{token.STRING, "Heading2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING3, "###"},
		{token.SPACE, " "},
		{token.STRING, "Heading3"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		// t.Logf("tests[%d] - got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if !bytes.Equal(tok.Literal, []byte(tt.expectedLiteral)) {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer2(t *testing.T) {
	input, err := ioutil.ReadFile("../testdata/2.md.golden")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.HEADING1, "#"},
		{token.SPACE, " "},
		{token.STRING, "Heading1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING2, "##"},
		{token.SPACE, " "},
		{token.STRING, "Heading1_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING1, "#"},
		{token.SPACE, " "},
		{token.STRING, "Heading2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING2, "##"},
		{token.SPACE, " "},
		{token.STRING, "Heading2_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB1, "\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB1, "\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1_2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING1, "#"},
		{token.SPACE, " "},
		{token.STRING, "Heading3"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING2, "##"},
		{token.SPACE, " "},
		{token.STRING, "Heading3_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB1, "\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB2, "\t\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1_1_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		// t.Logf("tests[%d] - got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if !bytes.Equal(tok.Literal, []byte(tt.expectedLiteral)) {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer3(t *testing.T) {
	input, err := ioutil.ReadFile("../testdata/3.md.golden")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.HEADING1, "#"},
		{token.SPACE, " "},
		{token.STRING, "Heading1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING2, "##"},
		{token.SPACE, " "},
		{token.STRING, "Heading2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB1, "\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "Nest List1_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB2, "\t\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "Nest List1_1_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB2, "\t\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "Nest List1_1_2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB1, "\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "Nest List1_2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List3"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HEADING3, "###"},
		{token.SPACE, " "},
		{token.STRING, "Heading3"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB1, "\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "Nest List1_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB2, "\t\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "Nest List1_1_1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB2, "\t\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "Nest List1_1_2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.TAB1, "\t"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "Nest List1_2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.STRING, "List3"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.CITATION, ">"},
		{token.SPACE, " "},
		{token.STRING, "Description1"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.CITATION, ">"},
		{token.SPACE, " "},
		{token.STRING, "Description2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.CITATION, ">"},
		{token.SPACE, " "},
		{token.STRING, "Description3_1 > Description3_2"},
		{token.LINE_FEED_CODE_N, "\n"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		// t.Logf("tests[%d] - got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if !bytes.Equal(tok.Literal, []byte(tt.expectedLiteral)) {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
