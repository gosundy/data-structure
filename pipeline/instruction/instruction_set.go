package instruction

import "io"

type Source struct {
}

func NewSource() *Source {
	return &Source{}
}
func (source *Source) Process(Function func(readers []io.Reader, wirters []io.Writer) error) error {
	return nil
}
