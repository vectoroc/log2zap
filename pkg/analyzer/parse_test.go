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
		{"error string without key", "http request: %s", []formatVar{{raw: "%s", value: "%s"}}},
		{"err with quotes", "failed err='%v'", []formatVar{{raw: "err='%v'", key: "err", value: "%v"}}},
		{"err with escaped quotes", `failed err=\"%v\"`, []formatVar{{raw: `err=\"%v\"`, key: "err", value: "%v"}}},
		{"value with few names", "import has failed domain_id=%d err='%v'", []formatVar{
			{raw: "domain_id=%d", key: "domain_id", value: "%d"},
			{raw: "err='%v'", key: "err", value: "%v"},
		}},
		{"longer value placeholder", `failed err=\"%#v\"`, []formatVar{{raw: `err=\"%#v\"`, key: "err", value: "%#v"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.expected, ParseFormat(test.format))
		})
	}
}

func TestParseStaticVars(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected []formatVar
	}{
		{"string without static vars", "some message err=%v", nil},
		{"one static var", "done process=download", []formatVar{{raw: "process=download", key: "process", value: "download"}}},
		{"var with single quotes", "starting task='import'", []formatVar{{raw: "task='import'", key: "task", value: "import"}}},
		{"var with double quotes", `starting task="import"`, []formatVar{{raw: `task="import"`, key: "task", value: "import"}}},
		{"var with escaped quotes", `starting task=\"import\"`, []formatVar{{raw: `task=\"import\"`, key: "task", value: "import"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.expected, ParseStaticVars(test.format))
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
		{"dot ended msg", "task finished.", "task finished"},
		{"dot in the middle", "task finished. exiting", "task finished. exiting"},
		{"err with quotes", "failed err='%v'", "failed"},
		{"value with few names", "import has failed domain_id=%d err='%v'", "import has failed"},
		{"placeholder inside value string", "domain %s imported", "domain imported"},
		{"space before dot", "domain err='%s'. exiting", "domain. exiting"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vars := ParseFormat(test.format)
			assert.EqualValues(t, test.expected, CleanUpFormatString(test.format, vars))
		})
	}
}
