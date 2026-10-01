package trace

import (
	"context"
	"strings"
	"testing"
)

func TestNormalizeAndContextIsolation(t *testing.T) {
	if Normalize("incoming-42") != "incoming-42" {
		t.Fatal("safe ID changed")
	}
	for _, invalid := range []string{"", strings.Repeat("a", 129), "header\ninjection"} {
		if id := Normalize(invalid); id == invalid || len(id) == 0 || len(id) > 128 {
			t.Fatal(id)
		}
	}
	root := context.Background()
	ctx := WithID(root, "one")
	if ID(root) != "" || ID(ctx) != "one" {
		t.Fatal("context leaked")
	}
}
