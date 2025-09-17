/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec_test

import (
	"fmt"
	"testing"

	"go.osspkg.com/casecheck"

	"go.osspkg.com/ioutils/codec"
)

func TestFile_Blob_EncodeDecode1(t *testing.T) {
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
		Ext:  codec.ExtJSON,
	}
	casecheck.NoError(t, b.Encode(model1, model2))
	casecheck.Equal(t, `{"Data1":{"AA":"123","BB":true},"Data2":{"CC":"qwer","DD":-100}}`, string(b.Blob))

	b = &codec.BlobEncoder{
		Blob: nil,
		Ext:  codec.ExtYAML,
	}
	casecheck.NoError(t, b.Encode(model1, model2))
	casecheck.Equal(t, "data-1:\n    aa: \"123\"\n    bb: true\n\ndata-2:\n    cc: qwer\n    dd: -100", string(b.Blob))

	b = &codec.BlobEncoder{
		Blob: nil,
		Ext:  codec.ExtToml,
	}
	casecheck.NoError(t, b.Encode(model1, model2))
	casecheck.Equal(t, "[Data1]\n  AA = \"123\"\n  BB = true\n\n[Data2]\n  CC = \"qwer\"\n  DD = -100", string(b.Blob))

	b = &codec.BlobEncoder{
		Blob: nil,
		Ext:  codec.ExtXML,
	}
	casecheck.NoError(t, b.Encode(model1, model2))
	fmt.Printf("%#v\n", string(b.Blob))
	casecheck.Equal(t, "<TestData1><Data1><AA>123</AA><BB>true</BB></Data1></TestData1>\n\n<TestData2><Data2><CC>qwer</CC><DD>-100</DD></Data2></TestData2>", string(b.Blob))
}

func TestFile_Blob_EncodeDecode2_DuplicatKey(t *testing.T) {
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
		Data1 TestDataItem2 `yaml:"data-1"`
	}

	model1 := &TestData1{Data1: TestDataItem1{AA: "123", BB: true}}
	model2 := &TestData2{Data1: TestDataItem2{CC: "qwer", DD: -100}}

	be := &codec.BlobEncoder{
		Blob: nil,
		Ext:  codec.ExtYAML,
	}
	casecheck.NoError(t, be.Encode(model1, model2))
	casecheck.Equal(t, "data-1:\n    aa: \"123\"\n    bb: true\n\ndata-1:\n    cc: qwer\n    dd: -100", string(be.Blob))

	model01 := &TestData1{}
	model02 := &TestData2{}
	casecheck.Error(t, be.Decode(model01, model02))
}
