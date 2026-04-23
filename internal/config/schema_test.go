package config

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// TestConfigSchemaMatchesStruct guards against drift between the
// hand-maintained docs/ggc-config.schema.json and the Config struct.
// If this fails, regenerate the schema (docs/ggc-config.schema.json)
// so the top-level YAML keys match.
func TestConfigSchemaMatchesStruct(t *testing.T) {
	const schemaPath = "../../docs/ggc-config.schema.json"

	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	schemaKeys := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties {
		schemaKeys = append(schemaKeys, k)
	}
	sort.Strings(schemaKeys)

	structKeys := topLevelYAMLKeys(reflect.TypeOf(Config{}))
	sort.Strings(structKeys)

	if !reflect.DeepEqual(schemaKeys, structKeys) {
		t.Fatalf("drift between Config struct and %s:\n  struct keys: %v\n  schema keys: %v\nRun docs update and realign.",
			schemaPath, structKeys, schemaKeys)
	}
}

func topLevelYAMLKeys(t reflect.Type) []string {
	keys := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("yaml")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.SplitN(tag, ",", 2)[0]
		if name == "" {
			continue
		}
		keys = append(keys, name)
	}
	return keys
}
