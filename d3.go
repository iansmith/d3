package d3

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

//d3root is the primary d3 object (function) in javascript land.
var d3root = js.Global.Get("d3")

//Selector is a special kind of string used for selecting elements
//from the dom. Roughly, it's a CSS selector.
type Selector string

//TagName is a html tag name, like "div"
type TagName string

//PropertyName is the name of a CSS property like "width"
type PropertyName string

//=================================================================

//Selection is the d3 concept of zero or more selected elements.
type Selection interface {
	Append(TagName) Selection
	SelectAll(Selector) Selection
	Data(js.Object) Selection
	Enter() Selection
	Style(PropertyName, func(js.Object) string) Selection
	StyleS(PropertyName, string) Selection
	Text(func(js.Object) string) Selection
	TextS(string) Selection
	Attr(PropertyName, int64) Selection
	AttrF(PropertyName, float64) Selection
	AttrS(PropertyName, string) Selection
	AttrFunc(PropertyName, func(js.Object) int64) Selection
	AttrFuncF(PropertyName, func(js.Object) float64) Selection
	AttrFunc2S(PropertyName, func(js.Object, int64) string) Selection
	AttrFuncS(PropertyName, func(js.Object) string) Selection
	Call(Axis) Selection
}

//selectionImpl is the implementation of Selection.
type selectionImpl struct {
	obj js.Object
}

//SelectAll finds all DOM elements that match selector in the current
//selection.
func (self *selectionImpl) SelectAll(n Selector) Selection {
	return &selectionImpl{
		self.obj.Call("selectAll", string(n)),
	}
}

//Data provides a set of data for a data join to work with.  The argument
//must be a JS array.
func (self *selectionImpl) Data(arr js.Object) Selection {
	return &selectionImpl{
		self.obj.Call("data", arr),
	}
}

//Enter is the case of adding new elements to a data join.
func (self *selectionImpl) Enter() Selection {
	return &selectionImpl{
		self.obj.Call("enter"),
	}
}

//Append adds elements to the current selection.  The parameter is
//a tag name to be added.
func (self *selectionImpl) Append(n TagName) Selection {
	return &selectionImpl{
		self.obj.Call("append", string(n)),
	}
}

//Style modifies the CSS attribute prop of the selection.  The function
//is passed each element of the data set to use in computing the value.
func (self *selectionImpl) Style(prop PropertyName, f func(js.Object) string) Selection {
	console.Log("calling style", self.obj, prop, f)
	return &selectionImpl{
		self.obj.Call("style", string(prop), f),
	}
}

//StyleConst modifies the CSS attribute prop of the selection to be a
//constant value.
func (self *selectionImpl) StyleS(prop PropertyName, value string) Selection {
	return &selectionImpl{
		self.obj.Call("style", string(prop), value),
	}
}

//Text modifies the text portion of the selected elements to be the return
//values of the function.  The function is called for each value in the
//dataset.
func (self *selectionImpl) Text(f func(js.Object) string) Selection {
	return &selectionImpl{
		self.obj.Call("text", f),
	}
}

//TextS modifies the text portion of the selected elements to be
//a constant value.
func (self *selectionImpl) TextS(v string) Selection {
	return &selectionImpl{
		self.obj.Call("text", v),
	}
}

//Attr sets an attribute of the selection to a particular value.
func (self *selectionImpl) Attr(p PropertyName, v int64) Selection {
	return &selectionImpl{
		self.obj.Call("attr", string(p), v),
	}
}

//AttrF sets an attribute of the selection to a particular value.
func (self *selectionImpl) AttrF(p PropertyName, v float64) Selection {
	return &selectionImpl{
		self.obj.Call("attr", string(p), v),
	}
}

//Attr sets an attribute of the selection to a string value.
func (self *selectionImpl) AttrS(p PropertyName, v string) Selection {
	return &selectionImpl{
		self.obj.Call("attr", string(p), v),
	}
}

//AttrFunc2S sets an attribute to a function of two variables with the
//second being the already extracted integer.
func (self *selectionImpl) AttrFunc2S(p PropertyName, v func(js.Object, int64) string) Selection {
	return &selectionImpl{
		self.obj.Call("attr", string(p), v),
	}
}

//AttrFuncS sets an attribute to a function of the data object
func (self *selectionImpl) AttrFuncS(p PropertyName, v func(js.Object) string) Selection {
	return &selectionImpl{
		self.obj.Call("attr", string(p), v),
	}
}

//AttrFunc sets an attribute to a function of one variable that produces an int.
func (self *selectionImpl) AttrFunc(p PropertyName, v func(js.Object) int64) Selection {
	return &selectionImpl{
		self.obj.Call("attr", string(p), v),
	}
}

