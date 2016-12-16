// Copyright 2016 Aleksey Blinov.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aspect

import (
	"errors"
	"reflect"
)

// Assertable must be implemented by any type for its instances
// to be acceptable as aspects.
// This is normally accomplished by including Aspect as an anonymous
// structure member of the aspect type.
type Assertable interface {
	As(ifaces ...interface{}) interface{}
}

// Aspect types wishing to have their entity reference
// automatically updated when attached to an entity
// should implement EntityAssignable interface.
type EntityAssignable interface {
	NoteAssignedTo(e Entity)
}

// If a type implements aspect interfaces other than its own
// reference type, it can implement AspectExposing to have their
// instances automatically mapped as entity aspects of those interfaces.
type AspectExposing interface {
	ExposedAspects() []interface{}
}

// Types representing aspect-enabled entities must implement Aspectual
// interface. Normally Entity type is sufficient for this purpose.
type Aspectual interface {
	As(ifaces ...interface{}) interface{}
	AddAspect(iface interface{}, imp Assertable)
}

// Default type for representing entity instances. Custom types
// can also be used for entity instances as long as they implement
// Aspectual interface. Additionally, for runtime-based aspect
// interface assertions to work, the type must be convertible to
// runtime.AspectMap.
type Entity map[uintptr]interface{}

// Aspect type can be included as an anonymous field in structures
// representing an aspect.
type Aspect struct {
	Entity
}

// NewEntity creates and returns a new aspect-enabled entity instance.
// Supplied aspects are added to the new instance as implementing
// their respected interfaces. Aspects are notified about the assignment
// as described in EntityAssignable interface.
// Additionally, any supplied aspects that implement AspectExposing interface
// get attached as representing the additional returned interfaces.
func NewEntity(aspects ...Assertable) Entity {
	var res Entity = make(map[uintptr]interface{})
	for _, v := range aspects {
		res.AddAspect(v, v)
		if ae, ok := v.(AspectExposing); ok {
			l := ae.ExposedAspects()
			for _, a := range l {
				res.AddAspect(a, v)
			}
		}
	}
	return res
}

// Implementation of Assertable interface.
func (r Entity) As(ifaces ...interface{}) interface{} {
	return r.as(ifaces)
}

// Implementation of Aspectual interface.
func (r Entity) AddAspect(iface interface{}, imp Assertable) {
	// TODO Optimize using unsafe.
	k := reflect.ValueOf(reflect.TypeOf(iface).Elem()).Pointer()
	if v, ok := r[k]; ok {
		if v != imp {
			panic(errBound)
		}
	} else  {
		r[k] = imp;
		if a, ok := imp.(EntityAssignable); ok {
			a.NoteAssignedTo(r)
		}
	}
}

// Implementation of Assertable interface.
func (r *Aspect) As(ifaces ...interface{}) interface{} {
	return r.Entity.as(ifaces)
}

// Implementation of EntityAssignable interface.
func (a *Aspect) NoteAssignedTo(r Entity) {
	// if a.Entity != nil && a.Entity != r {
	// 	panic(errForeign)
	// }
	a.Entity = r
}

var errMult = errors.New("attempt to assert multiple aspects at once")
var errNonAspect = errors.New("attempt to assert non-aspect")
var errForeign = errors.New("attempt to attach foriegn aspect")
var errBound = errors.New("attempt to attach to already bound aspect")

func (r Entity) as(ifaces []interface{}) interface{} {
	switch len(ifaces) {
	default:
		panic(errMult)
	case 0:
		if caster != nil {
			return caster(r)
		}
		return r
	case 1:
		switch c := ifaces[0].(type) {
		default:
			// TODO Optimize using unsafe.
			t := reflect.ValueOf(reflect.TypeOf(c).Elem()).Pointer()
			return r[t]
		case reflect.Type:
			t := reflect.ValueOf(c).Pointer()
			return r[t]
		}
	}
}

var caster func(Entity) interface{}
