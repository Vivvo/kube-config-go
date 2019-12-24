package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"./models"
	"./settings"
	"./utils"
	"github.com/go-yaml/yaml"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var kubeSystemWatcher = core.NewQFileSystemWatcher(nil)
var driveIntoNamespaces = true
var allowNamespaceFiltering = true

//RGC: https://github.com/ahmetb/kubectx/blob/master/kubens
func switchNamespace(ctx string, namespace string) error {
	kubeSystemWatcher.BlockSignals(true)
	err := utils.SetContext(ctx)
	if err == nil {
		kubeSystemWatcher.BlockSignals(false)
		err = utils.SetNamespace(ctx, namespace)
	}
	return err
}

func setupMenu(qApp *widgets.QApplication) *widgets.QMenu {
	menu := widgets.NewQMenu(nil)
	currentContext := utils.GetCurrentContext()
	currentNamespace := utils.GetCurrentNamespace()
	b, err := ioutil.ReadFile(utils.GetKubeConfigPath())

	if err == nil {
		config := models.KubeConfigStruct{}
		yaml.Unmarshal(b, &config)

		for _, cluster := range config.Clusters {
			submenuAction := menu.AddAction(cluster.Name)
			submenuAction.SetCheckable(true)
			if cluster.Name == currentContext {
				submenuAction.SetChecked(true)
			}
			submenu := widgets.NewQMenu2(cluster.Name,nil)
			submenuAction.ConnectTriggered(func(triggered bool) {
				//ctx := action.Data().ToString()
				//name := action.Text()
				fmt.Printf("switchNamespace:\n")
				//err := switchNamespace(ctx, name)
				//fmt.Printf(err.Error())
			})
			namespace := utils.GetNamespaces(cluster.Name)
			for _, item := range namespace.Items {
				namespace := item.Metadata.Name

				actionMenu := submenu.AddMenu2(namespace)
				action := actionMenu.MenuAction()
				action.SetCheckable(true)
				action.ConnectTriggered(func(trigggered bool) {
					//ctx := action.Data().ToString()
					//name := action.Text()
					fmt.Printf("switchNamespace:\n")
					//err := switchNamespace(ctx, name)
					//fmt.Printf(err.Error())
				})

				if cluster.Name == currentContext && namespace == currentNamespace {
					action.SetChecked(true)
				}
				action.SetData(core.NewQVariant12(cluster.Name))

				if driveIntoNamespaces {
					podsMenu := widgets.NewQMenu2("Pods", actionMenu)
					pod := utils.GetPods(cluster.Name, namespace)

					for _, pod := range pod.Items {
						podMenu := podsMenu.AddMenu2(pod.Metadata.Name)

						if len(pod.Spec.Containers) > 1 {
							for _, container := range pod.Spec.Containers {
								containerMenu := podMenu.AddMenu2(container.Name)
								containerMenu.SetToolTip(container.Image)
								for _, port := range container.Ports {
									portAction := containerMenu.AddAction(fmt.Sprintf("%s -> %d", port.Name, int(port.ContainerPort)))

									portAction.ConnectTriggered(func(triggered bool) {
										print("port triggered")
									})
								}
							}
						} else {
							for _, port := range pod.Spec.Containers[0].Ports {
								portAction := podMenu.AddAction(fmt.Sprintf("%s -> %d", port.Name, int(port.ContainerPort)))
								portAction.ConnectTriggered(func(triggered bool) {
									fmt.Printf("port triggered\n")
								})
							}
						}
					}
					actionMenu.AddMenu(podsMenu)
				}

				if allowNamespaceFiltering {
					actionMenu.AddAction("Search Logs...")
					//namespaceFilter.SetData(core.NewQVariant23()
					//namespaceFilter.ConnectTriggered(func(action *widgets.QAction) {

					//	})
				}
			}
			submenu.AddSeparator()
			submenu.AddAction("Filter Logs...")
			submenuAction.SetMenu(submenu)
			menu.AddMenu(submenu)
		}

		menu.AddSeparator()
		menu.AddAction("Add KubeConfig").ConnectTriggered(func(checked bool) {
		})

		menu.AddSeparator()

		menu.AddAction("Settings").ConnectTriggered(func(checked bool) {
			dialog, _ := settings.Dialog(qApp)
			dialog.Exec()
			dialog.SetFocus2()
			dialog.Raise()
		})

		menu.AddAction("Exit").ConnectTriggered(func(checked bool) {
			qApp.Exit(1)
		})

		return menu
	}
	return menu
}

func main() {
	qApp := widgets.NewQApplication(len(os.Args), os.Args)
	qApp.SetQuitOnLastWindowClosed(false)
	qApp.SetApplicationDisplayName("Kube Configuration Plugin")
	qApp.SetWindowIcon(gui.NewQIcon5(":images/icon.svg"))
	qApp.SetOrganizationName("")
	if _, err := os.Stat(utils.GetKubeConfigPath()); os.IsNotExist(err) {
		panic("Cannot find config path")
	}
	_, err := exec.LookPath(utils.GetKubeCtl())
	if err != nil {
		panic(fmt.Sprintf("Cannot find kubectl: %s", err.Error()))
	}

	systemTray := widgets.NewQSystemTrayIcon(qApp)
	systemTray.SetIcon(gui.NewQIcon5(":images/icon.svg"))
	systemTray.SetContextMenu(setupMenu(qApp))
	systemTray.SetVisible(true)
	kubeSystemWatcher.AddPath(utils.GetKubeConfigPath())
	kubeSystemWatcher.ConnectFileChanged(func(path string) {
		menu := setupMenu(qApp)
		systemTray.SetContextMenu(menu)

		systemTray.ShowMessage(
			"Kubectl Configuration has changed",
			fmt.Sprintf("The kubectl configuration changed.  The current context/namespace is: %s - %s", utils.GetCurrentContext(), utils.GetCurrentNamespace()),
			widgets.QSystemTrayIcon__Information,
			10000,
		)
	})
	qApp.Exec()
}
