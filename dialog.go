// Copyright 2024 The Contaiman Author
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package formui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ResizeDialog(weight, height int, d dialog.Dialog) dialog.Dialog {
	d.Resize(fyne.NewSize(float32(weight), float32(height)))
	return d
}

func (w Window) ShowError(err error) { dialog.NewInformation("Error", err.Error(), w.Window).Show() }

func (w Window) NewDialog(title string, object fyne.CanvasObject) dialog.Dialog {
	return dialog.NewCustom(title, "Close", object, w.Window)
}

func (w Window) NewDialogWithConfirm(title, confirm, dismiss string, callback func(), object fyne.CanvasObject) dialog.Dialog {
	return dialog.NewCustomConfirm(title, confirm, dismiss, object, func(b bool) {
		if b {
			callback()
		}
	}, w.Window)
}

func (w Window) NewFormDialog(title string, items ...*widget.FormItem) dialog.Dialog {
	return dialog.NewCustom(title, "Close", widget.NewForm(items...), w.Window)
}

func (w Window) NewFormDialogWithConfirm(title, confirm, dismiss string, onConfirm func(), items ...*widget.FormItem) dialog.Dialog {
	return dialog.NewForm(title, confirm, dismiss, items, func(b bool) {
		if b {
			onConfirm()
		}
	}, w.Window)
}

func NewListView(items []string) *widget.List {
	return widget.NewList(func() int {
		return len(items)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		object.(*widget.Label).SetText(items[id])
	})
}

func NewListSelect(onSelected func(i int, v string), items []string) *widget.List {
	list := widget.NewList(func() int {
		return len(items)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		object.(*widget.Label).SetText(items[id])
	})
	list.OnSelected = func(id widget.ListItemID) { onSelected(id, items[id]) }
	return list
}

func (w Window) NewFileOpenDialog(onSelected func(path string)) dialog.Dialog {
	return dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
		if err != nil {
			w.ShowError(err)
			return
		}
		if closer != nil {
			onSelected(closer.URI().Path())
			_ = closer.Close()
		}
	}, w.Window)
}

func (w Window) NewDirOpenDialog(onSelected func(path string)) dialog.Dialog {
	return dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			w.ShowError(err)
			return
		}
		if uri != nil {
			list, err := uri.List()
			if err != nil {
				w.ShowError(err)
				return
			}
			onSelected(list[0].Path())
		}
	}, w.Window)
}

func (w Window) NewFileOpenEntry(entry *widget.Entry) *fyne.Container {
	return container.NewVBox(entry, widget.NewButton("...", func() {
		w.NewFileOpenDialog(func(path string) {
			entry.SetText(path)
		}).Show()
	}))
}

func (w Window) NewDirOpenEntry(entry *widget.Entry) *fyne.Container {
	return container.NewVBox(entry, widget.NewButton("...", func() {
		w.NewDirOpenDialog(func(path string) {
			entry.SetText(path)
		}).Show()
	}))
}

func (w Window) NewEntryWithFill(entry *widget.Entry, title string, callback func(d dialog.Dialog), items []string) *fyne.Container {
	return container.NewVBox(
		entry,
		widget.NewButton("...", func() {
			callback(w.NewDialog(title, NewListSelect(func(i int, v string) {
				entry.SetText(v)
			}, items)))
		}),
	)
}
