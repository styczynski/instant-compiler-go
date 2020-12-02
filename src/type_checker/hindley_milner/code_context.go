package hindley_milner

import (
	"fmt"
)

type CodeContext struct {
	Source *Expression
	Builtin bool
	Scheme *Scheme
	Name string
}

func CreateCodeContext(source Expression) CodeContext {
	return CodeContext{
		Source:  &source,
		Builtin: false,
		Scheme: nil,
		Name: "",
	}
}

func CreateBuilinCodeContext(name string, scheme *Scheme) CodeContext {
	return CodeContext{
		Source:  nil,
		Builtin: true,
		Scheme: scheme,
		Name: name,
	}
}

func (c CodeContext) IsEmpty() bool {
	return c.Source == nil && !c.Builtin
}

func (c CodeContext) IsBuiltin() bool {
	return c.Builtin
}

func (c CodeContext) String() string {
	if c.IsEmpty() {
		return "<No info>"
	} else if c.IsBuiltin() {
		return fmt.Sprintf("<internal: %s>", c.Name)
	} else {
		return fmt.Sprintf("<%#v>", (*c.Source))
	}
}


func CopyContextTo(t Type, src ...Type) Type {
	var context *CodeContext = nil
	for _, srcT := range src {
		if !srcT.GetContext().IsEmpty() {
			c := srcT.GetContext()
			context = &c
		}
	}
	if context != nil {
		return t.WithContext(*context)
	}
	return t
}