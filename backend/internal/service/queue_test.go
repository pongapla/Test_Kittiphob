package service

import "testing"

func TestFormatQueueNumber(t *testing.T) {
	tests := []struct {
		name  string
		index int
		want  string
	}{
		{"FirstQueue", 0, "A0"},
		{"SecondQueue", 1, "A1"},
		{"EndOfA", 9, "A9"},
		{"StartOfB", 10, "B0"},
		{"EndOfB", 19, "B9"},
		{"StartOfC", 20, "C0"},
		{"LastQueue", 259, "Z9"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatQueueNumber(tt.index)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestFormatQueueNumberInvalidIndex(t *testing.T) {
	for _, index := range []int{-1, 260, 999} {
		_, err := FormatQueueNumber(index)

		if err == nil {
			t.Errorf("expected error for index %d", index)
		}
	}
}

func TestAllQueueNumbersAreUnique(t *testing.T) {
	seen := make(map[string]bool)

	for i := 0; i < 260; i++ {
		number, err := FormatQueueNumber(i)
		if err != nil {
			t.Fatal(err)
		}

		if seen[number] {
			t.Fatalf("duplicate queue number: %s", number)
		}

		seen[number] = true
	}

	if len(seen) != 260 {
		t.Errorf("expected 260 unique queues, got %d", len(seen))
	}
}
