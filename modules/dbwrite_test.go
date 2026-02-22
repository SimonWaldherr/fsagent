package modules

import "testing"

func TestDBQuote(t *testing.T) {
	got := dbQuote("a'b")
	if got != "'a''b'" {
		t.Fatalf("unexpected quoted value: %s", got)
	}
}

func TestDBWriteName(t *testing.T) {
	if (DBWrite{}).Name() != "dbwrite" {
		t.Fatalf("unexpected module name")
	}
}
