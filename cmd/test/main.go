// Copyright 2024 JetERA Creative
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fyne.io/fyne/v2/app"
	. "github.com/jellytea/formui"
)

func main() {
	w := Window{Window: app.New().NewWindow("Test")}

	items := []any{"Item 1"}

	var listEdit *ListEdit
	listEdit = NewListEdit(items, func(e any) string {
		return e.(string)
	}, []string{"Add"}, func() {
		listEdit.Append("Item new")
	})

	w.SetContent(listEdit.Widget)

	w.ShowAndRun()
}
