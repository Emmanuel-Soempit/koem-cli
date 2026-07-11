package impls_test

import (
	"testing"

	"github.com/Emmanuel-Soempit/koem/impls"
)

func TestAddLabel(t *testing.T) {
	label := &impls.Label{}
	err := label.AddLabel("test", []string{"105", "106"})
	if err != nil {
		t.Errorf("AddLabel() error = %v", err)
	}
}
