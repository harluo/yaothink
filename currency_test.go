package yaothink

import "testing"

func TestParseCurrency(t *testing.T) {
	for _, test := range []struct {
		code  string
		want  Currency
		valid bool
	}{
		{code: "USD", want: USD, valid: true},
		{code: "cny", want: CNY, valid: true},
		{code: "ZZZ", valid: false},
	} {
		got, err := ParseCurrency(test.code)
		if test.valid && (err != nil || got != test.want) {
			t.Errorf("ParseCurrency(%q) = %v, %v；期望 %v", test.code, got, err, test.want)
		} else if !test.valid && err == nil {
			t.Errorf("ParseCurrency(%q) 应返回错误", test.code)
		}
	}
}
