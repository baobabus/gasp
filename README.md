# gasp

Aspect-oriented data modeling in go.

[![Build Status (Linux)](https://travis-ci.org/baobabus/gasp.svg?branch=master)](https://travis-ci.org/baobabus/gasp)

### Sample model
```go
package main

import (
	"fmt"
	"github.com/baobabus/gasp/aspect"
)

type Person interface {
	aspect.Assertable
	Name() string
}

type ChessPlayer interface {
	aspect.Assertable
	Rank() uint8
}

type Employee interface {
	aspect.Assertable
	Rank() string
}

// implements Person interface below
type person struct {
	aspect.Aspect
	name string
}

// implements ChessPlayer interface below
type chessPlayer struct {
	aspect.Aspect
	rank uint8
}

// implements Employee interface below
type employee struct {
	aspect.Aspect
	rank string
}

var (
	i1 = aspect.NewEntity(&person{name: "Bob"})
	i2 = aspect.NewEntity(&person{name: "Carl"}, &chessPlayer{rank: 5})
	i3 = aspect.NewEntity(&person{name: "Sophie"}, &employee{rank: "executive"})
	i4 = aspect.NewEntity(&person{name: "Daphne"}, &employee{rank: "director"}, &chessPlayer{rank: 7})
)

func main() {

	p1 := i1.As().(Person) // *person
	n1 := i1.As().(Person).Name() // "Bob"
	c1, ok := i1.As().(ChessPlayer) // nil, false

	c2, ok := i2.As().(ChessPlayer) // *chessPlayer, true
	cr2 := c2.Rank() // 5
	n2 := c2.As().(Person).Name() // "Carl"
	e2, ok := c2.As().(Employee) // nil, false

	e4 := i4.As().(Employee) // *employee
	n4 := e4.As().(Person).Name() // "Daphne"
	p4 := e4.As().(Person) // *person
	cr4 := p4.As().(ChessPlayer).Rank() // 7

	// Shut up the compiler
	foo := []interface{}{p1, n1, c1, c2, cr2, n2, e2, e4, n4, p4, cr4, ok}
	fmt.Printf("%v\n", foo)
}

func (a *person) ExposedAspects() []interface{} {
	return []interface{}{(*Person)(nil)}
}

func (a *chessPlayer) ExposedAspects() []interface{} {
	return []interface{}{(*ChessPlayer)(nil)}
}

func (a *employee) ExposedAspects() []interface{} {
	return []interface{}{(*Employee)(nil)}
}

func (a *person) Name() string {
	return a.name
}

func (a *chessPlayer) Rank() uint8 {
	return a.rank
}

func (a *employee) Rank() string {
	return a.rank
}
```
