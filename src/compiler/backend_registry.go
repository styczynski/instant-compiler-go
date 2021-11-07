package compiler

type CompilerBackendFactory interface {
	CreateBackend() CompilerBackend
	BackendName() string
}

var backendFactoryRegistry = map[string]CompilerBackendFactory{}

func RegisterCompilerBackendFactory(factory CompilerBackendFactory) {
	backendFactoryRegistry[factory.BackendName()] = factory
}

func CreateCompilerBackend(name string) CompilerBackend {
	return backendFactoryRegistry[name].CreateBackend()
}
