package brainit

import (
	"reflect"
	"testing"
)

func TestElement_Next(t *testing.T) {
	testMemory := NewMemory(nil)
	type args struct {
		current *Element
		v       interface{}
	}
	tests := []struct {
		name string
		args args
		want *Element
	}{
		{
			name: "next with value number",
			args: args{
				current: testMemory.Current,
				v:       10,
			},
			want: &Element{
				nextElement: nil,
				prevElement: testMemory.Front,
				list:        testMemory,
				Value:       10,
			},
		},
		{
			name: "next with exist",
			args: args{
				current: testMemory.Front,
				v:       1000,
			},
			want: &Element{
				nextElement: nil,
				prevElement: testMemory.Front,
				list:        testMemory,
				Value:       10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.current.Next(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_Prev(t *testing.T) {
	testMemory := NewMemory(nil)
	type args struct {
		current *Element
		v       interface{}
	}
	tests := []struct {
		name string
		args args
		want *Element
	}{
		{
			name: "prev with value number",
			args: args{
				current: testMemory.Current,
				v:       10,
			},
			want: &Element{
				nextElement: testMemory.Back,
				prevElement: nil,
				list:        testMemory,
				Value:       10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.current.Prev(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_RemoveUntil(t *testing.T) {
	testMemory := NewMemory(0)
	testMemory.Current = testMemory.Current.Next(1).Next(2).Next(3).Next(4).Prev(nil)

	if testMemory.Len != 5 {
		t.Errorf("Len problem")
	}

	testMemory.RemoveUntil(testMemory.Current)

	if testMemory.Len != 2 {
		t.Errorf("Len problem after erase")
	}

	want := &Element{
		nextElement: testMemory.Current.nextElement,
		prevElement: nil,
		list:        testMemory,
		Value:       3,
	}

	if !reflect.DeepEqual(testMemory.Current, want) {
		t.Errorf("Element.Prev() = %v, want %v", testMemory.Current, want)
	}
}
