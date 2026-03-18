package db

import "testing"

func TestStatusIDToName(t *testing.T) {
	// Arrange

	tests := []struct {
		name string
		id   int
		want PRStatus
	}{
		{name: "open", id: 1, want: PRStatusOpen},
		{name: "merged", id: 2, want: PRStatusMerged},
		{name: "default", id: 999, want: PRStatusOpen},
	}

	// Act

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.id
			want := tt.want

			got := StatusIDToName(id)

			// Assert
			if got != want {
				t.Fatalf("StatusIDToName(%d) = %q, want %q", id, got, want)
			}
		})
	}
}