//AttrFuncF sets an attribute to a function of one variable that produces a float.
func (self *selectionImpl) AttrFuncF(p PropertyName, v func(js.Object) float64) Selection {
	return &selectionImpl{
		self.obj.Call("attr", string(p), v),
	}
}

//Call is a wrapper over the d3 selection call() method. No idea how it works.
func (self *selectionImpl) Call(a Axis) Selection {
	s := a.(*axisImpl)
	return &selectionImpl{
		self.obj.Call("call", s.obj),
	}
}

//=================================================================

/*
 * Utilities
 */

//Max is d3.max with a function that passes over each object in the array
//supplied as the first argument. If the second function is nil, we assume
//that the array contains (JS) integers.
func Max(v js.Object, fn ExtractorFunc) int64 {
	if fn != nil {
		return int64(d3root.Call("max", v, fn).Int())
	}
	return int64(d3root.Call("max", v).Int())
}

//MaxF is d3.max with a function that passes over each object in the array
//supplied as the first argument. If the second function is nil, we assume
//that the array contains (JS) floats.
func MaxF(v js.Object, fn ExtractorFuncF) float64 {
	if fn != nil {
		return d3root.Call("max", v, fn).Float()
	}
	return d3root.Call("max", v).Float()
}

//Select is d3.select() and creates a selection from the selector.
func Select(n Selector) Selection {
	return &selectionImpl{
		d3root.Call("select", string(n)),
	}
}

//ScaleLinear is d3.scale.linear and returns a LinearScale
func ScaleLinear() LinearScale {
	return &linearScaleImpl{
		d3root.Get("scale").Call("linear"),
	}
}

//ScaleLinear is d3.scale.ordinal and returns a Ordinal
func ScaleOrdinal() OrdinalScale {
	return &ordinalScaleImpl{
		d3root.Get("scale").Call("ordinal"),
	}
}

//NewAxis creates a new axis object
func NewAxis() Axis {
	return &axisImpl{
		d3root.Get("svg").Call("axis"),
	}
}

//FilterFunc converts "raw" objects to their formatted counterparts.
//Raw version has only string fields, but formatted version should have
//the parsed values. If this func returns nil, that item is ignored.
type FilterFunc func(js.Object) js.Object

//ExtractorFunc is a fun that can pull the int value from an object
type ExtractorFunc func(js.Object) int64

//ExtractorFuncF is a fun that can pull the float value from an object
type ExtractorFuncF func(js.Object) float64

//ExtractorFuncO is a fun that can pull the named (usually ordinal) value from the object
type ExtractorFuncO func(js.Object) js.Object

//TSV loads a tab separated value called filename from the server.
//Each loaded element is passed through filter func and the final
//result is handed to callback.
func TSV(filename string, filter FilterFunc, callback func(js.Object, js.Object)) {
	d3root.Call("tsv", filename, filter, callback)
}

//=================================================================

//LinearScale is a wrapper around the d3 concept of the same name.
type LinearScale interface {
	Domain([]int64) LinearScale
	DomainF([]float64) LinearScale
	Range([]int64) LinearScale
	Linear(js.Object, ExtractorFunc) int64
	LinearF(js.Object, ExtractorFuncF) float64
	Invert(js.Object, ExtractorFunc) int64
	Func(ExtractorFunc) func(js.Object) int64
	FuncF(ExtractorFuncF) func(js.Object) float64
}

//linearScaleImpl is the implementation of LinearScale.
type linearScaleImpl struct {
	obj js.Object
}

//Domain sets the domain of the linear domain.
func (self *linearScaleImpl) Domain(d []int64) LinearScale {
	return &linearScaleImpl{
		self.obj.Call("domain", d),
	}
}

//Domain sets the domain of the linear domain.
func (self *linearScaleImpl) DomainF(d []float64) LinearScale {
	in := js.Global.Get("Array").New()
	for i := 0; i < len(d); i++ {
		in.SetIndex(i, d[i])
	}

	return &linearScaleImpl{
		self.obj.Call("domain", in),
	}
}

//Range sets the range of the linear domain.
func (self *linearScaleImpl) Range(d []int64) LinearScale {
	return &linearScaleImpl{
		self.obj.Call("range", d),
	}
}

//Linear calls the scale to interpolate a value into its range.
func (self *linearScaleImpl) Linear(obj js.Object, fn ExtractorFunc) int64 {
	if fn != nil {
		return int64(self.obj.Invoke(fn(obj)).Int())
	}
	return int64(self.obj.Invoke(obj.Int()).Int())
}

//LinearF calls the scale to interpolate a value into its range and returns
//the results as a float. If the extractor function is specified, it should
//pull the floating point input value out of the provided object.
func (self *linearScaleImpl) LinearF(obj js.Object, fn ExtractorFuncF) float64 {
	if fn != nil {
		return self.obj.Invoke(fn(obj)).Float()
	}
	return self.obj.Invoke(obj.Float()).Float()
}

