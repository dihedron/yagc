package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/BurntSushi/toml"
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
func (*JSON[K, V]) Decode(data []byte) (map[K]V, error) {
	m := map[K]V{}
	err := json.Unmarshal(data, &m)
	return m, err
}

// YAML encodes/decodes cache data in YAML format.
type YAML[K comparable, V any] struct{}

// Encode encodes cache data in YAML format.
func (*YAML[K, V]) Encode(data map[K]V) ([]byte, error) {
	return yaml.Marshal(data)
}

// Decode decodes cache data from YAML format.
func (*YAML[K, V]) Decode(data []byte) (map[K]V, error) {
	m := map[K]V{}
	err := yaml.Unmarshal(data, &m)
	return m, err
}

// TOML encodes/decodes cache data in TOML format.
type TOML[K comparable, V any] struct {
	buffer  bytes.Buffer
	encoder *toml.Encoder
}

// Encode encodes cache data in TOML format.
func (t *TOML[K, V]) Encode(data map[K]V) ([]byte, error) {
	t.buffer.Reset()
	if t.encoder == nil {
		t.encoder = toml.NewEncoder(&t.buffer)
	}
	err := t.encoder.Encode(data)
	return t.buffer.Bytes(), err
}

// Decode decodes cache data from TOML format.
func (*TOML[K, V]) Decode(data []byte) (map[K]V, error) {
	m := map[K]V{}
	_, err := toml.Decode(string(data), &m)
	return m, err
}

// GOB encodes/decodes cache data in self-describing binary format.
type GOB[K comparable, V any] struct{}

// Encode encodes cache data in self-describing binary format.
func (g *GOB[K, V]) Encode(data map[K]V) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(&data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Decode decodes cache data from self-describing binary format.
func (g *GOB[K, V]) Decode(data []byte) (map[K]V, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	m := map[K]V{}
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
