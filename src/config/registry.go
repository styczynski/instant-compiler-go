package config

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type EntityType int64

const (
	ENTITY_COMPILER_BACKEND EntityType = iota
	ENTITY_SUMMARIZER
	ENTITY_RUNNER
	ENTITY_STATUS_UPDATER
	ENTITY_COMPILER_PIPELINE
)

type EntityConfig interface {
	String(name string) string
	Int(name string) int
}

type EntityConfigProxy struct {
	prefix string
	c      EntityConfig
}

func (c EntityConfigProxy) String(name string) string {
	return c.c.String(fmt.Sprintf("%s%s", c.prefix, name))
}

func (c EntityConfigProxy) Int(name string) int {
	return c.c.Int(fmt.Sprintf("%s%s", c.prefix, name))
}

type EntityArgSpec struct {
	prefix string
	flags  []cli.Flag
}

func CreateEmptyEntityArgSpec() *EntityArgSpec {
	return &EntityArgSpec{
		prefix: "",
		flags:  []cli.Flag{},
	}
}

func (argSpec *EntityArgSpec) CliFlags() []cli.Flag {
	return argSpec.flags
}

func (argSpec *EntityArgSpec) AddString(name string, value string, usage string) *EntityArgSpec {
	argSpec.flags = append(argSpec.flags, &cli.StringFlag{
		Name:  fmt.Sprintf("%s%s", argSpec.prefix, name),
		Usage: usage,
		Value: value,
	})
	return argSpec
}

func (argSpec *EntityArgSpec) AddInt(name string, value int, usage string) *EntityArgSpec {
	argSpec.flags = append(argSpec.flags, &cli.IntFlag{
		Name:  fmt.Sprintf("%s%s", argSpec.prefix, name),
		Usage: usage,
		Value: value,
	})
	return argSpec
}

type EntityFactory interface {
	CreateEntity(c EntityConfig) interface{}
	EntityName() string
	Params(argSpec *EntityArgSpec)
}

var entityFactoryRegistry = map[EntityType]map[string]EntityFactory{}

func RegisterEntityFactory(t EntityType, factory EntityFactory) {
	if _, ok := entityFactoryRegistry[t]; !ok {
		entityFactoryRegistry[t] = map[string]EntityFactory{}
	}
	entityFactoryRegistry[t][factory.EntityName()] = factory
}

func CreateEntity(t EntityType, name string, c EntityConfig) interface{} {
	if t == ENTITY_COMPILER_PIPELINE {
		return entityFactoryRegistry[t][name].CreateEntity(c)
	}
	return entityFactoryRegistry[t][name].CreateEntity(EntityConfigProxy{
		prefix: fmt.Sprintf("%s-", name),
		c:      c,
	})
}

func GetEntityNamesList(t EntityType) []string {
	names := []string{}
	for _, Entity := range entityFactoryRegistry[t] {
		names = append(names, Entity.EntityName())
	}
	return names
}

func GetEntityParams() []cli.Flag {
	params := []cli.Flag{}
	for t, _ := range entityFactoryRegistry {
		for _, entity := range entityFactoryRegistry[t] {
			argSpec := &EntityArgSpec{
				prefix: fmt.Sprintf("%s-", entity.EntityName()),
				flags:  []cli.Flag{},
			}
			if t == ENTITY_COMPILER_PIPELINE {
				argSpec.prefix = ""
			}
			entity.Params(argSpec)
			params = append(params, argSpec.CliFlags()...)
		}
	}
	return params
}
