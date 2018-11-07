// Code generated by gocc; DO NOT EDIT.

package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
)

const numNTSymbols = 22

type (
	gotoTable [numStates]gotoRow
	gotoRow   [numNTSymbols]int
)

var gotoTab = gotoTable{}

func init() {
	tab := [][]int{}
	data := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xec, 0x97, 0x51, 0x8b, 0x12, 0x51,
		0x14, 0xc7, 0xef, 0x7f, 0x90, 0x08, 0x31, 0x91, 0x28, 0xf1, 0x13, 0x44, 0x8f, 0x3d, 0x84, 0x88,
		0x44, 0x88, 0x44, 0x44, 0x48, 0x44, 0x84, 0x88, 0x4f, 0x3d, 0x45, 0x44, 0x44, 0xf4, 0x1c, 0x12,
		0x11, 0x51, 0xe6, 0x73, 0x0f, 0x11, 0x22, 0x22, 0x22, 0x22, 0x22, 0x12, 0x7d, 0x92, 0x9e, 0xa2,
		0x22, 0x62, 0xd9, 0x8f, 0xb0, 0xec, 0xc3, 0x72, 0xce, 0xa2, 0x8e, 0xeb, 0x0a, 0xf7, 0x9c, 0xd1,
		0x9d, 0x99, 0x65, 0x76, 0x99, 0x7d, 0x98, 0x61, 0x7f, 0xbf, 0xff, 0x1c, 0x0f, 0xce, 0xbd, 0xf7,
		0xe0, 0x25, 0xfe, 0xec, 0xc0, 0xe1, 0x96, 0x01, 0x37, 0x8d, 0x49, 0xf1, 0xa7, 0xd9, 0x7f, 0x4d,
		0x83, 0x84, 0x31, 0x74, 0xf5, 0x0a, 0xb7, 0x0c, 0xff, 0xc9, 0xc2, 0x81, 0xe5, 0x2f, 0x6b, 0x83,
		0x27, 0xc0, 0xc9, 0x54, 0x1a, 0x19, 0x4b, 0x3a, 0x67, 0xc7, 0x21, 0x76, 0x12, 0x31, 0x7c, 0xfd,
		0x5a, 0x00, 0x45, 0x6e, 0x6c, 0xdb, 0x49, 0x61, 0x7e, 0xcd, 0x17, 0x6f, 0xdd, 0x2e, 0x95, 0x57,
		0xb8, 0xe2, 0xbf, 0x13, 0x05, 0xd7, 0x50, 0x5d, 0xc7, 0xf5, 0xe0, 0x6a, 0x47, 0x1d, 0x37, 0xe6,
		0xd7, 0x37, 0xfc, 0x96, 0xdf, 0xf1, 0x7b, 0xfe, 0xc0, 0x1f, 0x23, 0xd9, 0x25, 0xc0, 0x5f, 0xfc,
		0x57, 0x29, 0xb8, 0xa5, 0xbe, 0xad, 0x2f, 0xae, 0x23, 0xde, 0xb6, 0x73, 0x70, 0x47, 0xe0, 0x5d,
		0x81, 0xf7, 0x04, 0xde, 0x17, 0xf8, 0x40, 0xe0, 0x43, 0x3b, 0x07, 0x8f, 0x04, 0x3e, 0xb6, 0x73,
		0xf0, 0x44, 0xe0, 0xd3, 0xb2, 0x9f, 0x6d, 0xc3, 0x3f, 0xb6, 0x88, 0x9f, 0x2f, 0xbc, 0xd8, 0x37,
		0xfc, 0xd3, 0xba, 0x71, 0xf8, 0xf7, 0xe2, 0xf6, 0x8b, 0xff, 0xf2, 0x3f, 0xfe, 0xcf, 0x3b, 0xbc,
		0x1b, 0x46, 0x5b, 0xbc, 0xcf, 0x7b, 0x01, 0x94, 0x89, 0xb1, 0xed, 0xd5, 0x82, 0x0f, 0x2c, 0x6f,
		0x76, 0xe9, 0x68, 0xf3, 0xe3, 0xd2, 0x7d, 0x04, 0x04, 0x23, 0x95, 0x03, 0xc1, 0x51, 0x5c, 0x42,
		0x71, 0x17, 0x14, 0x77, 0x51, 0x71, 0x49, 0xd1, 0x81, 0x90, 0x92, 0x55, 0x5a, 0x52, 0x20, 0x64,
		0x44, 0x73, 0x39, 0x84, 0xb9, 0x52, 0x03, 0x21, 0xe7, 0xbf, 0xcc, 0x19, 0xc5, 0x8b, 0x6f, 0x97,
		0x70, 0xd3, 0x3e, 0xb9, 0x09, 0x79, 0xeb, 0x63, 0x84, 0xe2, 0xd6, 0x9f, 0x4e, 0x28, 0xf9, 0xed,
		0x36, 0xc6, 0x21, 0x61, 0x77, 0xd2, 0x80, 0x70, 0xc7, 0x32, 0x69, 0x56, 0xf6, 0xae, 0x62, 0x41,
		0xb8, 0xa7, 0xc9, 0xfb, 0x9a, 0xac, 0x68, 0xf2, 0x81, 0x26, 0x1f, 0x6a, 0xf2, 0x91, 0x2c, 0x41,
		0x78, 0xac, 0xb8, 0xaa, 0xe8, 0x40, 0xa8, 0xc9, 0xaa, 0x1e, 0xc2, 0x8c, 0x9e, 0x1d, 0x51, 0x4f,
		0x36, 0xcf, 0x13, 0x9e, 0xb9, 0xf7, 0xa7, 0x84, 0xe7, 0x84, 0x17, 0x84, 0x97, 0x84, 0x57, 0x84,
		0xd7, 0x91, 0x5c, 0x7b, 0x31, 0xde, 0x1c, 0xcf, 0x56, 0x42, 0x23, 0x22, 0xbd, 0x9c, 0x3e, 0x5e,
		0x0e, 0xab, 0xaf, 0x01, 0xfc, 0xcc, 0x5c, 0x6e, 0x12, 0x10, 0xda, 0xc2, 0x16, 0x39, 0x16, 0xe9,
		0x78, 0x45, 0x40, 0xe8, 0x7a, 0x05, 0x7a, 0x5e, 0x81, 0xbe, 0x57, 0x60, 0xe0, 0x15, 0x18, 0x7a,
		0x05, 0x46, 0x7a, 0x00, 0x84, 0xb1, 0xae, 0x27, 0x9a, 0x06, 0x61, 0xaa, 0xc9, 0xef, 0xf1, 0x19,
		0x14, 0x15, 0x7c, 0x18, 0x00, 0x00, 0xff, 0xff, 0x0a, 0x64, 0x14, 0x88, 0x32, 0x15, 0x00, 0x00,
	}
	buf, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&tab); err != nil {
		panic(err)
	}
	for i := 0; i < numStates; i++ {
		for j := 0; j < numNTSymbols; j++ {
			gotoTab[i][j] = tab[i][j]
		}
	}
}
