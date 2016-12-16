// Copyright 2016 Aleksey Blinov.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Functionality provided by this source is dependent on
// certain go runtime modifications.
//
// When implemented, shorthand assertions are possible:
// 
//   aspect2 := aspect1.As().(Aspect2)
//
// instead of:
//
//   aspect2 := aspect1.As((*Aspect2)(nil)).(Aspect2)
//
// or worse:
//
//   aspect2 := aspect1.As(reflect.TypeOf((*Aspect2)(nil)).Elem()).(Aspect2)
//
// Here Aspect1 and Aspect2 are interface types. Variable aspect1 is of type Aspect1.
// 
// Specific runtime changes:
//
//   1. Definition of type AspectMap:
//      type AspectMap map[uintptr]interface{}
//
//   2. Modification of assertE2I() and assertE2I2() functions.
//      Modified functions check if the subject type is AspectMap,
//      and if so, attempt to lookup an alternative candidate
//      (an aspect) in the map by the target interface.
//      If the lookup is successful, the subject is replaced
//      by the mapped aspect.
//
//      Runtime performance impact in non-aspectual cases is negligible.
//      Only a single additional function call is made, where a simple
//      pointer equality test is made. For additional optimization,
//      the function code can be easily inlined.
//      
//      func assertE2I(inter *interfacetype, e eface, r *iface) {
//      	t := e._type
//      	if t == nil {
//      		// explicit conversions require non-nil interface value.
//      		panic(&TypeAssertionError{"", "", inter.typ.string(), ""})
//      	}
// +>    	switchE2A(&e, &t, inter)
//      	r.tab = getitab(inter, t, false)
//      	r.data = e.data
//      }
//
//      func assertE2I2(inter *interfacetype, e eface, r *iface) bool {
//      	if testingAssertE2I2GC {
//      		GC()
//      	}
//      	t := e._type
//      	if t == nil {
//      		if r != nil {
//      			*r = iface{}
//      		}
//      		return false
//      	}
// +>     	switchE2A(&e, &t, inter)
//      	tab := getitab(inter, t, true)
//      	if tab == nil {
//      		if r != nil {
//      			*r = iface{}
//      		}
//      		return false
//      	}
//      	if r != nil {
//      		r.tab = tab
//      		r.data = e.data
//      	}
//      	return true
//      }
//
//      func switchE2A(e *eface, t **_type, inter *interfacetype) {
//      	if e._type == aspm {
//      		ex := *(*AspectMap)(e.data)
//      		if v, ok := ex[uintptr(unsafe.Pointer(&inter.typ))]; ok {
//      			*e = *efaceOf(&v)
//      			if t != nil {
//      				*t = e._type
//      			}
//      		}
//      	}
//      }
//
//      var aspm *_type = func() *_type {
//      	var m AspectMap = make(map[uintptr]interface{})
//      	var i interface{} = m
//      	e := efaceOf(&i)
//      	return e._type.typeOff(e._type.ptrToThis)
//      }()

// +build aspectrt

package aspect

import (
	"runtime"
)

// In order to benefit from runtime support for aspect-oriented
// type assertion, aspectual entity must be cast to runtime.AspectMap.
func init() {
	caster = func(r Entity) interface{} {
		var c runtime.AspectMap = runtime.AspectMap(r)
		var res interface{} = &c
		return res;
	}
}