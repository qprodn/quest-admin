// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package slices

import (
	"github.com/samber/lo"
)

func ToMap[T any, K comparable, V any](s []T, t func(e T) (K, V)) map[K]V {
	return lo.SliceToMap(s, t)
}

func Transform[T any, R any](s []T, iteratee func(e T, idx int) R) []R {
	return lo.Map(s, iteratee)
}

func Uniq[T comparable, Slice ~[]T](s Slice) Slice {
	return lo.Uniq(s)
}

func Map[T any, R any](s []T, f func(item T, index int) R) []R {
	return lo.Map(s, f)
}

func Contains[T comparable](s []T, v T) bool {
	return lo.Contains(s, v)
}

func Filter[T any](s []T, predicate func(item T, index int) bool) []T {
	return lo.Filter(s, predicate)
}

func FilterMap[T, R any](collection []T, callback func(item T, index int) (R, bool)) []R {
	return lo.FilterMap(collection, callback)
}

func ForEach[T any](collection []T, iteratee func(item T, index int)) {
	lo.ForEach(collection, iteratee)
}
