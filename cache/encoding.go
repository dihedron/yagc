package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Encoding[K comparable, V any] interface {
	Encode(data map[K]V) ([]byte, error)
	Decode(data []byte) (map[K]V, error)
}

// JSON encodes/decodes cache data in JSON format.
type JSON[K comparable, V any] struct {
	Pretty bool
}

// Encode encodes cache data in JSON format.
func (j *JSON[K, V]) Encode(data map[K]V) ([]byte, error) {
	if j.Pretty {
		return json.MarshalIndent(data, "", "  ")
	} else {
		return json.Marshal(data)
	}
}

// Decode decodes cache data from JSON format.
func (j *JSON[K, V]) Decode(data []byte) (map[K]V, error) {
	m := map[K]V{}
	err := json.Unmarshal(data, &m)
	return m, err
}

// YAML encodes/decodes cache data in YAML format.
type YAML[K comparable, V any] struct{}

// Encode encodes cache data in YAML format.
func (y *YAML[K, V]) Encode(data map[K]V) ([]byte, error) {
	return yaml.Marshal(data)
}

// Decode decodes cache data from YAML format.
func (y *YAML[K, V]) Decode(data []byte) (map[K]V, error) {
	m := map[K]V{}
	err := yaml.Unmarshal(data, &m)
	return m, err
}

// GOB encodes/decodes cache data in self-describing binary format.
type GOB[K comparable, V any] struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
}

type Wrapper[K comparable, V any] struct {
	Map map[K]V
}

// Encode encodes cache data in self-describing binary format.
func (g *GOB[K, V]) Encode(data map[K]V) ([]byte, error) {
	var buffer bytes.Buffer
	gob.RegisterName("Wrapper", Wrapper[K, V]{})
	g.encoder = gob.NewEncoder(&buffer)
	m := Wrapper[K, V]{
		Map: data,
	}
	if err := g.encoder.Encode(&m); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Decode decodes cache data from self-describing binary format.
func (g *GOB[K, V]) Decode(data []byte) (map[K]V, error) {
	gob.RegisterName("Wrapper", Wrapper[K, V]{})
	g.decoder = gob.NewDecoder(bytes.NewReader(data))
	m := Wrapper[K, V]{}
	if err := g.decoder.Decode(&m); err != nil {
		return nil, err
	}
	return m.Map, nil
}
