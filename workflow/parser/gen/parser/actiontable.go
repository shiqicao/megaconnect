// Code generated by gocc; DO NOT EDIT.

package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
)

type (
	actionTable [numStates]actionRow
	actionRow   struct {
		canRecover bool
		actions    [numSymbols]action
	}
)

var actionTab = actionTable{}

func init() {
	tab := []struct {
		CanRecover bool
		Actions    []struct {
			Index  int
			Action int
			Amount int
		}
	}{}
	data := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xec, 0x9a, 0xef, 0x8b, 0x5b, 0x59,
		0xfd, 0xc7, 0xcf, 0xe7, 0xf6, 0x7e, 0xf3, 0x1d, 0xaf, 0xb3, 0xb3, 0x43, 0x09, 0xc3, 0x30, 0x0c,
		0xc3, 0x38, 0x0c, 0x43, 0x29, 0xa5, 0x94, 0x52, 0xca, 0x50, 0x6b, 0xe8, 0xd6, 0xb1, 0xd4, 0x52,
		0xca, 0x32, 0x0c, 0x25, 0xc4, 0x18, 0x62, 0x9a, 0x86, 0x6c, 0x88, 0x69, 0xc8, 0xde, 0x0d, 0x21,
		0xc6, 0xb2, 0xae, 0x55, 0xcb, 0x22, 0xb2, 0x2c, 0xcb, 0xb2, 0x88, 0x88, 0x88, 0x48, 0x91, 0xc5,
		0x47, 0x3e, 0x12, 0x1f, 0x88, 0x0f, 0xc4, 0x07, 0xe2, 0x5f, 0x20, 0xb2, 0x2c, 0xcb, 0x22, 0x22,
		0x22, 0x22, 0x22, 0x3e, 0xc8, 0x95, 0x7b, 0x5e, 0xe7, 0x9c, 0x4d, 0xb3, 0x37, 0x85, 0x26, 0xb9,
		0xd9, 0x82, 0x43, 0xa1, 0xe7, 0x70, 0x7f, 0xbc, 0x5f, 0x9f, 0x9c, 0xcc, 0x3d, 0x9f, 0xf7, 0xe7,
		0x73, 0xf3, 0x5c, 0xf4, 0x1d, 0x4f, 0xbc, 0xe8, 0xa1, 0x92, 0xe8, 0x35, 0xa5, 0x4e, 0x45, 0xdf,
		0x38, 0x21, 0x5e, 0xf4, 0x9a, 0x12, 0x4f, 0x82, 0xcf, 0x7f, 0xa5, 0x75, 0x78, 0xf7, 0xce, 0xbd,
		0xee, 0xdd, 0x8e, 0x78, 0x4a, 0xfe, 0xff, 0x85, 0x3b, 0xe1, 0x4b, 0xf7, 0x5a, 0x2f, 0x4b, 0xf4,
		0x6d, 0xa5, 0xd4, 0xe7, 0xa2, 0x6f, 0x79, 0x22, 0x67, 0xbf, 0xf4, 0xe5, 0x97, 0xc3, 0xce, 0x2b,
		0x77, 0xc2, 0xed, 0xaf, 0x6d, 0x7f, 0xb1, 0x55, 0xbd, 0xdb, 0xdb, 0x7e, 0xa9, 0x15, 0x7e, 0x76,
		0x9b, 0x2b, 0xcd, 0xfc, 0xab, 0xf7, 0x5e, 0x69, 0x85, 0xf1, 0x7c, 0xfb, 0xeb, 0xf1, 0xad, 0x12,
		0x3d, 0x50, 0xea, 0x4c, 0xf4, 0xcd, 0x18, 0xf3, 0x40, 0xc9, 0x09, 0xf9, 0x3f, 0x7d, 0xa3, 0xf8,
		0x4a, 0x32, 0xdc, 0xc7, 0x54, 0xdf, 0x26, 0xbe, 0x52, 0x6a, 0xf8, 0x85, 0x17, 0xa2, 0x87, 0x6a,
		0x28, 0x9f, 0xf1, 0x44, 0xfc, 0xf8, 0x9f, 0x52, 0x9e, 0x88, 0xa7, 0xff, 0xcf, 0x88, 0x2f, 0x19,
		0x3d, 0x5b, 0x12, 0x5f, 0x96, 0x94, 0xf2, 0x7c, 0x09, 0xc4, 0x93, 0x15, 0x25, 0x2b, 0xe2, 0x4b,
		0x56, 0xc9, 0xba, 0xf8, 0xb2, 0xa6, 0x64, 0x4b, 0x7c, 0x59, 0x37, 0xa7, 0x7d, 0xd9, 0x78, 0xd2,
		0x69, 0x2f, 0x0e, 0x60, 0x85, 0x61, 0x9d, 0x61, 0x4b, 0x0f, 0xf6, 0xf4, 0x12, 0xa7, 0x97, 0x38,
		0xbd, 0xc4, 0xe9, 0x8f, 0xd8, 0x01, 0xa7, 0x03, 0x4e, 0x07, 0x9c, 0x0e, 0xdc, 0xe9, 0x65, 0x4e,
		0x2f, 0x73, 0x7a, 0x99, 0xd3, 0xcb, 0xee, 0xd3, 0x6c, 0xb9, 0xd9, 0xb6, 0x9b, 0xed, 0xf0, 0x89,
		0xe3, 0x7f, 0x4e, 0x26, 0x83, 0x4c, 0x06, 0x99, 0x0c, 0x32, 0x2c, 0xc5, 0xaa, 0xf8, 0xb2, 0xeb,
		0x16, 0x65, 0x4f, 0xcf, 0xb6, 0xc5, 0x97, 0x53, 0x4e, 0xf0, 0xb4, 0x52, 0x9e, 0xa7, 0x67, 0x67,
		0x94, 0x56, 0xfb, 0x88, 0x75, 0x5e, 0xcf, 0x4e, 0x8a, 0x2f, 0xfb, 0x7a, 0xb6, 0x21, 0xbe, 0x5c,
		0xd2, 0xb3, 0x78, 0xe1, 0x2e, 0xeb, 0x1b, 0xe3, 0x3b, 0x36, 0x95, 0x6c, 0xea, 0x41, 0x79, 0x27,
		0x64, 0x47, 0x3c, 0xb9, 0xa0, 0x64, 0x97, 0x61, 0x4f, 0x0f, 0x1c, 0xf7, 0xe5, 0x4a, 0x7c, 0xdc,
		0x97, 0xab, 0xf1, 0x71, 0x5f, 0x0e, 0xdc, 0xf5, 0xe7, 0xb9, 0xfe, 0x3c, 0xd7, 0xc7, 0xd8, 0x8c,
		0x0e, 0xe0, 0x9a, 0xd2, 0x61, 0x5f, 0x57, 0x72, 0x4d, 0x7c, 0xc9, 0x2b, 0xb9, 0x2e, 0xbe, 0x14,
		0x94, 0xdc, 0x10, 0x5f, 0x8a, 0x4a, 0x6e, 0x8a, 0x2f, 0x25, 0xbd, 0x0a, 0xf1, 0x55, 0x65, 0x25,
		0xa7, 0xc4, 0x97, 0xaa, 0x92, 0xd3, 0xe2, 0x4b, 0x4d, 0xc9, 0x19, 0xf1, 0xa5, 0xee, 0x16, 0x69,
		0x8d, 0x45, 0x5a, 0x63, 0x91, 0xd6, 0x58, 0xa4, 0x35, 0xf7, 0x19, 0xb6, 0xe2, 0xcf, 0xe0, 0x4b,
		0xc3, 0xad, 0x55, 0x73, 0x6c, 0x19, 0xec, 0xec, 0x79, 0xc9, 0x8a, 0x27, 0xf7, 0x95, 0xac, 0x31,
		0xec, 0x32, 0xec, 0x31, 0x9c, 0x63, 0x38, 0xcf, 0x70, 0x81, 0xe1, 0x22, 0xc3, 0x3e, 0xc3, 0x25,
		0x86, 0xcb, 0x0c, 0x39, 0x86, 0x2b, 0x0c, 0x57, 0x19, 0x0e, 0xf4, 0x60, 0xbf, 0x97, 0x0e, 0xdf,
		0xcb, 0x75, 0xfd, 0x51, 0x62, 0x78, 0x36, 0x86, 0xfb, 0xd2, 0x67, 0x35, 0x07, 0xac, 0xe6, 0x7d,
		0x17, 0xdb, 0x80, 0xd8, 0x06, 0xc4, 0x36, 0x20, 0xb6, 0x01, 0xb1, 0x0d, 0x88, 0x6d, 0x40, 0x6c,
		0x03, 0x62, 0x1b, 0x10, 0xdb, 0x80, 0xd8, 0x06, 0xc4, 0x36, 0x20, 0xb6, 0x01, 0xb1, 0x0d, 0x88,
		0x6d, 0x40, 0x6c, 0x03, 0xa5, 0xbc, 0x40, 0xb3, 0x6e, 0xc3, 0xba, 0x0d, 0xeb, 0x36, 0xac, 0xdb,
		0x31, 0xcb, 0x97, 0xe8, 0xd5, 0x18, 0xe6, 0xc7, 0x9b, 0x89, 0x5c, 0x88, 0xc7, 0x07, 0x31, 0xce,
		0xd7, 0x4f, 0xfe, 0x7e, 0x3c, 0x3e, 0x8c, 0x81, 0xbe, 0x44, 0xaf, 0x2b, 0xe5, 0x2d, 0x6b, 0xb9,
		0x2a, 0x72, 0x55, 0xe4, 0xaa, 0xc8, 0x55, 0x09, 0xbd, 0x4a, 0xe8, 0x55, 0x42, 0xaf, 0x12, 0x7a,
		0x95, 0xd0, 0xab, 0x84, 0x5e, 0x8d, 0x43, 0xf7, 0x25, 0xfa, 0x6e, 0x1c, 0xbb, 0x2f, 0xd1, 0xf7,
		0x94, 0xf2, 0x56, 0xb4, 0x70, 0x03, 0xe1, 0x06, 0xc2, 0x0d, 0x84, 0x1b, 0x08, 0x37, 0x10, 0x6e,
		0x20, 0xdc, 0x40, 0xb8, 0x81, 0x70, 0x03, 0xe1, 0x06, 0x6b, 0xd2, 0x60, 0x4d, 0x1a, 0xf1, 0x9a,
		0xf8, 0x12, 0xbd, 0x11, 0x2f, 0x8a, 0x2f, 0xd1, 0x9b, 0x6e, 0xe9, 0xdb, 0x60, 0xda, 0x60, 0xda,
		0x60, 0xda, 0x60, 0xda, 0x60, 0xda, 0x60, 0xda, 0x60, 0xda, 0x60, 0xda, 0x60, 0xda, 0x60, 0xda,
		0x60, 0xda, 0x2c, 0x7d, 0x9b, 0xa5, 0x6f, 0xc7, 0x4b, 0xef, 0x4b, 0xf4, 0x96, 0x83, 0x85, 0xc0,
		0x42, 0x60, 0x21, 0xb0, 0x10, 0x58, 0x08, 0x2c, 0x04, 0x16, 0x02, 0x0b, 0x81, 0x85, 0xc0, 0x42,
		0x60, 0x21, 0xb0, 0x10, 0x58, 0x08, 0x2c, 0xe4, 0x7b, 0x0e, 0x1d, 0xab, 0x0b, 0xab, 0x0b, 0xab,
		0x0b, 0xab, 0x0b, 0xab, 0x0b, 0xab, 0x0b, 0xab, 0x0b, 0xab, 0x0b, 0xab, 0x0b, 0xab, 0x0b, 0xab,
		0x0b, 0xab, 0x0b, 0xab, 0x0b, 0xab, 0x0b, 0xab, 0xeb, 0x58, 0x3d, 0x58, 0x3d, 0x58, 0x3d, 0x58,
		0x3d, 0x58, 0x3d, 0x58, 0x3d, 0x58, 0x3d, 0x58, 0x3d, 0x58, 0x3d, 0x58, 0x3d, 0x58, 0x3d, 0x58,
		0x3d, 0x58, 0x3d, 0x58, 0x3d, 0x58, 0x3d, 0xc7, 0xea, 0xc3, 0xea, 0xc3, 0xea, 0xc3, 0xea, 0xc3,
		0xea, 0xc3, 0xea, 0xc3, 0xea, 0xc3, 0xea, 0xc3, 0xea, 0xc3, 0xea, 0xc3, 0xea, 0xc3, 0xea, 0xc3,
		0xea, 0xc3, 0xea, 0xc3, 0xea, 0xbb, 0xdd, 0x2c, 0x7a, 0x9b, 0xed, 0x2c, 0x7a, 0x87, 0xfd, 0x2c,
		0x7a, 0xc4, 0x86, 0x16, 0xfd, 0x8c, 0x1d, 0x2d, 0x7a, 0x97, 0x2d, 0x2d, 0xfa, 0x79, 0xd2, 0x96,
		0xcc, 0x36, 0xb5, 0xce, 0x56, 0xbb, 0xee, 0x0e, 0x5c, 0xe4, 0xc0, 0x45, 0x77, 0x60, 0x9f, 0x03,
		0xfb, 0xee, 0xc0, 0x25, 0x0e, 0x5c, 0x72, 0xbb, 0xd8, 0x19, 0x77, 0xea, 0x8c, 0x92, 0x6c, 0x0c,
		0xfc, 0x85, 0xdb, 0x89, 0xcf, 0xb2, 0x10, 0x67, 0x59, 0x88, 0xb3, 0xee, 0xf8, 0x39, 0x8e, 0x9f,
		0xe3, 0xf8, 0x39, 0x97, 0x0e, 0xa2, 0x5f, 0xba, 0x7c, 0x10, 0xfd, 0xca, 0xe9, 0x1e, 0x80, 0x3c,
		0xd0, 0xe7, 0x34, 0xe1, 0xd7, 0x6e, 0x5f, 0x8d, 0x7e, 0xf3, 0xb4, 0xdb, 0xfb, 0xf1, 0xc5, 0xff,
		0x53, 0x17, 0xf3, 0x47, 0x1a, 0xfd, 0x51, 0xef, 0xdd, 0xe9, 0x65, 0x57, 0xb9, 0x35, 0x29, 0xc9,
		0x9e, 0xd0, 0x69, 0x35, 0x7a, 0x8f, 0xbc, 0x1a, 0xbd, 0x1f, 0x5f, 0xea, 0x4b, 0xf4, 0x81, 0x0b,
		0x28, 0x9d, 0x94, 0xaa, 0x03, 0x8a, 0x33, 0xeb, 0xa7, 0x92, 0x72, 0xe9, 0x87, 0x26, 0x97, 0xfe,
		0xd9, 0xe4, 0xd2, 0xbf, 0x98, 0x5c, 0xfa, 0x57, 0x93, 0x4b, 0xff, 0x66, 0x72, 0xe9, 0xdf, 0xd1,
		0xb9, 0xad, 0x94, 0xf7, 0xe9, 0xd9, 0x92, 0xe8, 0x3f, 0x4c, 0x12, 0xfd, 0x27, 0x8a, 0x55, 0xa5,
		0xbc, 0xe7, 0xe6, 0x99, 0x3d, 0xff, 0x65, 0xb2, 0xe7, 0xbf, 0xd1, 0x6f, 0xb8, 0xe5, 0x4d, 0x25,
		0x6d, 0xfe, 0x07, 0x4a, 0xdb, 0x51, 0xd2, 0xc9, 0x97, 0x1a, 0x12, 0x3a, 0x48, 0x3a, 0x89, 0x52,
		0x43, 0xba, 0x0e, 0x92, 0x4e, 0x86, 0xd4, 0x90, 0x9e, 0x83, 0xa4, 0x93, 0x1a, 0x35, 0x64, 0xca,
		0x0c, 0x69, 0xb2, 0xce, 0xd0, 0x65, 0x9d, 0x0d, 0xb2, 0xce, 0x86, 0xdb, 0x43, 0x86, 0xa2, 0xdc,
		0x75, 0x43, 0x5d, 0x2f, 0x71, 0xe1, 0x29, 0xa5, 0x73, 0xff, 0x29, 0xc7, 0x1d, 0xc6, 0x45, 0xd3,
		0x92, 0x9e, 0x2c, 0x41, 0x1e, 0xc6, 0x99, 0xf6, 0xba, 0x9e, 0x6c, 0xc0, 0x1e, 0xc6, 0x65, 0xce,
		0x4d, 0x3d, 0xd9, 0x72, 0xf6, 0x21, 0x87, 0x7d, 0xc8, 0x91, 0x1d, 0x73, 0xac, 0x51, 0x8e, 0x35,
		0xca, 0xb1, 0x46, 0x39, 0xd6, 0x28, 0xc7, 0x1a, 0xe5, 0x58, 0xa3, 0x1c, 0x6b, 0x94, 0x63, 0x8d,
		0x72, 0xac, 0x51, 0x8e, 0x35, 0xca, 0xb1, 0x46, 0x39, 0xd6, 0x28, 0xe7, 0xa2, 0xbe, 0x46, 0x85,
		0x32, 0x1c, 0x29, 0xd0, 0x86, 0xba, 0x1c, 0xf4, 0xcc, 0x7c, 0x97, 0x2d, 0xec, 0xa6, 0x33, 0xe7,
		0x87, 0x44, 0x77, 0x48, 0x74, 0x87, 0x44, 0x77, 0x38, 0x85, 0x39, 0x47, 0xee, 0x08, 0xb9, 0x23,
		0xe4, 0x8e, 0x90, 0x3b, 0x9a, 0xda, 0xeb, 0xe7, 0x91, 0xcb, 0x23, 0x97, 0x47, 0x2e, 0xcf, 0xda,
		0xe5, 0x59, 0xbb, 0x3c, 0x6b, 0x97, 0x67, 0xed, 0xf2, 0xac, 0x5d, 0x9e, 0xb5, 0xcb, 0x27, 0x78,
		0x7d, 0x84, 0x0b, 0x08, 0x17, 0x10, 0x2e, 0x20, 0x5c, 0x40, 0xb8, 0x80, 0x70, 0x01, 0xe1, 0x02,
		0xc2, 0x05, 0x84, 0x0b, 0x08, 0x17, 0x26, 0x0a, 0x17, 0x11, 0x2e, 0x22, 0x5c, 0x44, 0xb8, 0x88,
		0x70, 0x11, 0xe1, 0x22, 0xc2, 0x45, 0x84, 0x8b, 0x08, 0x17, 0x11, 0x2e, 0x4e, 0x14, 0x2e, 0x21,
		0x5c, 0x42, 0xb8, 0x84, 0x70, 0x09, 0xe1, 0x12, 0xc2, 0x25, 0x84, 0x4b, 0x08, 0x97, 0x10, 0x2e,
		0x21, 0x5c, 0x9a, 0x28, 0x5c, 0x46, 0xb8, 0x8c, 0x70, 0x19, 0xe1, 0x32, 0xc2, 0x65, 0x84, 0xcb,
		0x08, 0x97, 0x11, 0x2e, 0x23, 0x5c, 0x46, 0xb8, 0x3c, 0x51, 0xb8, 0x82, 0x70, 0x05, 0xe1, 0x0a,
		0xc2, 0x15, 0x84, 0x2b, 0x08, 0x57, 0x10, 0xae, 0x20, 0x5c, 0x41, 0xb8, 0x82, 0x70, 0x65, 0x62,
		0xa1, 0x56, 0x43, 0xb8, 0x86, 0x70, 0x0d, 0xe1, 0x1a, 0xc2, 0x35, 0x84, 0x6b, 0x08, 0xd7, 0x10,
		0xae, 0x21, 0x5c, 0x43, 0xb8, 0xc6, 0x13, 0x55, 0xe3, 0x89, 0xaa, 0x25, 0x14, 0x6a, 0x60, 0xea,
		0x60, 0xea, 0x60, 0xea, 0x60, 0xea, 0x60, 0xea, 0x60, 0xea, 0x60, 0xea, 0x60, 0xea, 0x60, 0xea,
		0x60, 0xea, 0x60, 0xea, 0x60, 0xea, 0x13, 0xeb, 0xc1, 0x26, 0x98, 0x26, 0x98, 0x26, 0x98, 0x26,
		0x98, 0x26, 0x98, 0x26, 0x98, 0x26, 0x98, 0x26, 0x98, 0x26, 0x98, 0x26, 0x98, 0x26, 0x98, 0x26,
		0xfb, 0x43, 0x93, 0xfd, 0xa1, 0xf9, 0xb1, 0x7a, 0xb0, 0x05, 0xac, 0x05, 0xac, 0x05, 0xac, 0x05,
		0xac, 0x05, 0xac, 0x05, 0xac, 0x05, 0xac, 0x05, 0xac, 0x05, 0xac, 0x05, 0xac, 0x05, 0xac, 0x05,
		0xac, 0x05, 0xac, 0xf5, 0x31, 0x58, 0x07, 0x58, 0x07, 0x58, 0x07, 0x58, 0x07, 0x58, 0x07, 0x58,
		0x07, 0x58, 0x07, 0x58, 0x07, 0x58, 0x07, 0x58, 0x07, 0x58, 0x07, 0x58, 0x07, 0x58, 0x07, 0x58,
		0x87, 0x9d, 0xaf, 0x33, 0xb2, 0x77, 0x9f, 0x9e, 0x2e, 0x45, 0x4c, 0x71, 0x0b, 0x9f, 0x2e, 0xde,
		0xce, 0xd6, 0xcc, 0xb8, 0x6b, 0xc6, 0x3d, 0x33, 0x9e, 0x33, 0xe3, 0x79, 0x33, 0x5e, 0x30, 0xe3,
		0x45, 0x33, 0xee, 0x9b, 0xf1, 0x92, 0x19, 0x2f, 0x9b, 0x31, 0x67, 0xc6, 0x2b, 0x66, 0xbc, 0x6a,
		0xc6, 0x03, 0xc6, 0x45, 0x7d, 0xc0, 0xe3, 0x5b, 0x9e, 0xcd, 0x5b, 0x6c, 0x46, 0x3f, 0x98, 0x58,
		0x81, 0x0c, 0xe5, 0x9a, 0x33, 0x03, 0x97, 0xf1, 0x3a, 0x97, 0x47, 0x2c, 0xcc, 0x75, 0x6b, 0x61,
		0x6e, 0x58, 0x0b, 0x53, 0xb0, 0x16, 0xa6, 0x68, 0x2d, 0x4c, 0xc9, 0x5a, 0x98, 0xb2, 0xeb, 0x8b,
		0xee, 0xd0, 0x17, 0xdd, 0xa1, 0x2f, 0xba, 0x43, 0x5f, 0x74, 0xc7, 0x81, 0x4e, 0x8f, 0xb4, 0x08,
		0x9e, 0xd7, 0x47, 0xee, 0x83, 0x5e, 0x7c, 0xeb, 0x33, 0x9e, 0x5d, 0x01, 0x6e, 0x1a, 0xc9, 0x43,
		0x12, 0x44, 0x3c, 0xa9, 0xbb, 0xf8, 0x06, 0x5c, 0x92, 0x7a, 0xfb, 0x33, 0xa0, 0x3a, 0xdb, 0x9c,
		0xd4, 0xfe, 0x1c, 0x52, 0x1a, 0xc5, 0x93, 0x26, 0xa6, 0x68, 0xc8, 0xc6, 0x1b, 0x4f, 0xda, 0xd8,
		0xa2, 0x21, 0xfb, 0x61, 0x3c, 0x09, 0x75, 0x6e, 0x0d, 0x28, 0xc0, 0x36, 0x67, 0x6f, 0x82, 0x0e,
		0x29, 0x22, 0xe2, 0x09, 0x2e, 0x3e, 0xa0, 0xf6, 0xda, 0x9c, 0x77, 0x1b, 0x74, 0x88, 0xab, 0x8f,
		0x27, 0x03, 0xf7, 0x25, 0xb4, 0x01, 0xa5, 0xd8, 0x08, 0x1d, 0x9a, 0x8e, 0x77, 0x40, 0xdd, 0xb5,
		0xb9, 0x88, 0x4e, 0x68, 0x40, 0xf9, 0xb5, 0xb9, 0x88, 0x4e, 0x68, 0x40, 0x15, 0xb6, 0xb9, 0x88,
		0x4e, 0x68, 0x40, 0x31, 0xb6, 0xf9, 0x4c, 0x76, 0x42, 0xc5, 0xec, 0x04, 0xb6, 0x81, 0x38, 0x94,
		0x07, 0x7a, 0x7e, 0x56, 0xcf, 0x1f, 0x8e, 0x17, 0x40, 0xf1, 0xe4, 0x75, 0x73, 0xd0, 0x93, 0x5b,
		0x6c, 0x20, 0xb7, 0x5c, 0x1d, 0x9b, 0x4e, 0x8d, 0xa6, 0xeb, 0xd8, 0x9c, 0xeb, 0xdd, 0x8c, 0x95,
		0x5a, 0x4f, 0xd9, 0xbb, 0x39, 0x74, 0x3a, 0x63, 0x35, 0xd6, 0x53, 0xea, 0x1c, 0xb9, 0x1e, 0xd0,
		0xd4, 0xc5, 0xd5, 0x58, 0x0f, 0x28, 0xef, 0x14, 0xa7, 0xae, 0xaa, 0xc6, 0x14, 0x0b, 0x4e, 0x71,
		0xea, 0x72, 0x6a, 0x4c, 0xb1, 0xe8, 0x14, 0xa7, 0xae, 0xa3, 0xc6, 0x14, 0x4b, 0x4e, 0x71, 0xea,
		0x02, 0x6a, 0x4c, 0xb1, 0xec, 0x14, 0xa7, 0xae, 0x9c, 0xc6, 0x14, 0x2b, 0xae, 0x3b, 0x37, 0xa7,
		0x92, 0x69, 0xac, 0x3b, 0x57, 0x73, 0xfa, 0x73, 0xaa, 0x95, 0xc6, 0xf4, 0xeb, 0xee, 0x01, 0x4d,
		0xa5, 0x48, 0x32, 0xdd, 0xbf, 0xa6, 0xa3, 0xa4, 0x52, 0x1d, 0x19, 0x4a, 0xcb, 0x51, 0xd2, 0x29,
		0x8b, 0x34, 0xa4, 0xe3, 0x20, 0x0b, 0xa8, 0x4d, 0x34, 0x51, 0xd7, 0x28, 0xce, 0x02, 0x66, 0x3f,
		0x31, 0x0b, 0xb8, 0x0b, 0x7c, 0xd7, 0x5a, 0xc0, 0x37, 0xac, 0x05, 0x7c, 0x73, 0xd4, 0x02, 0x66,
		0x17, 0x68, 0x01, 0xb3, 0x93, 0x2d, 0xe0, 0x5b, 0xd6, 0x02, 0xbe, 0x6d, 0x2d, 0xe0, 0x3b, 0xd6,
		0x02, 0x7e, 0xdf, 0x5a, 0xc0, 0x1f, 0x58, 0x0b, 0xf8, 0xc3, 0x51, 0x0b, 0x98, 0x9d, 0x87, 0x05,
		0xfc, 0x91, 0xb5, 0x80, 0x3f, 0x1e, 0xb5, 0x80, 0xd9, 0xf9, 0x5b, 0xc0, 0x9f, 0x58, 0x0b, 0xf8,
		0xd3, 0x51, 0x0b, 0x98, 0x4d, 0xdb, 0x02, 0x3e, 0x1a, 0xb5, 0x80, 0xd9, 0x05, 0x5a, 0xc0, 0xec,
		0x02, 0x2d, 0x60, 0x76, 0x81, 0x16, 0x30, 0xfb, 0x4c, 0x5a, 0x40, 0x4a, 0xd1, 0x3d, 0xc2, 0xdb,
		0x1b, 0xe9, 0x0b, 0xbd, 0x3b, 0x4b, 0x0b, 0xff, 0xf8, 0xc6, 0xe3, 0x1b, 0x17, 0x7c, 0xa3, 0x6d,
		0xf2, 0xfc, 0xf6, 0x09, 0x4d, 0x9e, 0xdf, 0xb9, 0x3f, 0xf8, 0xab, 0x94, 0x64, 0x57, 0x27, 0xe7,
		0xc2, 0x55, 0x9a, 0x37, 0xab, 0x34, 0x6f, 0x56, 0x69, 0xde, 0xac, 0x8e, 0x44, 0xf8, 0x7b, 0x1b,
		0xe1, 0x1f, 0x6c, 0x84, 0x7f, 0xb2, 0x11, 0xbe, 0x67, 0x23, 0x7c, 0xdf, 0x46, 0xf8, 0x81, 0xf9,
		0xbd, 0x86, 0x27, 0x27, 0x5d, 0x01, 0xf5, 0x22, 0xc8, 0x17, 0x47, 0x1e, 0xba, 0x0f, 0x67, 0x69,
		0x3a, 0x1d, 0xdf, 0x78, 0x7c, 0xe3, 0x82, 0x6f, 0x34, 0x0f, 0x9d, 0xa7, 0x26, 0x3f, 0x74, 0x9e,
		0xe7, 0x92, 0x60, 0x8e, 0x87, 0x2e, 0xe5, 0x57, 0xba, 0x78, 0xc7, 0x43, 0x58, 0x89, 0x2f, 0x68,
		0xa7, 0x68, 0x1f, 0x22, 0x7a, 0x84, 0x68, 0xe2, 0x6b, 0xda, 0xa9, 0x7b, 0x92, 0x79, 0x44, 0x67,
		0x7c, 0x59, 0x3b, 0xd6, 0x93, 0x44, 0xba, 0x80, 0xf4, 0x8c, 0xaf, 0x6b, 0x13, 0xa5, 0x8b, 0x48,
		0xcf, 0xf8, 0xc2, 0x36, 0x51, 0xba, 0x84, 0xf4, 0x8c, 0xaf, 0x6c, 0x13, 0xa5, 0xcb, 0x48, 0xcf,
		0xf8, 0xd2, 0x36, 0x51, 0xba, 0x82, 0xf4, 0x8c, 0xaf, 0x6d, 0x13, 0x5b, 0xcb, 0x35, 0xa4, 0xe7,
		0xfa, 0xe2, 0x76, 0xac, 0xb5, 0x0c, 0xa8, 0x0e, 0x68, 0xae, 0xaf, 0x6e, 0x13, 0x7b, 0xd8, 0x4d,
		0x40, 0x29, 0xbe, 0xbc, 0x1d, 0xed, 0x61, 0xb7, 0xa0, 0xa5, 0xf8, 0xf6, 0x76, 0x94, 0xd6, 0x81,
		0x96, 0xf2, 0xeb, 0x5b, 0x58, 0xd1, 0xab, 0xc0, 0x16, 0xfa, 0x32, 0xd5, 0x98, 0x15, 0xcf, 0xd7,
		0x61, 0x64, 0xe8, 0x2c, 0x04, 0x9f, 0x48, 0xd3, 0x22, 0x86, 0xdf, 0xe0, 0xc0, 0x0d, 0xd3, 0xb4,
		0xf0, 0x96, 0x4c, 0xd3, 0xc2, 0x0b, 0x5c, 0x7c, 0x03, 0x65, 0x7b, 0x17, 0xe9, 0x36, 0x2d, 0x32,
		0xf4, 0x27, 0x82, 0x89, 0x4d, 0x0b, 0x6f, 0xd9, 0xe4, 0x08, 0x6f, 0xc5, 0xe4, 0x08, 0x6f, 0xd5,
		0xe4, 0x08, 0xef, 0xa4, 0xc9, 0x11, 0x5e, 0xd6, 0xe4, 0x08, 0x6f, 0x4d, 0x6f, 0x2e, 0x19, 0xba,
		0x0e, 0xc1, 0x1c, 0x9a, 0x16, 0xde, 0xba, 0xd9, 0x5c, 0xbc, 0x0d, 0xfd, 0xcc, 0x67, 0xe8, 0x33,
		0x04, 0x73, 0x6f, 0x5a, 0x78, 0x9b, 0xe6, 0x99, 0xf7, 0xb6, 0xdc, 0x97, 0xd0, 0x56, 0xb6, 0x77,
		0x91, 0x5a, 0xd3, 0xc2, 0xdb, 0x76, 0xb4, 0x50, 0xd9, 0xde, 0x45, 0xba, 0x4d, 0x8b, 0x0c, 0x9d,
		0x85, 0x60, 0x11, 0x4d, 0x8b, 0x0c, 0x9d, 0x85, 0x60, 0x11, 0x4d, 0x8b, 0x0c, 0x9d, 0x85, 0xe0,
		0x99, 0x6c, 0x5a, 0x38, 0x3b, 0x99, 0x5d, 0xa0, 0x9d, 0xcc, 0x4e, 0xb6, 0x93, 0x4f, 0xdd, 0x8a,
		0x74, 0x76, 0x32, 0x3b, 0xd9, 0x4e, 0x4e, 0xd9, 0xdf, 0xcc, 0x2b, 0xfb, 0x13, 0xc0, 0xd9, 0xec,
		0xe4, 0x63, 0xfd, 0x4d, 0x67, 0x27, 0xb3, 0xf3, 0xb0, 0x93, 0x09, 0xd2, 0x45, 0x65, 0x7f, 0x06,
		0x38, 0x9b, 0x9d, 0x4c, 0x90, 0x2e, 0x29, 0xfb, 0x43, 0xc0, 0xd9, 0xec, 0x64, 0x82, 0x74, 0x59,
		0xd9, 0x9f, 0x02, 0xce, 0x66, 0x27, 0x13, 0xa4, 0x2b, 0xca, 0xfe, 0x18, 0x70, 0x36, 0x3b, 0x99,
		0xd0, 0xa6, 0xae, 0x29, 0xfb, 0x73, 0xc0, 0x79, 0xda, 0xc9, 0xc7, 0xda, 0xd4, 0xce, 0x4e, 0x66,
		0xe7, 0x6f, 0x27, 0x13, 0xfa, 0xe1, 0x4d, 0x65, 0x7f, 0x12, 0x98, 0x9e, 0x9d, 0x7c, 0x34, 0x6a,
		0x27, 0xb3, 0x69, 0xdb, 0xc9, 0x47, 0xa3, 0x76, 0x32, 0xbb, 0x48, 0x3b, 0x99, 0x5d, 0xb4, 0x9d,
		0xf4, 0xed, 0x66, 0xbe, 0xc2, 0xb0, 0xce, 0xb0, 0x65, 0xf6, 0x5f, 0xeb, 0x36, 0x77, 0x67, 0x69,
		0xbd, 0x1d, 0xdf, 0x78, 0x7c, 0xe3, 0x82, 0x6f, 0xb4, 0xad, 0xb1, 0xdc, 0x13, 0x5a, 0x63, 0x57,
		0x9c, 0xd5, 0xca, 0x29, 0xfb, 0x14, 0xa4, 0xeb, 0x65, 0x32, 0xd8, 0x96, 0x60, 0xa2, 0x97, 0x79,
		0xfa, 0x0a, 0x05, 0xd1, 0x23, 0x65, 0x2d, 0x4d, 0x82, 0x97, 0x99, 0xb6, 0xec, 0xc9, 0x2b, 0x6b,
		0x69, 0x66, 0xf2, 0x32, 0x8f, 0x97, 0x3d, 0x48, 0x17, 0x94, 0xb5, 0x34, 0x33, 0x79, 0x99, 0x24,
		0xe9, 0xa2, 0xb2, 0x96, 0x66, 0x26, 0x2f, 0x93, 0x24, 0x5d, 0x52, 0xd6, 0xd2, 0xcc, 0xe4, 0x65,
		0x92, 0xa4, 0xcb, 0xca, 0x5a, 0x9a, 0x99, 0xbc, 0x4c, 0x92, 0x74, 0x45, 0x59, 0x4b, 0x33, 0x93,
		0x97, 0x49, 0xaa, 0x5e, 0x6b, 0xca, 0x5a, 0x9a, 0x39, 0x7a, 0x99, 0xc7, 0xab, 0x57, 0x40, 0x75,
		0x65, 0x2d, 0xcd, 0x1c, 0xbd, 0x4c, 0x52, 0x99, 0xdc, 0x54, 0xd6, 0xd2, 0xa4, 0xe6, 0x65, 0x46,
		0xca, 0xe4, 0x96, 0xb2, 0x96, 0x26, 0x35, 0x2f, 0x33, 0x42, 0x33, 0xbd, 0x9b, 0xd4, 0xbd, 0x4c,
		0xc6, 0x78, 0x8c, 0x60, 0xc1, 0x5e, 0xe6, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe2, 0x12, 0xe8,
		0x9c, 0xe6, 0x45, 0x00, 0x00,
	}
	buf, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&tab); err != nil {
		panic(err)
	}

	for i, row := range tab {
		actionTab[i].canRecover = row.CanRecover
		for _, a := range row.Actions {
			switch a.Action {
			case 0:
				actionTab[i].actions[a.Index] = accept(true)
			case 1:
				actionTab[i].actions[a.Index] = reduce(a.Amount)
			case 2:
				actionTab[i].actions[a.Index] = shift(a.Amount)
			}
		}
	}
}
