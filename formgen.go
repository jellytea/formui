// Copyright 2024 JetERA Creative
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package formui

import (
	"fyne.io/fyne/v2/widget"
	"reflect"
)

func BuildFormFromStruct(v any) *widget.Form {
	var formItems []*widget.FormItem

	t := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(v)))
	for _, f := range reflect.VisibleFields(t) {
		key := f.Tag.Get("formui")
		if key == "" {
			key = f.Name
		}
		v, _ := t.FieldByName(f.Name)
		entry := widget.NewEntry()
		entry.OnChanged = func(s string) {
			switch f.Type.Kind() {
			case reflect.String:
				reflect.ValueOf(v).SetString(entry.Text)
			default:
			}
		}
		formItems = append(formItems, widget.NewFormItem(key, entry))
	}

	return widget.NewForm(formItems...)
}
