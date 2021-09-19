package inmemorycache

import (
	"testing"
)

var initData = []struct {
	key      string
	value    string
	expected string
}{
	{
		key:      "abc",
		value:    "1",
		expected: "1",
	},
	{
		key:      "abc-1",
		value:    "2",
		expected: "2",
	},
	{
		key:      "xyz",
		value:    "3",
		expected: "3",
	},
	{
		key:      "xyz-1",
		value:    "4",
		expected: "4",
	},
	{
		key:      "ro",
		value:    "5",
		expected: "5",
	},
}

var prefixMatchData = []struct {
	name          string
	prefix        string
	total         int
	expectedKeys  []string
	expectedCount int
}{
	{
		name:          "PrefixLimited",
		prefix:        "abc",
		total:         1,
		expectedKeys:  []string{"abc"},
		expectedCount: 1,
	},
	{
		name:          "PrefixAllKeys",
		prefix:        "xyz",
		total:         -1,
		expectedKeys:  []string{"xyz", "xyz-1"},
		expectedCount: 2,
	},
	{
		name:          "PrefixNoMatch",
		prefix:        "c",
		total:         1,
		expectedKeys:  []string{},
		expectedCount: 0,
	},
	{
		name:          "PrefixMatchTerminal",
		prefix:        "abc-1",
		total:         -1,
		expectedKeys:  []string{"abc-1"},
		expectedCount: 1,
	},
	{
		name:          "PrefixMatchTerminalLenOne",
		prefix:        "ro",
		total:         -1,
		expectedKeys:  []string{"ro"},
		expectedCount: 1,
	},
}

func TestInsert(t *testing.T) {
	tr := New()
	for _, v := range initData {
		tr.Set(v.key, v.value)
		val, ok := tr.Get(v.key)
		if !ok {
			t.Fatalf("Expected %v, got %v", true, ok)
		}
		if v.expected != val {
			t.Fatalf("Expected %v, got %v", v.expected, val)
		}
	}
}

func TestGetNotFound(t *testing.T) {
	// Not Found Key check
	tr := New()
	v, ok := tr.Get("keyword")
	if ok != false {
		t.Fatalf("Expected %v got %v", false, ok)
	}
	if v != nil {
		t.Fatalf("Expected %v got %v", nil, v)
	}
}

func TestGetSuggestion(t *testing.T) {
	tr := New()
	for _, v := range initData {
		tr.Set(v.key, v.value)
	}
	for _, val := range prefixMatchData {
		resKeys := tr.PrefixMatch(val.prefix, val.total)
		if len(resKeys) != val.expectedCount {
			t.Fatalf("expected %v, got %v", val.expectedCount, len(resKeys))
		}
		for k, v := range resKeys {
			if v != val.expectedKeys[k] {
				t.Fatalf("expected %v, got %v", val.expectedKeys[k], v)
			}
		}
	}
}
