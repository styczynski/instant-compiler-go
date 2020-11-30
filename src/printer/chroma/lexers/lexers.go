// Package lexers contains the registry of all lexers.
//
// Sub-packages contain lexer implementations.
package lexers

// nolint: golint
import (
	"github.com/alecthomas/chroma"
	"github.com/styczynski/latte-compiler/src/printer/chroma/lexers/internal"
	_ "github.com/styczynski/latte-compiler/src/printer/chroma/lexers/latte"
)

// Registry of Lexers.
var Registry = internal.Registry

// Names of all lexers, optionally including aliases.
func Names(withAliases bool) []string { return internal.Names(withAliases) }

// Get a Lexer by name, alias or file extension.
func Get(name string) chroma.Lexer { return internal.Get(name) }

// MatchMimeType attempts to find a lexer for the given MIME type.
func MatchMimeType(mimeType string) chroma.Lexer { return internal.MatchMimeType(mimeType) }

// Match returns the first lexer matching filename.
func Match(filename string) chroma.Lexer { return internal.Match(filename) }

// Analyse text content and return the "best" lexer..
func Analyse(text string) chroma.Lexer { return internal.Analyse(text) }

// Register a Lexer with the global registry.
func Register(lexer chroma.Lexer) chroma.Lexer { return internal.Register(lexer) }

// Fallback lexer if no other is found.
var Fallback = internal.Fallback
