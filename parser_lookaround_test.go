package textparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_LookingAtRune(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "basic",
			fields: fields{
				input:    []rune("a b"),
				position: 0,
			},
			args: args{'a'},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.LookingAtRune(tt.args.r), "LookingAtRune(%v)", tt.args.r)
		})
	}
}
