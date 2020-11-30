package parser

import "github.com/alecthomas/participle/v2/lexer/stateful"

// A custom lexer for LatteProgram files. This illustrates a relatively complex Regexp lexer, as well
// as use of the Unquote filter, which unquotes string tokens.
var iniLexer = stateful.MustSimple([]stateful.Rule{
	{"TypeName", `string|int`, nil},
	{"Comment", `(?i)//[^\n]*`, nil},
	{"String", `"(\\"|[^"])*"`, nil},
	{"Int", `[-+]?(\d*\.)?\d+`, nil},
	{"Ident", `[a-zA-Z_]\w*`, nil},
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
	{"EOL", `[\n\r]+`, nil},
	{"whitespace", `[ \t]+`, nil},
})