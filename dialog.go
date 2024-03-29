// Copyright 2024 JetERA Creative
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

func (w Window) NewDialogWithWait(title string, object fyne.CanvasObject) (dialog.Dialog, func()) {
	var (
		d           dialog.Dialog
		closeButton *widget.Button
	)

	closeButton = widget.NewButton("Close", func() { d.Hide() })
	closeButton.Hide()

	d = dialog.NewCustomWithoutButtons(title, container.NewBorder(nil, closeButton, nil, nil, object), w.Window)

	return d, func() { closeButton.Show() }
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

type ListEdit struct {
	Elements []any
	Widget   fyne.CanvasObject
	Update   func()
}

func (l *ListEdit) Append(e any) {
	l.Elements = append(l.Elements, e)
	l.Update()
}

func NewListEdit(elements []any, toString func(e any) string, labels []string, add ...func()) *ListEdit {
	listEdit := &ListEdit{
		Elements: elements,
		Widget:   nil,
		Update:   nil,
	}

	list := container.NewHScroll(NewListView(nil))

	selected := -1

	var buttons []fyne.CanvasObject

	for i, label := range labels {
		f := add[i]
		buttons = append(buttons, widget.NewButton(label, func() {
			f()
			listEdit.Update()
		}))
	}

	buttons = append(buttons, widget.NewButton("Remove", func() {
		if selected < 0 {
			return
		}
		listEdit.Elements = append(listEdit.Elements[:selected], listEdit.Elements[selected+1:]...)
		listEdit.Update()
		selected = -1
	}))

	listEdit.Widget = container.NewBorder(nil, container.NewHBox(buttons...), nil, nil, list)
	listEdit.Update = func() {
		var items []string
		for _, e := range listEdit.Elements {
			items = append(items, toString(e))
		}
		list.Content = NewListSelect(func(i int, v string) { selected = i }, items)
		list.Refresh()
	}

	listEdit.Update()

	return listEdit
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
