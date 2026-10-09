/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"path/filepath"
	"testing"

	"go.osspkg.com/casecheck"
)

func TestFile_File_EncodeDecode(t *testing.T) {
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

	filename := filepath.Join(t.TempDir(), "config.yaml")

	model1 := &TestData1{Data1: TestDataItem1{AA: "123", BB: true}}
	model2 := &TestData2{Data2: TestDataItem2{CC: "qwer", DD: -100}}

	err := FileEncoder(filename).Encode(model1, model2)
	casecheck.NoError(t, err)

	model11 := &TestData1{}
	model22 := &TestData2{}

	err = FileEncoder(filename).Decode(model11, model22)
	casecheck.NoError(t, err)

	casecheck.Equal(t, model1, model11)
	casecheck.Equal(t, model2, model22)
}
