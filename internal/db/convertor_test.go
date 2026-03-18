package db

import "testing"

func TestStatusIDToName(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want PRStatus
	}{
		{name: "open", id: 1, want: PRStatusOpen},
		{name: "merged", id: 2, want: PRStatusMerged},
		{name: "default", id: 999, want: PRStatusOpen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusIDToName(tt.id); got != tt.want {
				t.Fatalf("StatusIDToName(%d) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}
