package textparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_MustReadToMatchingString(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	type args struct {
		open  string
		close string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantLookingAt string
	}{
		{
			name: "basic",
			fields: fields{
				input:    []rune("((abc))"),
				position: 2,
			},
			args: args{
				"((",
				"))",
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.MustReadToMatchingString(tt.args.open, tt.args.close), "ReadToMatchingString(%v, %v)", tt.args.open, tt.args.close)
			assert.Truef(t, x.LookingAtString(tt.wantLookingAt), "wrong position, looking at %v", x.CurrentContext())
		})
	}
}

func TestParser_MustReadToMatchingStringSkipDelims(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	type args struct {
		open  string
		close string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantLookingAt string
	}{
		{
			name: "basic",
			fields: fields{
				input:    []rune("((abc))d"),
				position: 0,
			},
			args: args{
				"((",
				"))",
			},
			want:          "abc",
			wantLookingAt: "d",
		},
		{
			name: "quote (same open and close rune)",
			fields: fields{
				input:    []rune("he said 'yes', sir"),
				position: 8,
			},
			args: args{
				"'",
				"'",
			},
			want:          "yes",
			wantLookingAt: ", sir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.MustReadToMatchingStringSkipDelims(tt.args.open, tt.args.close), "ReadToMatchingStringSkipDelims(%v, %v)", tt.args.open, tt.args.close)
			assert.Truef(t, x.LookingAtString(tt.wantLookingAt), "wrong position, looking at %v", x.CurrentContext())
		})
	}
}

func TestParser_MustReadToMatchingRune(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	type args struct {
		open  rune
		close rune
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantLookingAt string
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
			want:          "ab",
			wantLookingAt: "]c[de]f)",
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
			want:          "a[b]c",
			wantLookingAt: "]def)",
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
			want:          "",
			wantLookingAt: "]def",
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
			assert.Equalf(t, tt.want, x.MustReadToMatchingRune(tt.args.open, tt.args.close), "ReadToMatchingRune(%v, %v)", tt.args.open, tt.args.close)
			assert.Truef(t, x.LookingAtString(tt.wantLookingAt), "wrong position, looking at %v", x.CurrentContext())
		})
	}
}

func TestParser_MustReadToMatchingRuneSkipDelims(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	type args struct {
		open  rune
		close rune
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantLookingAt string
	}{
		{
			name: "basic",
			fields: fields{
				input:    []rune("(abc)d"),
				position: 0,
			},
			args: args{
				'(',
				')',
			},
			want:          "abc",
			wantLookingAt: "d",
		},
		{
			name: "quote (same open and close rune)",
			fields: fields{
				input:    []rune("he said 'yes', sir"),
				position: 8,
			},
			args: args{
				'\'',
				'\'',
			},
			want:          "yes",
			wantLookingAt: ", sir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			assert.Equalf(t, tt.want, x.MustReadToMatchingRuneSkipDelims(tt.args.open, tt.args.close), "ReadToMatchingRuneSkipDelims(%v, %v)", tt.args.open, tt.args.close)
			assert.Truef(t, x.LookingAtString(tt.wantLookingAt), "wrong position, looking at %v", x.CurrentContext())
		})
	}
}

func TestParser_ReadToMatchingRuneEscaped(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	type args struct {
		open   rune
		close  rune
		escape rune
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantLookingAt string
	}{
		{
			name: "nothing",
			fields: fields{
				input:    []rune(`my "" is`),
				position: 4,
			},
			args: args{
				'"',
				'"',
				'\\',
			},
			want:          ``,
			wantLookingAt: `" is`,
		},
		{
			name: "escape only",
			fields: fields{
				input:    []rune(`my (\)) is`),
				position: 4,
			},
			args: args{
				'(',
				')',
				'\\',
			},
			want:          `)`,
			wantLookingAt: `) is`,
		},
		{
			name: "basic",
			fields: fields{
				input:    []rune(`my "best \" attempt" is`),
				position: 4,
			},
			args: args{
				'"',
				'"',
				'\\',
			},
			want:          `best " attempt`,
			wantLookingAt: `" is`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			actual, err := x.ReadToMatchingRuneEscaped(tt.args.open, tt.args.close, tt.args.escape)
			assert.Nil(t, err)
			assert.Equalf(t, tt.want, actual, "ReadToMatchingRuneEscaped(%v, %v, %v)", tt.args.open, tt.args.close, tt.args.escape)
			assert.Truef(t, x.LookingAtString(tt.wantLookingAt), "wrong position, looking at %v", x.CurrentContext())
		})
	}
}

func TestParser_ReadToMatchingRuneEscapedSkipDelims(t *testing.T) {
	type fields struct {
		input    []rune
		position int
	}
	type args struct {
		open   rune
		close  rune
		escape rune
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantLookingAt string
	}{
		{
			name: "nothing",
			fields: fields{
				input:    []rune(`my "" is`),
				position: 3,
			},
			args: args{
				'"',
				'"',
				'\\',
			},
			want:          ``,
			wantLookingAt: ` is`,
		},
		{
			name: "escape only",
			fields: fields{
				input:    []rune(`my (\)) is`),
				position: 3,
			},
			args: args{
				'(',
				')',
				'\\',
			},
			want:          `)`,
			wantLookingAt: ` is`,
		},
		{
			name: "basic",
			fields: fields{
				input:    []rune(`my "best \" attempt" is`),
				position: 3,
			},
			args: args{
				'"',
				'"',
				'\\',
			},
			want:          `best " attempt`,
			wantLookingAt: ` is`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Parser{
				input:    tt.fields.input,
				position: tt.fields.position,
			}
			actual, err := x.ReadToMatchingRuneEscapedSkipDelims(tt.args.open, tt.args.close, tt.args.escape)
			assert.Nil(t, err)
			assert.Equalf(t, tt.want, actual, "ReadToMatchingRuneEscapedSkipDelims(%v, %v, %v)", tt.args.open, tt.args.close, tt.args.escape)
			assert.Truef(t, x.LookingAtString(tt.wantLookingAt), "wrong position, looking at %v", x.CurrentContext())
		})
	}
}
