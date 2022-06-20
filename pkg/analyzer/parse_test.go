package analyzer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected []formatVar
	}{
		{"string without placeholders", "string without placeholders", nil},
		{"error string without name", "http request: %s", []formatVar{{raw: "%s", format: "%s"}}},
		{"err with quotes", "failed err='%v'", []formatVar{{raw: "err='%v'", name: "err", format: "%v"}}},
		{"err with escaped quotes", `failed err=\"%v\"`, []formatVar{{raw: `err=\"%v\"`, name: "err", format: "%v"}}},
		{"format with few names", "import has failed domain_id=%d err='%v'", []formatVar{
			{raw: "domain_id=%d", name: "domain_id", format: "%d"},
			{raw: "err='%v'", name: "err", format: "%v"},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.expected, parseFormat(test.format))
		})
	}
}

func TestCleanUpFormatString(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{"string without placeholders", "string without placeholders", "string without placeholders"},
		{"err with quotes", "failed err='%v'", "failed"},
		{"format with few names", "import has failed domain_id=%d err='%v'", "import has failed"},
		{"placeholder inside format string", "domain %s imported", "domain imported"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vars := parseFormat(test.format)
			assert.EqualValues(t, test.expected, cleanUpFormatString(test.format, vars))
		})
	}
}
