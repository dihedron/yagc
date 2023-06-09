package cache

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Persistence defines the behaviour of how a Cache contents
// get persisted.
type Persistence interface {
	Write(data []byte) error
	Read() ([]byte, error)
}

// File persists the encoded data, and reads it back from a
// given file.
type File struct {
	Path string
}

// Write writes data to the given file.
func (f *File) Write(data []byte) error {
	return os.WriteFile(f.Path, data, 0644)
}

// Read reads data back from the given file.
func (f *File) Read() ([]byte, error) {
	return os.ReadFile(f.Path)
}

// Console persists the encoded data to the console; it cannot read
// it back though...
type Console struct {
	Writer io.Writer
}

// Write writes data to the given file.
func (c *Console) Write(data []byte) error {
	_, err := fmt.Fprintf(c.Writer, "%s\n", string(data))
	return err
}

// Read reads data back from the given file.
func (*Console) Read() ([]byte, error) {
	return nil, errors.New("not implemented")
}

// Discard does not persist data anywhere, nor can it recover it.
type Discard struct{}

// Write discards data written to it.
func (*Discard) Write(_ []byte) error {
	return nil
}

// Read always returns no data and a generic "not implemented" error.
func (*Discard) Read() ([]byte, error) {
	return nil, errors.New("not implemented")
}
