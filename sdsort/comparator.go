package sdsort

import (
	dsutils "github.com/emirpasic/gods/utils"
)

type Comparator = dsutils.Comparator

var (
	// byte
	ByteComparator Comparator = dsutils.ByteComparator

	// int
	IntComparator   Comparator = dsutils.IntComparator
	Int8Comparator  Comparator = dsutils.Int8Comparator
	Int16Comparator Comparator = dsutils.Int16Comparator
	Int32Comparator Comparator = dsutils.Int32Comparator
	Int64Comparator Comparator = dsutils.Int64Comparator

	// uint
	UIntComparator   Comparator = dsutils.UIntComparator
	UInt8Comparator  Comparator = dsutils.UInt8Comparator
	UInt16Comparator Comparator = dsutils.UInt16Comparator
	UInt32Comparator Comparator = dsutils.UInt32Comparator
	UInt64Comparator Comparator = dsutils.UInt64Comparator

	// float
	Float32Comparator Comparator = dsutils.Float32Comparator
	Float64Comparator Comparator = dsutils.Float64Comparator

	// string
	StringComparator Comparator = dsutils.StringComparator

	// time
	TimeComparator Comparator = dsutils.TimeComparator
)
