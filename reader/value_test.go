package reader

//
//import (
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestAtomicValueBool(t *testing.T) {
//	tvals := []interface{}{
//		true,
//		`true`,
//		`1`,
//		1,
//	}
//
//	for _, val := range tvals {
//		av := atomicValue{}
//		av.Store(val)
//
//		res, err := av.Bool()
//		assert.NoError(t, err)
//		assert.Equal(t, true, res)
//	}
//
//	fvals := []interface{}{
//		false,
//		`false`,
//		`0`,
//		0,
//	}
//
//	for _, val := range fvals {
//		av := atomicValue{}
//		av.Store(val)
//
//		res, err := av.Bool()
//		assert.NoError(t, err)
//		assert.Equal(t, false, res)
//	}
//}
//
//type testStruct struct {
//	Name string
//}
//
//func (ts testStruct) String() string {
//	return ts.Name
//}
//
//func TestAtomicValueString(t *testing.T) {
//	vals := []interface{}{
//		[]byte(`hello`),
//		testStruct{`hello`},
//	}
//
//	for _, val := range vals {
//		av := &atomicValue{}
//		av.Store(val)
//		res, err := av.String()
//
//		assert.NoError(t, err)
//		assert.Equal(t, `hello`, res)
//	}
//
//	intvals := []interface{}{
//		2,
//		int32(2),
//		int64(2),
//	}
//
//	for _, val := range intvals {
//		av := &atomicValue{}
//		av.Store(val)
//		res, err := av.String()
//
//		assert.NoError(t, err)
//		assert.Equal(t, `2`, res)
//	}
//
//	floatvals := []interface{}{
//		2.5,
//		float32(2.5),
//	}
//
//	for _, val := range floatvals {
//		av := &atomicValue{}
//		av.Store(val)
//		res, err := av.String()
//
//		assert.NoError(t, err)
//		assert.Equal(t, `2.5`, res)
//	}
//}
//
//func TestAtomicValueInt64(t *testing.T) {
//	intvals := []interface{}{
//		int32(100),
//		int64(100),
//		float32(100.01),
//		100.01, // float64
//		`100`,
//	}
//
//	for _, val := range intvals {
//		av := &atomicValue{}
//		av.Store(val)
//		res, err := av.Int64()
//
//		assert.NoError(t, err)
//		assert.Equal(t, int64(100), res)
//	}
//}
//
//func TestAtomicValueFloat64(t *testing.T) {
//	floatvals := []interface{}{
//		float32(5.0),
//		5,
//		int32(5),
//		int64(5),
//		`5.0`,
//	}
//
//	for _, val := range floatvals {
//		av := &atomicValue{}
//		av.Store(val)
//		res, err := av.Float64()
//
//		assert.NoError(t, err)
//		assert.Equal(t, 5.0, res)
//	}
//}
