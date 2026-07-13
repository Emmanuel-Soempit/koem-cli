package impls_test

import (
	"testing"

	"github.com/Emmanuel-Soempit/koem-cli/impls"
	"github.com/spf13/viper"
)

func TestAddLabel(t *testing.T) {
	label := &impls.Label{}
	err := label.AddLabel("test", []string{"3000", "5000"})
	if err != nil {
		t.Errorf("AddLabel() error = %v", err)
	}
	if label.Min != "3000" {
		t.Errorf("expected Min = 3000, got %s", label.Min)
	}
	if label.Max != "5000" {
		t.Errorf("expected Max = 5000, got %s", label.Max)
	}
}

func TestAddLabel_MinGreaterThanMax(t *testing.T) {
	label := &impls.Label{}
	err := label.AddLabel("test", []string{"5000", "3000"})
	if err == nil {
		t.Error("expected error when min >= max, got nil")
	}
}

func TestAddLabel_InvalidPort(t *testing.T) {
	label := &impls.Label{}
	err := label.AddLabel("test", []string{"abc", "3000"})
	if err == nil {
		t.Error("expected error for non-numeric port, got nil")
	}
}

func TestAddLabel_TooFewPorts(t *testing.T) {
	label := &impls.Label{}
	err := label.AddLabel("test", []string{"3000"})
	if err == nil {
		t.Error("expected error for fewer than 2 ports, got nil")
	}
}

func TestAddLabel_TooManyPorts(t *testing.T) {
	label := &impls.Label{}
	err := label.AddLabel("test", []string{"3000", "5000", "8000"})
	if err == nil {
		t.Error("expected error for more than 2 ports, got nil")
	}
}

func TestCheckOverlap(t *testing.T) {
	viper.Reset()
	viper.Set("labels.existing.min", "3000")
	viper.Set("labels.existing.max", "5000")

	tests := []struct {
		name    string
		min     int
		max     int
		wantErr bool
	}{
		{"no overlap — below existing", 1000, 2999, false},
		{"no overlap — above existing", 5001, 6000, false},
		{"overlap — inside existing", 3500, 4500, true},
		{"overlap — straddles min", 2000, 3500, true},
		{"overlap — straddles max", 4500, 6000, true},
		{"overlap — contains existing", 2000, 6000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := impls.CheckOverlap("newlabel", tt.min, tt.max)
			if tt.wantErr && err == nil {
				t.Errorf("expected overlap error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
