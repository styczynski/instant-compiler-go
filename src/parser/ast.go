package parser

type LatteProgram struct {
	TopDef *FnDef `@@`
}

type FnDef struct {
	ReturnType Type `@@`
	Name string `@Ident`
	Arg []*Arg `"(" @@ ( "," @@ )* ")"`
	Body Block `@@`
}

type Block struct {
	Statements []*Expression `"{" ( @@ ";" )* "}"`
}

type Arg struct {
	ArgumentType Type `@@`
	Name string `@Ident`
}

type Type struct {
	Name string `@TypeName`
}

type TopDef struct {
	FnDef `@Type @Ident "(" Arg* ")" Block`
}

type Expression struct {
	Equality *Equality `@@`
}

type Equality struct {
	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
}

type Comparison struct {
	Addition *Addition   `@@`
	Op       string      `[ @( ">" | ">" "=" | "<" | "<" "=" )`
	Next     *Comparison `  @@ ]`
}

type Addition struct {
	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
}

type Multiplication struct {
	Unary *Unary          `@@`
	Op    string          `[ @( "/" | "*" )`
	Next  *Multiplication `  @@ ]`
}

type Unary struct {
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	Primary *Primary `| @@`
}

type Primary struct {
	Number        *int64      `@Int`
	String        *string     `| @String`
	Bool          *bool       `| ( @"true" | "false" )`
	Nil           bool        `| @"nil"`
	SubExpression *Expression `| "(" @@ ")" `
}