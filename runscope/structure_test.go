package runscope

import (
	"github.com/hashicorp/terraform/flatmap"
	"reflect"
	"testing"
)

func TestExpandStringList(t *testing.T) {
	expanded := flatmap.Expand(testConf(), "scripts").([]interface{})
	stringList := expandStringList(expanded)
	expected := []string{
		"log(\"hello 1\");",
		"log(\"hello 2\");",
	}

	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}

func testConf() map[string]string {
	return map[string]string{
		"scripts.#": "2",
		"scripts.0": "log(\"hello 1\");",
		"scripts.1": "log(\"hello 2\");",
	}
}
