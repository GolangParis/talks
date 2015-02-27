// STARTIMPORT OMIT

package rot13 // HLxxx

import (
	"reflect"
	"testing" // HLxxx
)

// ENDIMPORT OMIT

func TestRot13(t *testing.T) {
	for _, table := range []struct {
		str  []byte
		want []byte
	}{
		{str: []byte("Lbh penpxrq gur pbqr!"), want: []byte("You cracked the code!")},
		{str: []byte("hello"), want: []byte("uryyb")},
	} {
		o := make([]byte, len(table.str))
		for i, b := range table.str {
			o[i] = rot13(b) // HLxxx
		}

		if !reflect.DeepEqual(o, table.want) {
			t.Fatalf("invalid rot13\ngot =%q\nwant=%q\n", string(o), string(table.want))
		}
	}
}
