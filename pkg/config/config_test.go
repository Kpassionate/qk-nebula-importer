package config

import (
	"encoding/json"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/vesoft-inc/nebula-importer/pkg/logger"
)

func TestYAMLParser(t *testing.T) {
	runnerLogger := logger.NewRunnerLogger("")
	yaml, err := Parse("../../examples/example.yaml", runnerLogger)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range yaml.Files {
		if strings.ToLower(*file.Type) != "csv" {
			t.Fatal("Error file type")
		}
		switch strings.ToLower(*file.Schema.Type) {
		case "edge":
		case "vertex":
			if file.Schema.Vertex == nil {
				continue
			}
			if len(file.Schema.Vertex.Tags) == 0 && !*file.CSV.WithHeader {
				t.Fatal("Empty tags in vertex")
			}
		default:
			t.Fatal("Error schema type")
		}
	}
}

type testYAML struct {
	Version *string `yaml:"version"`
	Files   *[]struct {
		Path *string `yaml:"path"`
	} `yaml:"files"`
}

var yamlStr string = `
version: beta
files:
  - path: ./file.csv
`

func TestTypePointer(t *testing.T) {
	ty := testYAML{}
	if err := yaml.Unmarshal([]byte(yamlStr), &ty); err != nil {
		t.Fatal(err)
	}
	t.Logf("yaml: %v, %v", *ty.Version, *ty.Files)
}

var jsonStr = `
{
  "version": "beta",
  "files": [
    { "path": "/tmp" },
    { "path": "/etc" }
	]
}
`

func TestJsonInYAML(t *testing.T) {
	conf := YAMLConfig{}
	if err := yaml.Unmarshal([]byte(jsonStr), &conf); err != nil {
		t.Fatal(err)
	}

	if conf.Version == nil || *conf.Version != "beta" {
		t.Fatal("Error version")
	}

	if conf.Files == nil || len(conf.Files) != 2 {
		t.Fatal("Error files")
	}

	paths := []string{"/tmp", "/etc"}
	for i, p := range paths {
		f := conf.Files[i]
		if f == nil || f.Path == nil || *f.Path != p {
			t.Fatalf("Error file %d path", i)
		}
	}
}

type Person struct {
	Name string `json:"name"`
}

type Man struct {
	Person
	Age int `json:"age"`
}

func TestJsonTypeEmbeding(t *testing.T) {
	man := Man{
		Person: Person{Name: "zhangsan"},
		Age:    18,
	}
	t.Logf("%v", man)
	b, _ := json.Marshal(man)
	t.Logf("%s", string(b))
}
