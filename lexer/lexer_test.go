package lexer

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/istsh/markdown-viewer/token"
)

type expected struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func compareGotAndWant(t *testing.T, goldenPath string, tests []expected) {
	input, err := ioutil.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}

	l := New(input)

	var tokens []token.Token
	//var i int
	for {
		//i++
		//fmt.Printf("[%d]\n", i)

		toks := l.NextTokens()
		tokens = append(tokens, toks...)

		//fmt.Printf("%v\n", toks)
		if toks[0].Type == token.EOF {
			break
		}
	}

	if len(tokens) != len(tests) {
		t.Errorf("not match length. len(gots)=%d, len(wants)=%d", len(tokens), len(tests))
	}

	for i, tt := range tests {
		tok := tokens[i]

		//if tt.expectedType == token.STRING {
		//	t.Logf("tests[%d] - got=%d(%q), want=%d(%q)", i, tok.Type, tok.Literal, tt.expectedType, tt.expectedLiteral)
		//} else {
		//	t.Logf("tests[%d] - got=%d, want=%d", i, tok.Type, tt.expectedType)
		//}

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%d, got=%d", i, tt.expectedType, tok.Type)
		}
		if tt.expectedType == token.STRING && !bytes.Equal(tok.Literal, []byte(tt.expectedLiteral)) {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer1(t *testing.T) {
	tests := []expected{
		{expectedType: token.HEADING1},
		{expectedType: token.STRING, expectedLiteral: "Heading1"},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.HEADING2},
		{expectedType: token.STRING, expectedLiteral: "Heading2"},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.HEADING3},
		{expectedType: token.STRING, expectedLiteral: "Heading3"},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.HEADING1},
		{expectedType: token.STRING, expectedLiteral: "Heading1"},
		{expectedType: token.SPACE},
		{expectedType: token.STRING, expectedLiteral: "Text"},
		{expectedType: token.SPACE},
		{expectedType: token.STRING, expectedLiteral: "#"},
		{expectedType: token.SPACE},
		{expectedType: token.STRING, expectedLiteral: "Heading1"},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.STRING, expectedLiteral: "#Heading1"},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.STRING, expectedLiteral: "##Heading2"},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.STRING, expectedLiteral: "###Heading3"},
		{expectedType: token.LINE_FEED_CODE},
		{expectedType: token.EOF},
	}

	compareGotAndWant(t, "../testdata/1.md.golden", tests)
}

