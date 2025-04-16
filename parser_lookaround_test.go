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

func TestParser_LookingAtWhitespace(t *testing.T) {
	basicInput := []rune(`a b	c
d`)

	type fields struct {
		input    []rune
		position int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "basic",
			fields: fields{
				input:    basicInput,
				position: 0,
			},
			want: false,
		},
		{
			name: "space",
			fields: fields{
				input:    basicInput,
				position: 1,
			},
			want: true,
		},
		{
			name: "tab",
			fields: fields{
				input:    basicInput,
				position: 3,
			},
			want: true,
		},
		{
			name: "newline",
			fields: fields{
				input:    basicInput,
				position: 5,
			},
			want: true,
		},
		{
			name: "end of input",
			fields: fields{
				input:    basicInput,
				position: 6,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.LookingAtWhitespace(), "LookingAtWhitespace()")
		})
	}
}

func TestParser_LookingAtDigit(t *testing.T) {
	basicInput := []rune(`a 4ö3`)

	type fields struct {
		input    []rune
		position int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "basic",
			fields: fields{
				input:    basicInput,
				position: 0,
			},
			want: false,
		},
		{
			name: "single",
			fields: fields{
				input:    basicInput,
				position: 2,
			},
			want: true,
		},
		{
			name: "in word",
			fields: fields{
				input:    basicInput,
				position: 4,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.LookingAtDigit(), "LookingAtDigit()")
		})
	}
}

func TestParser_LookingAtLetter(t *testing.T) {
	basicInput := []rune(`a 4b3`)

	type fields struct {
		input    []rune
		position int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "basic",
			fields: fields{
				input:    basicInput,
				position: 0,
			},
			want: true,
		},
		{
			name: "single",
			fields: fields{
				input:    basicInput,
				position: 1,
			},
			want: false,
		},
		{
			name: "in word",
			fields: fields{
				input:    basicInput,
				position: 3,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.LookingAtLetter(), "LookingAtLetter()")
		})
	}
}
