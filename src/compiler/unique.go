package compiler

import "fmt"

type UniqueNameGenerator interface {
	Next() string
	Copy() UniqueNameGenerator
}

type SeqNameGenerator struct {
	nameID int64
}

func (gen *SeqNameGenerator) Copy() UniqueNameGenerator {
	return &SeqNameGenerator{
		nameID: gen.nameID,
	}
}

func (gen *SeqNameGenerator) Next() string {
	id := gen.nameID
	gen.nameID++
	return fmt.Sprintf("_var_%d", id)
}

func CreateUniqueNameGenerator() UniqueNameGenerator {
	return &SeqNameGenerator{
		nameID: 1,
	}
}
