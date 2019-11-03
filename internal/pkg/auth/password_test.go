package auth

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	type args struct {
		userPassword   string
		passwordHashed string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive",
			args: args{
				userPassword:   "ABC",
				passwordHashed: string(HashPassword("ABC")),
			},
			want: true,
		},
		{
			name: "negative",
			args: args{
				userPassword:   "ABC",
				passwordHashed: "invalid_saltABC",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPassword(tt.args.userPassword, tt.args.passwordHashed); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
