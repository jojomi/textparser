package textparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_MustReadToMatching(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	type args struct {
		open  rune
		close rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				input:    []rune("(abc)"),
				position: 1,
			},
			args: args{
				'(',
				')',
			},
			want: "abc",
		},
		{
			name: "basic rune",
			fields: fields{
				input:    []rune("ß{ÄÖü}?-#Ü"),
				position: 2,
			},
			args: args{
				'-',
				'}',
			},
			want: "ÄÖü",
		},
		{
			name: "multiline",
			fields: fields{
				input:    []rune("-\n\n-"),
				position: 1,
			},
			args: args{
				'-',
				'-',
			},
			want: "\n\n",
		},
		{
			name: "empty",
			fields: fields{
				input:    []rune("{}"),
				position: 1,
			},
			args: args{
				'{',
				'}',
			},
			want: "",
		},
		{
			name: "repeated",
			fields: fields{
				input:    []rune("[ab]c[de]f)"),
				position: 1,
			},
			args: args{
				'[',
				']',
			},
			want: "ab",
		},
		{
			name: "nested",
			fields: fields{
				input:    []rune("([a[b]c]def)"),
				position: 2,
			},
			args: args{
				'[',
				']',
			},
			want: "a[b]c",
		},
		{
			name: "not inside",
			fields: fields{
				input:    []rune("]def"),
				position: 0,
			},
			args: args{
				'[',
				']',
			},
			want: "",
		},
		{
			name: "quote (same open and close rune)",
			fields: fields{
				input:    []rune("he said 'yes'"),
				position: 9,
			},
			args: args{
				'\'',
				'\'',
			},
			want: "yes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.MustReadToMatching(tt.args.open, tt.args.close), "ReadToMatching(%v, %v)", tt.args.open, tt.args.close)
		})
	}
}
