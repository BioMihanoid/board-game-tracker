package model

import "testing"

func TestFormatKopecks(t *testing.T) {
	tests := []struct {
		kopecks int64
		want    string
	}{
		{kopecks: 0, want: "0,00"},
		{kopecks: 5, want: "0,05"},
		{kopecks: 5999_00, want: "5 999,00"},
		{kopecks: 22999_94, want: "22 999,94"},
		{kopecks: 1_234_567_89, want: "1 234 567,89"},
		{kopecks: -1500_50, want: "-1 500,50"},
	}

	for _, test := range tests {
		if got := FormatKopecks(test.kopecks); got != test.want {
			t.Errorf("FormatKopecks(%d) = %q, ожидалось %q", test.kopecks, got, test.want)
		}
	}
}