//Invert calls the scale to interpolate a range value into its domain.
func (self *linearScaleImpl) Invert(obj js.Object, fn ExtractorFunc) int64 {
	if fn != nil {
		return int64(self.obj.Call("invert", fn(obj)).Int())
	}
	return int64(self.obj.Call("invert", obj.Int()).Int())
}

//Func returns a function wrapper around this scale so it can be used
//in the Go side as a function that can extract the values from the
//objects as integers.
func (self *linearScaleImpl) Func(fn ExtractorFunc) func(js.Object) int64 {
	if fn != nil {
		return func(obj js.Object) int64 {
			return int64(self.obj.Invoke(fn(obj)).Int())
		}
	}
	return func(obj js.Object) int64 {
		return int64(self.obj.Invoke(obj.Int()).Int())
	}
}

//FuncF returns a function wrapper around this scale so it can be used
//in the Go side as a function that can extract the values from the
//objects as floats.
func (self *linearScaleImpl) FuncF(fn ExtractorFuncF) func(js.Object) float64 {
	if fn != nil {
		return func(obj js.Object) float64 {
			return self.obj.Invoke(fn(obj)).Float()
		}
	}
	return func(obj js.Object) float64 {
		return self.obj.Invoke(obj.Float()).Float()
	}
}

//=================================================================

//OrdinalScale wraps d3.scale.ordinal
type OrdinalScale interface {
	Domain(js.Object) OrdinalScale
	RangeBands([]int64) OrdinalScale
	RangeBand() int64
	RangeBandF() float64
	RangeBands3([]int64, float64) OrdinalScale
	Ordinal(obj js.Object, fn ExtractorFuncO) int64
}

//ordinalScaleImpl is the implementation of LinearScale.
type ordinalScaleImpl struct {
	obj js.Object
}

//Domair should be an array of items.
func (self *ordinalScaleImpl) Domain(obj js.Object) OrdinalScale {
	return &ordinalScaleImpl{
		self.obj.Call("domain", obj),
	}
}

func (self *ordinalScaleImpl) RangeBands(b []int64) OrdinalScale {
	return &ordinalScaleImpl{
		self.obj.Call("rangeBands", b),
	}
}

func (self *ordinalScaleImpl) RangeBand() int64 {
	return int64(self.obj.Call("rangeBand").Int())
}

func (self *ordinalScaleImpl) RangeBandF() float64 {
	return self.obj.Call("rangeBand").Float()
}

func (self *ordinalScaleImpl) RangeBands3(b []int64, f float64) OrdinalScale {
	return &ordinalScaleImpl{
		self.obj.Call("rangeBands", b, f),
	}
}

//Ordinal calls the scale to interpolate a value into its range.  If the
//second function is nil, then we assume the ojbect is already in the ordinal
//domain.
func (self *ordinalScaleImpl) Ordinal(obj js.Object, fn ExtractorFuncO) int64 {
	if fn != nil {
		return int64(self.obj.Invoke(fn(obj)).Int())
	}
	return int64(self.obj.Invoke(obj).Int())
}

//=================================================================

type Edge int

const (
	BOTTOM = iota
	TOP
	RIGHT
	LEFT
)

func (self Edge) String() string {
	switch self {
	case BOTTOM:
		return "bottom"
	case LEFT:
		return "left"
	case RIGHT:
		return "right"
	case TOP:
		return "top"
	}
	panic("bad edge value?!?")
}

//Axis is a VERY thin wrapper over d3.svg.axis
type Axis interface {
	ScaleO(OrdinalScale) Axis
	Scale(LinearScale) Axis
	Orient(e Edge) Axis
	Ticks(int64, string) Axis
}

type axisImpl struct {
	obj js.Object
}

//ScaleO creates an axis, given an already created ordinal scale.
func (self *axisImpl) ScaleO(scale OrdinalScale) Axis {
	s := scale.(*ordinalScaleImpl)
	return &axisImpl{
		self.obj.Call("scale", s.obj),
	}
}

//Scale creates an axis, given an already created linear scale.
func (self *axisImpl) Scale(scale LinearScale) Axis {
	s := scale.(*linearScaleImpl)
	return &axisImpl{
		self.obj.Call("scale", s.obj),
	}
}

//Orient binds the axis to one of the four edges.
func (self *axisImpl) Orient(e Edge) Axis {
	return &axisImpl{
		self.obj.Call("orient", e.String()),
	}
}

//Ticks changes the way the ticks on the axis look.  You can optionally
//pass the 2nd parameter for formatting; use "" for no formatting.
func (self *axisImpl) Ticks(i int64, format string) Axis {
	if format == "" {
		return &axisImpl{
			self.obj.Call("ticks", i),
		}
	}
	return &axisImpl{
		self.obj.Call("ticks", i, format),
	}
}
