package config

import "testing"

func TestDefaultsAreRepoRelative(t *testing.T) {
	paths := Defaults()
	for label, value := range map[string]string{
		"config":       paths.Config,
		"prd":          paths.PRD,
		"factory_root": paths.FactoryRoot,
	} {
		if value == "" {
			t.Fatalf("%s path is empty", label)
		}
		if value[0] == '/' {
			t.Fatalf("%s path should be repo-relative: %s", label, value)
		}
	}
}
