package main

import (
	"d3"
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
	_ "reflect"
	"strconv"
)

var (
	pickChart2 = d3.Selector(".part2_chart")
	pickChart3 = d3.Selector(".part3_chart")
	pickDiv    = d3.Selector("div")
	pickG      = d3.Selector("g")
	pickBar    = d3.Selector(".bar")

	divTag  = d3.TagName("div")
	gTag    = d3.TagName("g")
	rectTag = d3.TagName("rect")
	textTag = d3.TagName("text")

	propWidth      = d3.PropertyName("width")
	propHeight     = d3.PropertyName("height")
	propXform      = d3.PropertyName("transform")
	propX          = d3.PropertyName("x")
	propY          = d3.PropertyName("y")
	propDy         = d3.PropertyName("dy")
	propClass      = d3.PropertyName("class")
	propTextAnchor = d3.PropertyName("text-anchor")
)

//convert object with string fields into one with parsed fields
func filterIntData(obj js.Object) js.Object {
	result := js.Global.Get("Object").New()
	result.Set("name", obj.Get("name").String())
	s := obj.Get("value").String()
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		console.Error("unable to parse ", s, " in the dataset: IGNORED")
		return nil
	}
	result.Set("value", i)
	return result
}

//this func must be coordinated with the filterData func above
func extractValue(obj js.Object) int64 {
	return int64(obj.Get("value").Int())
}

//implement part2, after all improvements
func part2_bars(width int64, barHeight int64) {
	x := d3.ScaleLinear().Range([]int64{0, width})

	chart := d3.Select(pickChart2).Attr(propWidth, width)

	//read sample data from the server
	d3.TSV("sample.tsv", filterIntData, func(err js.Object, data js.Object) {
		if !err.IsNull() {
			console.Error(err)
			return
		}
		x.Domain([]int64{0, d3.Max(data, extractValue)})
		chart.Attr(propHeight, barHeight*int64(data.Length()))

		bar := chart.SelectAll(pickG).Data(data).Enter().Append(gTag)
		bar.AttrFunc2S(propXform, func(d js.Object, i int64) string {
			return fmt.Sprintf("translate(0,%d)", i*barHeight)
		})
		rect := bar.Append(rectTag)
		rect.AttrFunc(propWidth, x.Func(extractValue)).Attr(propHeight, barHeight-1)

		text := bar.Append(textTag)
		text.AttrFunc(propX, func(d js.Object) int64 {
			return x.Linear(d, extractValue) - int64(3)
		})

		text.Attr(propY, barHeight/2).AttrS(propDy, ".35em")
		text.Text(func(d js.Object) string {
			return fmt.Sprintf("%s:%d", d.Get("name"), extractValue(d))
		})

	})
}

func filterFloatData(obj js.Object) js.Object {
	result := js.Global.Get("Object").New()
	result.Set("letter", obj.Get("letter").String())
	s := obj.Get("frequency").String()
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		console.Error("unable to parse ", s, " in the dataset: IGNORED")
		return nil
	}
	result.Set("frequency", f)
	return result

}

//this func must be coordinated with the filterData func above
func extractFreq(obj js.Object) float64 {
	return obj.Get("frequency").Float()
}

func extractLetter(obj js.Object) js.Object {
	return obj.Get("letter")
}

//horrific way to do map(func)
func extractAllLetters(obj js.Object) js.Object {
	result := js.Global.Get("Array").New()
	for i := 0; i < obj.Length(); i++ {
		result.SetIndex(i, extractLetter(obj.Index(i)))
	}
	return result
}

func part3_bars(overall_width, overall_height, top, right, bottom, left int64) {

	width := overall_width - left - right
	height := overall_height - top - bottom

	x := d3.ScaleOrdinal().RangeBands3([]int64{0, width}, 0.1)
	y := d3.ScaleLinear().Range([]int64{height, 0})

	xAxis := d3.NewAxis().ScaleO(x).Orient(d3.BOTTOM)
	yAxis := d3.NewAxis().Scale(y).Orient(d3.LEFT).Ticks(10, "%")

	chart := d3.Select(pickChart3).Attr(propWidth, width+left+right).
		Attr(propHeight, height+top+bottom).Append(gTag).
		AttrS(propXform, fmt.Sprintf("translate(%d,%d)", left, top))

	d3.TSV("letter_freq.tsv", filterFloatData, func(err js.Object, data js.Object) {
		if !err.IsNull() {
			console.Error(err)
			return
		}
		x.Domain(extractAllLetters(data))
		y.DomainF([]float64{0.0, d3.MaxF(data, extractFreq)})

		//AXES
		chart.Append(gTag).AttrS(propClass, "x axis").AttrS("transform",
			fmt.Sprintf("translate(0,%d)", height)).Call(xAxis)

		yText := chart.Append(gTag).AttrS(propClass, "y axis").Call(yAxis).Append(textTag)
		yText.AttrS(propXform, "rotate(-90)").Attr(propY, 6).
			AttrS(propDy, "0.71em").StyleS(propTextAnchor, "end").TextS("Frequency")

		//BAR
		rect := chart.SelectAll(pickBar).Data(data).Enter().Append(rectTag)
		rect.AttrS(propClass, "bar")
		rect.AttrFunc(propX, func(d js.Object) int64 {
			return x.Ordinal(d, extractLetter)
		})
		rect.AttrFuncF(propY, func(d js.Object) float64 {
			return y.LinearF(d, extractFreq)
		})
		rect.AttrFuncF(propHeight, func(obj js.Object) float64 {
			return float64(height) - y.LinearF(obj, extractFreq)
		})
		rect.AttrF(propWidth, x.RangeBandF())
	})
}

func main() {
	js.Global.Get("window").Set("onload", func() {
		part2_bars(420, 20)
		part3_bars(960, 500, 20, 30, 30, 40)
	})
}
