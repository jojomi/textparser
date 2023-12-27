package textparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_RemainingRuneCount(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "basic",
			fields: fields{
				input:    []rune("abc"),
				position: 0,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.RemainingRuneCount(), "RemainingRuneCount()")
		})
	}
}
