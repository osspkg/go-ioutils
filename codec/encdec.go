/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"encoding/json"

	"go.osspkg.com/errors"
	"go.osspkg.com/syncing"
	"gopkg.in/yaml.v3"
)

const (
	EncoderYAML = ".yaml"
	EncoderJSON = ".json"
)

var (
	errBadFormat = errors.New("format is not a supported")

	_default = newEncoders().
			Add(".yml", yaml.Marshal, yaml.Unmarshal, mapMerge).
			Add(EncoderYAML, yaml.Marshal, yaml.Unmarshal, mapMerge).
			Add(EncoderJSON, json.Marshal, json.Unmarshal, mapMerge)
)

type (
	Codec struct {
		Encode func(in interface{}) ([]byte, error)
		Decode func(b []byte, out interface{}) error
		Merge  func(dst map[string]interface{}, src map[string]interface{}) error
	}
	encoders struct {
		list map[string]Codec
		mux  syncing.Lock
	}
)

func newEncoders() *encoders {
	return &encoders{
		list: make(map[string]Codec, 10),
		mux:  syncing.NewLock(),
	}
}

func AddCodec(ext string, c Codec) {
	_default.Add(ext, c.Encode, c.Decode, c.Merge)
}

func (v *encoders) Add(
	ext string,
	enc func(interface{}) ([]byte, error),
	dec func([]byte, interface{}) error,
	merge func(map[string]interface{}, map[string]interface{}) error,
) *encoders {
	v.mux.Lock(func() {
		v.list[ext] = Codec{
			Encode: enc,
			Decode: dec,
			Merge:  merge,
		}
	})
	return v
}

func (v *encoders) Get(ext string) (c Codec, err error) {
	v.mux.RLock(func() {
		var ok bool
		if c, ok = v.list[ext]; !ok {
			err = errBadFormat
			return
		}
	})
	return
}