//func TestLexer2(t *testing.T) {
//	tests := []expected{
//		{expectedType: token.HEADING1},
//		{expectedType: token.STRING, expectedLiteral: "Heading1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HEADING2},
//		{expectedType: token.STRING, expectedLiteral: "Heading1_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HEADING1},
//		{expectedType: token.STRING, expectedLiteral: "Heading2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HEADING2},
//		{expectedType: token.STRING, expectedLiteral: "Heading2_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB1},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB1},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HEADING1},
//		{expectedType: token.STRING, expectedLiteral: "Heading3"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HEADING2},
//		{expectedType: token.STRING, expectedLiteral: "Heading3_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB1},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB2},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.EOF},
//	}
//
//	compareGotAndWant(t, "../testdata/2.md.golden", tests)
//}
//
//func TestLexer3(t *testing.T) {
//	tests := []expected{
//		{expectedType: token.HEADING1},
//		{expectedType: token.STRING, expectedLiteral: "Heading1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HEADING2},
//		{expectedType: token.STRING, expectedLiteral: "Heading2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB1},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Nest"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB2},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Nest"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB2},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Nest"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB1},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Nest"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List3"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HEADING3},
//		{expectedType: token.STRING, expectedLiteral: "Heading3"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB1},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Nest"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB2},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Nest"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1_1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB2},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Nest"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_1_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.TAB1},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Nest"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List1_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HYPHEN},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "List3"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.CITATION1},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.CITATION1},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.CITATION1},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description3_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: ">"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description3_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.CITATION2},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description4"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.CITATION2},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description5"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.EOF},
//	}
//
//	compareGotAndWant(t, "../testdata/3.md.golden", tests)
//}
//
//func TestLexer4(t *testing.T) {
//	tests := []expected{
//		{expectedType: token.HEADING1},
//		{expectedType: token.STRING, expectedLiteral: "Heading1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.STRING, expectedLiteral: "Description1_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.BACK_QUOTE_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "back"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "quote"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.BACK_QUOTE_FINISH},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description1_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.STRING, expectedLiteral: "Description2_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.ASTERISK_ITALIC_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "italic"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.ASTERISK_ITALIC_FINISH},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description2_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.STRING, expectedLiteral: "Description3_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.ASTERISK_BOLD_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "bold"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.ASTERISK_BOLD_FINISH},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description3_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.STRING, expectedLiteral: "Description4_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.ASTERISK_ITALIC_BOLD_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "italic"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "&"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "bold"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.ASTERISK_ITALIC_BOLD_FINISH},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description4_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.ASTERISK_ITALIC_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "italic"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.ASTERISK_ITALIC_FINISH},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.ASTERISK_BOLD_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "bold"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.ASTERISK_BOLD_FINISH},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.ASTERISK_ITALIC_BOLD_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "italic"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "&"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "bold"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.ASTERISK_ITALIC_BOLD_FINISH},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.EOF},
//	}
//
//	compareGotAndWant(t, "../testdata/4.md.golden", tests)
//}
//
//func TestLexer5(t *testing.T) {
//	tests := []expected{
//		{expectedType: token.HEADING1},
//		{expectedType: token.STRING, expectedLiteral: "Heading1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.STRING, expectedLiteral: "Description1_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.BACK_QUOTE_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "back"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "quote"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.BACK_QUOTE_FINISH},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description1_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.STRING, expectedLiteral: "Description2_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.UNDER_SCORE_ITALIC_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "italic"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.UNDER_SCORE_ITALIC_FINISH},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description2_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.STRING, expectedLiteral: "Description3_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.UNDER_SCORE_BOLD_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "bold"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.UNDER_SCORE_BOLD_FINISH},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description3_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.STRING, expectedLiteral: "Description4_1"},
//		{expectedType: token.SPACE},
//		{expectedType: token.UNDER_SCORE_ITALIC_BOLD_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "italic"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "&"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "bold"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.UNDER_SCORE_ITALIC_BOLD_FINISH},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "Description4_2"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.UNDER_SCORE_ITALIC_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "italic"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.UNDER_SCORE_ITALIC_FINISH},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.UNDER_SCORE_BOLD_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "bold"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.UNDER_SCORE_BOLD_FINISH},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.UNDER_SCORE_ITALIC_BOLD_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "italic"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "&"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "bold"},
//		{expectedType: token.SPACE},
//		{expectedType: token.STRING, expectedLiteral: "area"},
//		{expectedType: token.UNDER_SCORE_ITALIC_BOLD_FINISH},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.EOF},
//	}
//
//	compareGotAndWant(t, "../testdata/5.md.golden", tests)
//}
//
//func TestLexer6(t *testing.T) {
//	tests := []expected{
//		{expectedType: token.HEADING1},
//		{expectedType: token.STRING, expectedLiteral: "Heading1"},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HORIZON},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HORIZON},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.HORIZON},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.LINK_TEXT_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "Google"},
//		{expectedType: token.LINK_TEXT_FINISH},
//		{expectedType: token.LINK_BEGIN},
//		{expectedType: token.STRING, expectedLiteral: "https://www.google.com/"},
//		{expectedType: token.LINK_FINISH},
//		{expectedType: token.LINE_FEED_CODE},
//		{expectedType: token.EOF},
//	}
//
//	compareGotAndWant(t, "../testdata/6.md.golden", tests)
//}
