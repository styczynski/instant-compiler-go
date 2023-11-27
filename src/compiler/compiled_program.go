package compiler

type CompiledProgram interface {
	ProgramToText() string
}

type CompiledProgramEmpty struct {}

func (CompiledProgramEmpty) ProgramToText() string {
	return ""
}
