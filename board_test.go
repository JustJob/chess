package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMoveFromStringPos(t *testing.T) {
	type args struct {
		pos string
	}
	tests := []struct {
		name    string
		args    args
		want    *Move
		wantErr bool
	}{
		{"happy", args{"a8 to b6"}, &Move{startIndex: 56, endIndex: 41}, false},
		{"happyAltFormat", args{"a8 b6"}, &Move{startIndex: 56, endIndex: 41}, false},
		{"same", args{"a8 to a8"}, &Move{startIndex: 56, endIndex: 56}, false},
		{"invalidRow", args{"a9 to b6"}, nil, true},
		{"invalidCol", args{"i4 to b6"}, nil, true},
		{"invalidCol2nd", args{"a3 to b10"}, nil, true},
		{"invalidRow2nd", args{"a3 to n6"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MoveFromStringPos(tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveFromStringPos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MoveFromStringPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyMove(t *testing.T) {
	b := NewStandardBoard()
	m, err := MoveFromStringPos("b2 b4")
	require.NoError(t, err)

	expectedStart := convertPosToBitmap(9)
	expectedEnd := convertPosToBitmap(25)

	require.NotEqual(t, uint64(0), b.pieces[WhiteIndex][PawnIndex]&expectedStart)
	require.Equal(t, uint64(0), b.pieces[WhiteIndex][PawnIndex]&expectedEnd)

	err = b.applyMove(m)
	require.NoError(t, err)

	require.Equal(t, uint64(0), b.pieces[WhiteIndex][PawnIndex]&expectedStart)
	require.NotEqual(t, uint64(0), b.pieces[WhiteIndex][PawnIndex]&expectedEnd)
}

func Test_addPiece(t *testing.T) {
	type args struct {
		bitmap   uint64
		charCode byte
		pieces   []byte
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"happy", args{1 << 5, 'a', make([]byte, 64)}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addPiece(tt.args.bitmap, tt.args.charCode, tt.args.pieces)
			for i := uint8(0); i < 64; i++ {
				if i == tt.want {
					if !reflect.DeepEqual(tt.args.pieces[i], tt.args.charCode) {
						t.Errorf("addPiece(): index %d wrong value %b", i, tt.args.pieces[i])
					}
				} else {
					if !reflect.DeepEqual(tt.args.pieces[i], byte(0)) {
						t.Errorf("addPiece(): index %d wrong value %b", i, tt.args.pieces[i])
					}
				}
			}
		})
	}
}
