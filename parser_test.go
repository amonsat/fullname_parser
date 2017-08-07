package fullname_parser

import (
	"reflect"
	"testing"
)

func TestParseFullname(t *testing.T) {
	tests := []struct {
		name           string
		fullname       string
		wantParsedName ParsedName
	}{
		{"base test", "Juan Xavier", ParsedName{First: "Juan", Last: "Xavier"}},
		{"title test", "Dr. Juan Xavier", ParsedName{Title: "Dr.", First: "Juan", Last: "Xavier"}},
		{"nick test", "Dr. Juan Xavier (Doc Vega)", ParsedName{Title: "Dr.", First: "Juan", Last: "Xavier", Nick: "Doc Vega"}},
		{"middle test", "Juan Q. Xavier", ParsedName{First: "Juan", Middle: "Q.", Last: "Xavier"}},
		{"suffixes test", "Juan Xavier III (Doc Vega), Jr.", ParsedName{First: "Juan", Last: "Xavier", Nick: "Doc Vega", Suffix: "III, Jr."}},
		{"full test", "de la Vega, Dr. Juan et Glova (Doc Vega) Q. Xavier III, Jr., Genius", ParsedName{Title: "Dr.", First: "Juan et Glova", Middle: "Q. Xavier", Last: "de la Vega", Nick: "Doc Vega", Suffix: "III, Jr., Genius"}},
		{"just last name", "Cotter", ParsedName{Last: "Cotter"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotParsedName := ParseFullname(tt.fullname); !reflect.DeepEqual(gotParsedName, tt.wantParsedName) {
				t.Errorf("ParseFullname() = %v, want %v", gotParsedName, tt.wantParsedName)
			}
		})
	}
}
