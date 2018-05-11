/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package stateleveldb

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/common/ledger/testutil"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
)

// TestEncodeString tests encoding and decoding a string value using old format
func TestEncodeDecodeStringOldFormat(t *testing.T) {
	bytesString1 := []byte("value1")
	version1 := version.NewHeight(1, 1)
	encodedValue := encodeValueOldFormat(bytesString1, version1)
	decodedValue, err := decodeValue(encodedValue)
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, decodedValue, &statedb.VersionedValue{Version: version1, Value: bytesString1})
}

// TestEncodeDecodeJSONOldFormat tests encoding and decoding a JSON value using old format
func TestEncodeDecodeJSONOldFormat(t *testing.T) {
	bytesJSON2 := []byte(`{"asset_name":"marble1","color":"blue","size":"35","owner":"jerry"}`)
	version2 := version.NewHeight(1, 1)
	encodedValue := encodeValueOldFormat(bytesJSON2, version2)
	decodedValue, err := decodeValue(encodedValue)
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, decodedValue, &statedb.VersionedValue{Version: version2, Value: bytesJSON2})
}

func TestEncodeDecodeOldAndNewFormat(t *testing.T) {
	testdata := []*statedb.VersionedValue{
		{
			Value:   []byte("value1"),
			Version: version.NewHeight(1, 2),
		},

		{
			Value:   []byte{},
			Version: version.NewHeight(50, 50),
		},
		{
			Value:    []byte{},
			Version:  version.NewHeight(50, 50),
			Metadata: []byte("sample-metadata"),
		},
	}

	for i, testdatum := range testdata {
		t.Run(fmt.Sprintf("testcase-newfmt-%d", i),
			func(t *testing.T) { testEncodeDecodeNewFormat(t, testdatum) },
		)
	}

	for i, testdatum := range testdata {
		t.Run(fmt.Sprintf("testcase-oldfmt-%d", i),
			func(t *testing.T) { testEncodeDecodeOldFormat(t, testdatum) },
		)
	}

}

func testEncodeDecodeNewFormat(t *testing.T, v *statedb.VersionedValue) {
	encodedNewFmt, err := encodeValue(v)
	testutil.AssertNoError(t, err, "")
	// encoding-decoding using new format should return the same versioned_value
	decodedFromNewFmt, err := decodeValue(encodedNewFmt)
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, decodedFromNewFmt, v)
}

func testEncodeDecodeOldFormat(t *testing.T, v *statedb.VersionedValue) {
	encodedOldFmt := encodeValueOldFormat(v.Value, v.Version)
	// decodeValue should be able to handle the old format
	decodedFromOldFmt, err := decodeValue(encodedOldFmt)
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, decodedFromOldFmt.Value, v.Value)
	testutil.AssertEquals(t, decodedFromOldFmt.Version, v.Version)
	testutil.AssertNil(t, decodedFromOldFmt.Metadata)
}
