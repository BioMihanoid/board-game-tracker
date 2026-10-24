package console

import "testing"

func TestParseKopecks(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{name: "целое", input: "5999", want: 5999_00},
		{name: "запятая", input: "5999,99", want: 5999_99},
		{name: "точка", input: "5999.99", want: 5999_99},
		{name: "пробелы в сумме", input: "5 999,50", want: 5999_50},
		{name: "пробелы по краям", input: "  1500  ", want: 1500_00},
		{name: "ноль", input: "0", want: 0},
		{name: "пусто", input: "", wantErr: true},
		{name: "мусор", input: "abc", wantErr: true},
		{name: "отрицательная", input: "-1", wantErr: true},
		{name: "три знака", input: "5999,999", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseKopecks(test.input)

			if test.wantErr {
				if err == nil {
					t.Fatalf("ожидалась ошибка, получено %d", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if got != test.want {
				t.Fatalf("получено %d, ожидалось %d", got, test.want)
			}
		})
	}
}
