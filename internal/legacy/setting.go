package legacy

type Setting struct {
	ID    int    `storm:"id,increment"`
	Name  string `storm:"unique"`
	Value interface{}
}
