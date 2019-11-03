package csrf

import (
	"github.com/google/uuid"
	"testing"
)

var testSecret = []byte("test_secret")
var testUUID = uuid.New()
var testJWT = JwtToken{Secret:testSecret}
var testCSRF, _ = testJWT.Create(testUUID)

func TestJwtToken_Check(t *testing.T) {
	type fields struct {
		Secret []byte
	}
	type args struct {
		s          uuid.UUID
		inputToken string
	}
	var tests = []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "positive",
			fields: fields{
				Secret: []byte("test_secret"),
			},
			args: args{
				s:          testUUID,
				inputToken: testCSRF,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "negative",
			fields: fields{
				Secret: []byte("test_secret"),
			},
			args: args{
				s:          testUUID,
				inputToken: "some_invalid_CSRF",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tk := &JwtToken{
				Secret: tt.fields.Secret,
			}
			got, err := tk.Check(tt.args.s, tt.args.inputToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Check() got = %v, want %v", got, tt.want)
			}
		})
	}
}
