package packfile

import (
	"bytes"
	"fmt"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/format/packfile"
	"github.com/go-git/go-git/v6/storage/memory"
)

// GitObject represents a git object from a packfile
type GitObject struct {
	Type    plumbing.ObjectType
	Content []byte
	Hash    plumbing.Hash
}

// PackfileWriter writes git objects to a packfile
type PackfileWriter struct{}

// NewPackfileWriter creates a new packfile writer
func NewPackfileWriter() *PackfileWriter {
	return &PackfileWriter{}
}

// WritePackfile writes a set of git objects to a packfile using go-git's encoder.
func (w *PackfileWriter) WritePackfile(objects []*GitObject) ([]byte, error) {
	ms := memory.NewStorage()

	var hashes []plumbing.Hash
	for _, obj := range objects {
		enc := ms.NewEncodedObject()
		enc.SetType(obj.Type)
		enc.SetSize(int64(len(obj.Content)))

		writer, err := enc.Writer()
		if err != nil {
			return nil, fmt.Errorf("getting writer: %w", err)
		}

		_, err = writer.Write(obj.Content)
		if err != nil {
			return nil, fmt.Errorf("writing content: %w", err)
		}

		hash, err := ms.SetEncodedObject(enc)
		if err != nil {
			return nil, fmt.Errorf("storing object: %w", err)
		}

		hashes = append(hashes, hash)
	}

	var buf bytes.Buffer
	encoder := packfile.NewEncoder(&buf, ms, false)

	_, err := encoder.Encode(hashes, 10)
	if err != nil {
		return nil, fmt.Errorf("encoding packfile: %w", err)
	}

	return buf.Bytes(), nil
}
