package reader

var (
	_ Value = (*NullValue)(nil)
)

type Value interface {
	Int64() (int64, error)
	String() (string, error)
	Float64() (float64, error)
	Bool() (bool, error)
	Bytes() ([]byte, error)
	Scan(v interface{}) error
}

type NullValue struct {
}

func (n *NullValue) Int64() (int64, error) {
	return 0, nil
}

func (n NullValue) String() (string, error) {
	return ``, nil
}

func (n NullValue) Float64() (float64, error) {
	return 0.0, nil
}

func (n NullValue) Bool() (bool, error) {
	return false, nil
}

func (n NullValue) Bytes() ([]byte, error) {
	return []byte{}, nil
}

func (n NullValue) Scan(v interface{}) error {
	return nil
}
