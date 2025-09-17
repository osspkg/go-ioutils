/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"encoding/json"
	"encoding/xml"

	"github.com/BurntSushi/toml"
	"go.osspkg.com/errors"
	"go.osspkg.com/syncing"
	"gopkg.in/yaml.v3"
)

const (
	ExtYAML = ".yaml"
	ExtJSON = ".json"
	ExtToml = ".toml"
	ExtXML  = ".xml"
)

var (
	ErrUnsupportedFormat = errors.New("format is not a supported")

	_default = newEncoders().
			Add(".yml", yaml.Marshal, yaml.Unmarshal, BytesJoin).
			Add(ExtYAML, yaml.Marshal, yaml.Unmarshal, BytesJoin).
			Add(ExtJSON, json.Marshal, json.Unmarshal, MapJoin).
			Add(ExtToml, toml.Marshal, toml.Unmarshal, BytesJoin).
			Add(ExtXML, xml.Marshal, xml.Unmarshal, BytesJoin)
)

type (
	Codec struct {
		Encode func(in interface{}) ([]byte, error)
		Decode func(b []byte, out interface{}) error
		Join   func(c Codec, dst *[]byte, src ...[]byte) error
	}
	encoders struct {
		list map[string]Codec
		mux  syncing.Lock
	}
)

func AddCodec(ext string, c Codec) {
	_default.Add(ext, c.Encode, c.Decode, c.Join)
}

func newEncoders() *encoders {
	return &encoders{
		list: make(map[string]Codec, 10),
		mux:  syncing.NewLock(),
	}
}

func (v *encoders) Add(
	ext string,
	enc func(interface{}) ([]byte, error),
	dec func([]byte, interface{}) error,
	join func(c Codec, dst *[]byte, src ...[]byte) error,
) *encoders {
	v.mux.Lock(func() {
		v.list[ext] = Codec{
			Encode: enc,
			Decode: dec,
			Join:   join,
		}
	})
	return v
}

func (v *encoders) Get(ext string) (c Codec, err error) {
	v.mux.RLock(func() {
		var ok bool
		if c, ok = v.list[ext]; !ok {
			err = ErrUnsupportedFormat
			return
		}
	})
	return
}
