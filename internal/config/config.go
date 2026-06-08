package config

const DefaultConfigPath = "lumyn.yaml"

type Paths struct {
	Config      string
	PRD         string
	FactoryRoot string
}

func Defaults() Paths {
	return Paths{
		Config:      DefaultConfigPath,
		PRD:         "docs/product/prd.md",
		FactoryRoot: ".factory/artifacts",
	}
}
