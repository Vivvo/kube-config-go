package settings

import (
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"../utils"
)

/*Dialog...*/
func Dialog(app *widgets.QApplication) (*widgets.QDialog, *core.QSettings) {
	dialog := widgets.NewQDialog(nil, core.Qt__Dialog)
	dialog.SetWindowTitle("KubeConfig Configuration")
	settings := core.NewQSettings4(utils.GetSettingsPath(), core.QSettings__IniFormat, nil)

	_ = settings.Value("icon/type", core.NewQVariant12("icon"))
	var rows = 0

	mainLayout := widgets.NewQGridLayout2()
	formLayout := widgets.NewQFormLayout(nil)

	comboBox := widgets.NewQComboBox(nil)
	comboBox.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding)
	var actions []*widgets.QAction

	action1 := widgets.NewQAction2("Icon", nil)
	action1.SetText("Icon")
	action1.SetData(core.NewQVariant12("icon"))
	actions = append(actions, action1)

	action2 := widgets.NewQAction2("Icon/Context", nil)
	action2.SetText("Icon/Context")
	action2.SetData(core.NewQVariant12("context"))
	actions = append(actions, action2)

	action3 := widgets.NewQAction2("Icon/Context/Namespace", nil)
	action3.SetText("Icon/Context/Namespace")
	action3.SetData(core.NewQVariant12("namespace"))
	actions = append(actions, action3)

	comboBox.AddItem(action1.Text(), action1.Data())
	comboBox.AddItem(action2.Text(), action2.Data())
	comboBox.AddItem(action3.Text(), action3.Data())

	comboBox.ConnectActivated2(func(activatedString string) {
		settings.SetValue("icon/type", core.NewQVariant12(activatedString))
	})

	formLayout.AddRow(widgets.NewQLabel2("Display as:", nil, core.Qt__Widget), comboBox)
	formLayout.SetHorizontalSpacing(0)
	formLayout.SetVerticalSpacing(0)
	formLayout.SetContentsMargins(0, 0, 0, 0)
	formLayout.SetSizeConstraint(widgets.QLayout__SetNoConstraint)

	mainLayout.AddItem(formLayout, 0, 0, 0, -1, core.Qt__AlignLeft)
	rows = rows + 1

	integrations := setupIntegrations()
	for _, obj := range integrations {
		rows = rows + 1
		obj.SetSizePolicy2(widgets.QSizePolicy__MinimumExpanding, widgets.QSizePolicy__MinimumExpanding)
		mainLayout.AddWidget3(obj, rows, 0, 1, -1, core.Qt__AlignLeft)
	}

	buttonBox := widgets.NewQDialogButtonBox3(widgets.QDialogButtonBox__Ok|widgets.QDialogButtonBox__Cancel, nil)
	buttonBox.ConnectAccepted(func() {

	})

	buttonBox.ConnectRejected(func() {
		dialog.Close()
	})

	mainLayout.AddWidget3(buttonBox, rows, 1, -1, -1, core.Qt__AlignLeft)

	dialog.SetLayout(mainLayout)
	//dialog.SetMaximumSize(mainLayout.MaximumSize())
	//dialog.Layout().SetSizeConstraint( widgets.QLayout__SetFixedSize )

	return dialog, settings
}

func setupIntegrations() []*widgets.QGroupBox {
	var objs []*widgets.QGroupBox

	if strings.Contains(core.NewQSysInfoFromPointer(nil).ProductType(), "osx") {
		spotlightFrame := widgets.NewQGroupBox2("Spotlight Integration", nil)
		spotlightFrame.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding)
		spotlightFrame.SetContentsMargins(0, 0, 0, 0)
		spotlightLayout := widgets.NewQFormLayout(nil)
		spotlightLayout.SetSizeConstraint(widgets.QLayout__SetNoConstraint)
		spotlightLayout.AddRow(widgets.NewQCheckBox(nil), widgets.NewQLabel2("Enable Spotlight Integration", nil, core.Qt__Widget))
		spotlightFrame.SetLayout(spotlightLayout)
		spotlightFrame.SetEnabled(false)
		objs = append(objs, spotlightFrame)

		alfredFrame := widgets.NewQGroupBox2("Alfred Integration", nil)
		alfredFrame.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding)
		alfredFrame.SetContentsMargins(0, 0, 0, 0)
		alfredLayout := widgets.NewQVBoxLayout()
		alfredLayout.SetSizeConstraint(widgets.QLayout__SetNoConstraint)
		alfredLayout.SetContentsMargins(0, 0, 0, 0)
		rowLayout := widgets.NewQHBoxLayout()
		rowLayout.AddWidget(widgets.NewQCheckBox(nil), 0, core.Qt__AlignLeft)
		rowLayout.AddWidget(widgets.NewQLabel2("Enable Alfred Integration", nil, core.Qt__Widget), 0, core.Qt__AlignLeft)
		alfredLayout.AddItem(rowLayout)
		//alfredLayout.AddRow(widgets.NewQCheckBox(nil), widgets.NewQLabel2("Enable Alfred Integration", nil, core.Qt__Widget))
		alfredFrame.SetLayout(alfredLayout)
		alfredFrame.SetEnabled(false)
		objs = append(objs, alfredFrame)
	}
	return objs
}
