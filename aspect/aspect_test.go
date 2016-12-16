// Copyright 2016 Aleksey Blinov.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aspect_test

import (
	"github.com/baobabus/gasp/aspect"
	"reflect"
	"testing"
)

type Person struct {
	aspect.Aspect
	Name string
}

type ChessPlayer struct {
	aspect.Aspect
	Rank uint8
}

type Employee struct {
	aspect.Aspect
	Rank string
}

var (
	p2 = &Person{Name: "Bob"}
	c2 = &ChessPlayer{Rank: 1}
	e2 = &Employee{Rank: "director"}
	i2 = aspect.NewEntity(p2, c2, e2)
)

func TestAssertAspectStatic2(t *testing.T) {
	_ = p2.Name
	_ = c2.Rank
	_ = e2.Rank
	_ = p2.As(reflect.TypeOf((*ChessPlayer)(nil)).Elem()).(*ChessPlayer).Rank
	_ = p2.As(reflect.TypeOf((*Employee)(nil)).Elem()).(*Employee).Rank
	_ = p2.As((*ChessPlayer)(nil)).(*ChessPlayer).Rank
	_ = p2.As((*Employee)(nil)).(*Employee).Rank
	_ = i2.As((*Person)(nil)).(*Person).Name
	_ = i2.As((*ChessPlayer)(nil)).(*ChessPlayer).Rank
	_ = i2.As((*Employee)(nil)).(*Employee).Rank
}

type Apt interface {
	aspect.Assertable
	NBeds() uint
} 

type apt struct {
	aspect.Aspect
	nbeds uint
}

func (this *apt) NBeds() uint {
	return this.nbeds
}

func (this *apt) ExposedAspects() []interface{} {
	return []interface{}{(*Apt)(nil)}
}

type Rentable interface {
	aspect.Assertable
	Rate() float32
}

type rentable struct {
	aspect.Aspect
	rate float32
}

func (this *rentable) Rate() float32 {
	return this.rate
}

func (this *rentable) ExposedAspects() []interface{} {
	return []interface{}{(*Rentable)(nil)}
}

var (
	apt1 = &apt{nbeds: 2}
	rentable1 = &rentable{rate: 1000.0}
	r1 = aspect.NewEntity(apt1, rentable1)
)

func TestAssertAspectStatic(t *testing.T) {
	_ = apt1.NBeds()
	_ = rentable1.Rate()
	_ = rentable1.As(reflect.TypeOf((*Apt)(nil)).Elem()).(Apt).NBeds()
	_ = rentable1.As((*Apt)(nil)).(Apt).NBeds()
	_ = rentable1.As((*Apt)(nil)).(Apt).NBeds()
}

func TestAssertAspectPanicMult(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got none")
		}
	}()
	_ = rentable1.As(reflect.TypeOf((*Apt)(nil)).Elem(), nil).(Apt).NBeds()
}
