// Copyright 2016 Aleksey Blinov.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build aspectrt

package aspect_test

import (
	"runtime"
	"testing"
)

// Aspect assertions in runtime only work for aspects
// expressed as interfaces. A compiler tweak is needed
// to support direct interface aspect assertions. This
// however would result in slight slowdown of the code
// requiring such assertions. We are not doing it for now.
//
// func TestAssertAspectDynamic2(t *testing.T) {
// 	_ = p2.As().(*ChessPlayer).Rank
// 	_ = p2.As().(*Employee).Rank
// 	_ = p2.As().(*ChessPlayer).Rank
// 	_ = p2.As().(*Employee).Rank
// 	_ = i2.As().(*Person).Name
// 	_ = i2.As().(*ChessPlayer).Rank
// 	_ = i2.As().(*Employee).Rank
// }

func TestAssertAspectDynamic(t *testing.T) {
	_ = rentable1.As().(Apt).NBeds()
	_ = apt1.As().(Rentable).Rate()
}

