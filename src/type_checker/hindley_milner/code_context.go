package hindley_milner

type CodeContext struct {
	Source *Expression
	Builtin bool
}

func (c CodeContext) IsEmpty() bool {
	return c.Source == nil
}

func (c CodeContext) IsBuiltin() bool {
	return c.Builtin
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