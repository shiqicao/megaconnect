// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"fmt"
	"io/ioutil"
	"unicode/utf8"

	"github.com/megaspacelab/megaconnect/workflow/parser/gen/util"
	"github.com/megaspacelab/megaconnect/workflow/parser/gen/token"
)

const (
	NoState    = -1
	NumStates  = 110
	NumSymbols = 128
)

type Lexer struct {
	src    []byte
	pos    int
	line   int
	column int
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:    src,
		pos:    0,
		line:   1,
		column: 1,
	}
	return lexer
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewLexer(src), nil
}

func (l *Lexer) Scan() (tok *token.Token) {
	fmt.Printf("Lexer.Scan() pos=%d\n", l.pos)
	tok = new(token.Token)
	if l.pos >= len(l.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = l.pos, l.line, l.column
		return
	}
	start, startLine, startColumn, end := l.pos, l.line, l.column, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
		fmt.Printf("\tpos=%d, line=%d, col=%d, state=%d\n", l.pos, l.line, l.column, state)
		if l.pos >= len(l.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(l.src[l.pos:])
			l.pos += size
		}

		nextState := -1
		if rune1 != -1 {
			nextState = TransTab[state](rune1)
		}
		fmt.Printf("\tS%d, : tok=%s, rune == %s(%x), next state == %d\n", state, token.TokMap.Id(tok.Type), util.RuneToString(rune1), rune1, nextState)
		fmt.Printf("\t\tpos=%d, size=%d, start=%d, end=%d\n", l.pos, size, start, end)
		if nextState != -1 {
			fmt.Printf("\t\taction:%s\n", ActTab[nextState].String())
		}
		state = nextState

		if state != -1 {

			switch rune1 {
			case '\n':
				l.line++
				l.column = 1
			case '\r':
				l.column = 1
			case '\t':
				l.column += 4
			default:
				l.column++
			}

			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				end = l.pos
			case ActTab[state].Ignore != "":
				start, startLine, startColumn = l.pos, l.line, l.column
				state = 0
				if start >= len(l.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = l.pos
			}
		}
	}
	if end > start {
		l.pos = end
		tok.Lit = l.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = start, startLine, startColumn
	fmt.Printf("Token at %s: %s \"%s\"\n", tok.String(), token.TokMap.Id(tok.Type), tok.Lit)

	return
}

func (l *Lexer) Reset() {
	l.pos = 0
}

/*
Lexer symbols:
0: 'm'
1: 'o'
2: 'n'
3: 'i'
4: 't'
5: 'o'
6: 'r'
7: 'e'
8: 'v'
9: 'e'
10: 'n'
11: 't'
12: 'a'
13: 'c'
14: 't'
15: 'i'
16: 'o'
17: 'n'
18: 'c'
19: 'h'
20: 'a'
21: 'i'
22: 'n'
23: 'c'
24: 'o'
25: 'n'
26: 'd'
27: 'i'
28: 't'
29: 'i'
30: 'o'
31: 'n'
32: 'v'
33: 'a'
34: 'r'
35: 'w'
36: 'o'
37: 'r'
38: 'k'
39: 'f'
40: 'l'
41: 'o'
42: 'w'
43: 'f'
44: 'i'
45: 'r'
46: 'e'
47: 'i'
48: 'n'
49: 't'
50: 's'
51: 't'
52: 'r'
53: 'i'
54: 'n'
55: 'g'
56: 'b'
57: 'o'
58: 'o'
59: 'l'
60: 'r'
61: 'u'
62: 'n'
63: 't'
64: 'r'
65: 'i'
66: 'g'
67: 'g'
68: 'e'
69: 'r'
70: 't'
71: 'r'
72: 'u'
73: 'e'
74: 'f'
75: 'a'
76: 'l'
77: 's'
78: 'e'
79: '"'
80: '\'
81: '"'
82: '"'
83: '0'
84: '\n'
85: '{'
86: '}'
87: ':'
88: ','
89: '|'
90: '|'
91: '&'
92: '&'
93: '='
94: '='
95: '='
96: '!'
97: '='
98: '>'
99: '>'
100: '='
101: '<'
102: '<'
103: '='
104: '+'
105: '-'
106: '*'
107: '/'
108: '.'
109: '('
110: ')'
111: '\'
112: 'n'
113: '\'
114: 'r'
115: '\'
116: 't'
117: ' '
118: '\t'
119: '\r'
120: '\n'
121: '/'
122: '/'
123: '1'-'9'
124: '0'-'9'
125: 'a'-'z'
126: 'A'-'Z'
127: .
*/
