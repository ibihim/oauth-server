// generated by gotemplate

package opt

import (
	"fmt"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// template type Optional(A)

// A 'gotemplate'-based type for providing optional semantics without using pointers.
type Uint struct {
	V       uint
	Defined bool
}

// Creates an optional type with a given value.
func OUint(v uint) Uint {
	return Uint{V: v, Defined: true}
}

// Get returns the value or given default in the case the value is undefined.
func (v Uint) Get(deflt uint) uint {
	if !v.Defined {
		return deflt
	}
	return v.V
}

// MarshalEasyJSON does JSON marshaling using easyjson interface.
func (v Uint) MarshalEasyJSON(w *jwriter.Writer) {
	if v.Defined {
		w.Uint(v.V)
	} else {
		w.RawString("null")
	}
}

// UnmarshalEasyJSON does JSON unmarshaling using easyjson interface.
func (v *Uint) UnmarshalEasyJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		l.Skip()
		*v = Uint{}
	} else {
		v.V = l.Uint()
		v.Defined = true
	}
}

// MarshalJSON implements a standard json marshaler interface.
func (v *Uint) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	v.MarshalEasyJSON(&w)
	return w.Buffer.BuildBytes(), w.Error
}

// UnmarshalJSON implements a standard json unmarshaler interface.
func (v *Uint) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{}
	v.UnmarshalEasyJSON(&l)
	return l.Error()
}

// IsDefined returns whether the value is defined, a function is required so that it can
// be used in an interface.
func (v Uint) IsDefined() bool {
	return v.Defined
}

// String implements a stringer interface using fmt.Sprint for the value.
func (v Uint) String() string {
	if !v.Defined {
		return "<undefined>"
	}
	return fmt.Sprint(v.V)
}
