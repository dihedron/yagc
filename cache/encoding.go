package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type Encoding[K comparable, V any] interface {
	Encode(data map[K]V) ([]byte, error)
	Decode(data []byte) (map[K]V, error)
}

type JSON[K comparable, V any] struct {
	Pretty bool
}

func (j *JSON[K, V]) Encode(data map[K]V) ([]byte, error) {
	if j.Pretty {
		return json.MarshalIndent(data, "", "  ")
	} else {
		return json.Marshal(data)
	}
}

func (j *JSON[K, V]) Decode(data []byte) (map[K]V, error) {
	m := map[K]V{}
	err := json.Unmarshal(data, &m)
	return m, err
}

type GOB[K comparable, V any] struct {
	buffer  bytes.Buffer
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func (g *GOB[K, V]) Encode(data map[K]V) ([]byte, error) {
	if g.encoder == nil {
		g.encoder = gob.NewEncoder(&g.buffer)
	}
	g.buffer.Reset()
	if err := g.encoder.Encode(data); err != nil {
		return nil, err
	}
	return g.buffer.Bytes(), nil
}

func (g *GOB[K, V]) Decode(data []byte) (map[K]V, error) {
	if g.encoder == nil {
		g.decoder =
			gob.NewDecoder(&g.buffer)
	}
	m := map[K]V{}
	g.buffer.ReadFrom(bytes.NewReader(data))
	if err := g.decoder.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
