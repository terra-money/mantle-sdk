package kvindex

import (
	"fmt"
	"github.com/terra-project/mantle/utils"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndexMap(t *testing.T) {
	var indexMapSlice []IndexMapEntry
	var testStruct interface{}

	// passing
	func() {
		type TestIndexStruct struct {
			Foo string
			Bar struct {
				Hello string `mantle:"index"`
			}
		}

		indexMapSlice = []IndexMapEntry{}
		testStruct = (*TestIndexStruct)(nil)
		assert.NotPanics(
			t,
			func() {
				indexMapSlice = createIndexMap(reflect.TypeOf(testStruct))
			},
		)

		assert.Equal(t, len(indexMapSlice), 1)
		assert.Equal(t, indexMapSlice[0].Type, reflect.String)
		assert.Equal(t, indexMapSlice[0].Name, "Hello")
		assert.Equal(t, indexMapSlice[0].Path, []string{"Bar", "Hello"})
	}()

	// panic
	func() {
		type TestIndexStructFailing struct {
			Foo struct {
				Hello string
			} `mantle:"index"` // should panic
			Bar struct {
				Hello string `mantle:"index"`
			}
		}
		indexMapSlice = []IndexMapEntry{}
		testStruct = (*TestIndexStructFailing)(nil)

		assert.Panics(
			t,
			func() {
				indexMapSlice = createIndexMap(reflect.TypeOf(testStruct))
			},
		)
	}()

	// panic due to disallowed character
	func() {
		type TestIndexStructFailingDueToDisallowedCharacter struct {
			Foo string
			Bar struct {
				Hello string `mantle:"index=$money"`
			}
		}

		indexMapSlice = []IndexMapEntry{}
		testStruct = (*TestIndexStructFailingDueToDisallowedCharacter)(nil)

		assert.Panics(
			t,
			func() {
				indexMapSlice = createIndexMap(reflect.TypeOf(testStruct))
			},
		)
	}()

	// complex
	func() {
		type TestIndexStructComplex struct {
			Foo struct {
				Bar struct {
					Hello struct {
						World int `mantle:"index"`
					}
				}
				Bar2 struct {
					Whatever struct {
						Leaf string `mantle:"index"`
					}
				}
			}
			Foo2 []struct {
				Hello string
				World struct {
					Leaf string `mantle:"index=foo2leaf"`
				}
			}
		}

		indexMapSlice = []IndexMapEntry{}
		testStruct = (*TestIndexStructComplex)(nil)

		assert.NotPanics(
			t,
			func() {
				indexMapSlice = createIndexMap(reflect.TypeOf(testStruct))
			},
		)

		assert.Equal(t, len(indexMapSlice), 3)
		assert.Equal(t, indexMapSlice[0].Type, reflect.Int)
		assert.Equal(t, indexMapSlice[0].Name, "World")
		assert.Equal(t, indexMapSlice[0].Path, []string{"Foo", "Bar", "Hello", "World"})
		assert.Equal(t, indexMapSlice[1].Type, reflect.String)
		assert.Equal(t, indexMapSlice[1].Name, "Leaf")
		assert.Equal(t, indexMapSlice[1].Path, []string{"Foo", "Bar2", "Whatever", "Leaf"})
		assert.Equal(t, indexMapSlice[2].Type, reflect.String)
		assert.Equal(t, indexMapSlice[2].Name, "foo2leaf")
		assert.Equal(t, indexMapSlice[2].Path, []string{"Foo2", "*", "World", "Leaf"})
	}()

	// slice entity
	func() {
		type TestSliceIndex []struct {
			Foo string `mantle:"index=foo"`
			Bar uint64 `mantle:"index=bar"`
		}

		indexMapSlice = []IndexMapEntry{}
		testStruct = (*TestSliceIndex)(nil)

		assert.NotPanics(
			t,
			func() {
				indexMapSlice = createIndexMap(reflect.TypeOf(testStruct))
			},
		)

		assert.Equal(t, reflect.String, indexMapSlice[0].Type)
		assert.Equal(t, "foo", indexMapSlice[0].Name)
		assert.Equal(t, []string{"*", "Foo"}, indexMapSlice[0].Path)
		assert.Equal(t, reflect.Uint64, indexMapSlice[1].Type)
		assert.Equal(t, "bar", indexMapSlice[1].Name)
		assert.Equal(t, []string{"*", "Bar"}, indexMapSlice[1].Path)
	}()

}

func TestNewKVIndex(t *testing.T) {
	var prefix []byte

	// duplicate index is not allowed
	func() {
		type TestIndexStructDuplicate struct {
			Foo struct {
				Hello string `mantle:"index"`
			}
			Bar struct {
				Hello string `mantle:"index"`
			}
		}
		testStruct := (*TestIndexStructDuplicate)(nil)
		assert.Panics(t, func() { NewKVIndex(reflect.TypeOf(testStruct)) })
	}()

	// get prefix
	func() {
		type TestIndexStructComplex struct {
			Foo struct {
				Bar struct {
					Hello struct {
						World int `mantle:"index"`
					}
				}
				Bar2 struct {
					Whatever struct {
						Leaf string `mantle:"index"`
					}
				}
			}
			Foo2 []struct {
				Hello string
				World struct {
					Leaf string `mantle:"index=foo2leaf"`
				}
			}
		}

		testStruct := (*TestIndexStructComplex)(nil)
		kvindex := NewKVIndex(reflect.TypeOf(testStruct))
		prefix = kvindex.GetPrefix("World")
		assert.Equal(t, prefix, []byte("TestIndexStructComplexWorld"))

		// get prefix fails because this specific index does not exist
		assert.Panics(t, func() { kvindex.GetPrefix("Fail") })

		// get cursor success
		cursor, err := kvindex.GetCursor("foo2leaf", "hello")
		assert.Nil(t, err)
		assert.Equal(t, cursor, []byte("TestIndexStructComplexfoo2leafhello"))

		// get cursor panics because index is not defined
		assert.Panics(t, func() { kvindex.GetCursor("foo2leafX", "hello") })

		// get cursor fails because cursor is nil
		assert.Panics(t, func() { kvindex.GetCursor("foo2leaf", nil) })

		// get cursor fails because index type is different
		assert.Panics(t, func() { kvindex.GetCursor("foo2leaf", 1) })

		// getcursor fails cursor is ptr type
		assert.Panics(t, func() { kvindex.GetCursor("foo2leaf", &[]byte{1}) })
	}()

	// test uint64
	func() {
		type TestIndexStructUint64 struct {
			Foo struct {
				Hello uint64 `mantle:"index=id"`
				World int    `mantle:"index=id2"`
			}
		}

		testStruct := (*TestIndexStructUint64)(nil)
		kvindex := NewKVIndex(reflect.TypeOf(testStruct))

		assert.Equal(t, kvindex.GetPrefix("id"), []byte("TestIndexStructUint64id"))
		assert.Equal(t, kvindex.GetPrefix("id2"), []byte("TestIndexStructUint64id2"))

		// uint64 cursor (note for numbers cursor is in big endian)
		cursor, err := kvindex.GetCursor("id", uint64(172))
		assert.Nil(t, err)
		assert.Equal(t, cursor, append([]byte("TestIndexStructUint64id"), []byte{0, 0, 0, 0, 0, 0, 0, 172}...))

		// TODO: fix me
		// int cursor (note for signed ints, we take bit flipped version of it so lex compare works)
		cursor, err = kvindex.GetCursor("id2", int(-333))
		assert.Nil(t, err)
		fmt.Println(cursor)
	}()

	// test slice entity
	func() {
		type TestSliceEntity struct {
			Foo string `mantle:"index=foo"`
			Bar uint64 `mantle:"index=bar"`
		}

		testStruct := (*TestSliceEntity)(nil)
		kvindex := NewKVIndex(reflect.TypeOf(testStruct))

		assert.Equal(t, []byte("TestSliceEntityfoo"), kvindex.GetPrefix("foo"))
		assert.Equal(t, []byte("TestSliceEntitybar"), kvindex.GetPrefix("bar"))

		cursor, err := kvindex.GetCursor("foo", "hello")
		assert.Nil(t, err)
		assert.Equal(t, cursor, utils.ConcatBytes([]byte("TestSliceEntityfoo"), []byte("hello")))

	}()

}
