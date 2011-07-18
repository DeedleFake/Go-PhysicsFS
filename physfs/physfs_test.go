package physfs

import (
	"fmt"
	"testing"
)

func TestSupportedArchiveTypes(t *testing.T) {
	if !IsInit() {
		err := Init()
		if err != nil {
			t.Fatalf("Error: %v\n", err.String())
		}
	}

	sat := SupportedArchiveTypes()
	fmt.Printf("SupportedArchiveTypes: %v:\n", len(sat))
	for i := range sat {
		fmt.Printf("\t%v: %v\n", i+1, sat[i])
	}

	err := Deinit()
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}
}
