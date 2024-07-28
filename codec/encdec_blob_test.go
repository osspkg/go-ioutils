/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec_test

import (
	"testing"

	"go.osspkg.com/casecheck"
	"go.osspkg.com/ioutils/codec"
)

func TestFile_Blob_EncodeDecode(t *testing.T) {
	type TestDataItem1 struct {
		AA string `yaml:"aa"`
		BB bool   `yaml:"bb"`
	}
	type TestData1 struct {
		Data1 TestDataItem1 `yaml:"data-1"`
	}
	type TestDataItem2 struct {
		CC string `yaml:"cc"`
		DD int    `yaml:"dd"`
	}
	type TestData2 struct {
		Data2 TestDataItem2 `yaml:"data-2"`
	}

	model1 := &TestData1{Data1: TestDataItem1{AA: "123", BB: true}}
	model2 := &TestData2{Data2: TestDataItem2{CC: "qwer", DD: -100}}

	b := &codec.BlobEncoder{
		Blob: nil,
		Ext:  codec.EncoderJSON,
	}
	casecheck.NoError(t, b.Encode(model1, model2))
	casecheck.Equal(t, `{"Data1":{"AA":"123","BB":true},"Data2":{"CC":"qwer","DD":-100}}`, string(b.Blob))

	b = &codec.BlobEncoder{
		Blob: nil,
		Ext:  codec.EncoderYAML,
	}
	casecheck.NoError(t, b.Encode(model1, model2))
	casecheck.Equal(t, "data-1:\n    aa: \"123\"\n    bb: true\ndata-2:\n    cc: qwer\n    dd: -100\n", string(b.Blob))
}
