package config

type EntityConstConfig struct {
	Strings map[string]string
	Ints    map[string]int
}

func (c EntityConstConfig) String(name string) string {
	return c.Strings[name]
}

func (c EntityConstConfig) Int(name string) int {
	return c.Ints[name]
}
