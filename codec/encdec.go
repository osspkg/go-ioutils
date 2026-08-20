/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"encoding/json"
	"encoding/xml"

	"github.com/BurntSushi/toml"
	"go.osspkg.com/errors"
	"go.osspkg.com/syncing"
	"go.osspkg.com/unic"
	"gopkg.in/yaml.v3"
)

const (
	ExtYAMLs = ".yml"
	ExtYAML  = ".yaml"
	ExtJSON  = ".json"
	ExtToml  = ".toml"
	ExtXML   = ".xml"
	ExtUnic  = ".unic"
	ExtConf  = ".conf"
)

var (
	ErrUnsupportedFormat = errors.New("format is not a supported")

	_default = newEncoders().
			Add(ExtYAMLs, SimpleCodec(yaml.Marshal, yaml.Unmarshal, BytesJoin)).
			Add(ExtYAML, SimpleCodec(yaml.Marshal, yaml.Unmarshal, BytesJoin)).
			Add(ExtJSON, SimpleCodec(json.Marshal, json.Unmarshal, MapJoin)).
			Add(ExtToml, SimpleCodec(toml.Marshal, toml.Unmarshal, BytesJoin)).
			Add(ExtXML, SimpleCodec(xml.Marshal, xml.Unmarshal, BytesJoin)).
			Add(ExtUnic, Codec{Encode: unic.Marshal, Decode: unic.Unmarshal}).
			Add(ExtConf, Codec{Encode: unic.Marshal, Decode: unic.Unmarshal})
)

type (
	Codec struct {
		Encode func(args ...any) ([]byte, error)
		Decode func(b []byte, args ...any) error
	}
	encoders struct {
		list map[string]Codec
		mux  syncing.Lock
	}
)

func AddCodec(ext string, c Codec) {
	_default.Add(ext, c)
}

func newEncoders() *encoders {
	return &encoders{
		list: make(map[string]Codec, 10),
		mux:  syncing.NewLock(),
	}
}

func (v *encoders) Add(ext string, c Codec) *encoders {
	v.mux.Lock(func() {
		v.list[ext] = c
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
