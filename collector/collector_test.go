package collector

import "testing"

func Test_parseDbValue(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want float64
		err  error
	}{
		{
			name: "test positive value in dB",
			str:  "55dB",
			want: 55,
		},
		{
			name: "test negative value in dB",
			str:  "-55dB",
			want: -55,
		},
		{
			name: "test positive value in dBm",
			str:  "55dBm",
			want: 55,
		},
		{
			name: "test negative value in dBm",
			str:  "-55dBm",
			want: -55,
		},
		{
			name: "test greater than value in dB",
			str:  "&gt;=30dB",
			want: 30,
		},
		{
			name: "test greater than value in dBm",
			str:  "&gt;=30dBm",
			want: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDbValue(tt.str)
			if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Fatalf("parseDbValue() error = %s, wantErr %s", err, tt.err)
				return
			}
			if tt.err != nil && err != nil {
				t.Fatal("Expected err got nil")
				return
			}
			if got != tt.want {
				t.Errorf("parseDbValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
