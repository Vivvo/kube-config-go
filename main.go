package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type KubeConfigStruct struct {
	CurrentContext string `yaml:"current-context"`

	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			APIVersion string `yaml:"api-version"`

			Server string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster   string `yaml:"cluster"`
			Namespace string `yaml:"namespace"`
			User      string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	Kind        string `yaml:"kind"`
	Preferences struct {
		Colors bool `yaml:"colors"`
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`

		User struct {
			Token string `yaml:"token"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func getKubeCtl() (string) {
	if strings.Contains(core.NewQSysInfoFromPointer(nil).ProductType(), "win") {
		return "kubectl.ext"
	} else {
		return "kubectl"
	}
}

func settingsDialog() (*widgets.QDialog) {
	return widgets.NewQDialog(nil, core.Qt__Dialog)
}

//RGC: https://github.com/ahmetb/kubectx/blob/master/kubens
func switchNamespace(ctx string, namespace string) {
	args := []string{
		"set-context",
		ctx,
		fmt.Sprintf("--namespace=\"%s\"", namespace),
	}
	cmd := exec.Command(getKubeCtl(), args...)
	_, _ = cmd.Output()
}

//RGC: https://github.com/ahmetb/kubectx/blob/master/kubens
func getCurrentContext() string{
	args := []string{
		"config",
		"current-context",
	}
	cmd := exec.Command(getKubeCtl(), args...)
	out, err := cmd.Output()
	if err == nil {
		fmt.Printf("current-context: %s", out)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func getCurrentNamespace() string{
	args := []string{
		"config",
		"view",
		"--minify",
		"--output",
		"jsonpath={..namespace}",
	}
	cmd := exec.Command(getKubeCtl(), args...)
	out, err := cmd.Output()
	if err == nil {
		fmt.Printf("current-namespace: %s\n", out)
	}
	return strings.TrimSuffix(string(out), "\n")
}

//RGC: https://github.com/ahmetb/kubectx/blob/master/kubens
func getNamespaces(cluster string) []string {
	args := []string{
		"get",
		"namespaces",
		fmt.Sprintf("--context=%s", cluster),
	}
	objs := make([]string, 0)
	cmd := exec.Command(getKubeCtl(), args...)
	out, err := cmd.Output()
	if err == nil {
		output := strings.Split(string(out), "\n")
		for i, o := range output {
			if i != 0 {
				objs = append(objs, strings.Split(o," ")[0])
			}
		}
	}
	return objs
}

func setupMenu() *widgets.QMenu {
	menu := widgets.NewQMenu(nil)
	currentContext := getCurrentContext()
	currentNamespace := getCurrentNamespace()
	fmt.Printf("%s,%s\n",currentContext,currentNamespace)
	b, err := ioutil.ReadFile(getKubeConfigPath())

	if err == nil {
		config := KubeConfigStruct{}
		yaml.Unmarshal(b, &config)

		for _, cluster := range config.Clusters {
			submenu := widgets.NewQMenu2(cluster.Name, nil)
			if cluster.Name == currentContext{
				submenu.SetStyleSheet("font-weight:bold;color:red")
			}
			namespaces := getNamespaces(cluster.Name)
			for _, namespace := range namespaces {
				action := submenu.AddAction(namespace)
				if cluster.Name == currentContext && namespace == currentNamespace{
					action.SetChecked(true)
				}
				action.SetData(core.NewQVariant12(cluster.Name))
			}
			submenu.ConnectTriggered(func(action *widgets.QAction) {
				fmt.Printf("%v - %v",action.Data().ToString(),action.Text())
			})
			menu.AddMenu(submenu)
		}
		menu.AddSeparator()
		menu.AddAction("Settings").ConnectTriggered(func(checked bool) {

		})

		menu.AddAction("Exit").ConnectTriggered(func(checked bool) {
			os.Exit(1)
		})

		return menu
	}
	return menu
}

func getPath() (string) {
	var path = ""
	if strings.Contains(core.NewQSysInfoFromPointer(nil).ProductType(), "win") {
		path = fmt.Sprintf("%s\\.kube-config-go", core.NewQStandardPathsFromPointer(nil).WritableLocation(core.QStandardPaths__ConfigLocation))
	} else {
		path = fmt.Sprintf("%s/.kube-config-go", core.NewQStandardPathsFromPointer(nil).WritableLocation(core.QStandardPaths__HomeLocation))
	}
	return path
}

func getKubeConfigPath() (string) {
	var path = ""
	if strings.Contains(core.NewQSysInfoFromPointer(nil).ProductType(), "win") {
		path = fmt.Sprintf("%s\\.kube\\config", core.NewQStandardPathsFromPointer(nil).WritableLocation(core.QStandardPaths__ConfigLocation))
	} else {
		path = fmt.Sprintf("%s/.kube/config", core.NewQStandardPathsFromPointer(nil).WritableLocation(core.QStandardPaths__HomeLocation))
	}
	return path
}

func main() {
	qApp := widgets.NewQApplication(len(os.Args), os.Args)

	if _, err := os.Stat(getPath()); os.IsNotExist(err) {
		os.Mkdir(getPath(), 0777)
	}

	systemTray := widgets.NewQSystemTrayIcon(nil)
	systemTray.SetIcon(gui.NewQIcon5(":images/icon.svg"))
	systemTray.SetContextMenu(setupMenu())
	systemTray.SetVisible(true)

	kubeSystemWatcher := core.NewQFileSystemWatcher(nil)
	kubeSystemWatcher.AddPath(getKubeConfigPath())
	kubeSystemWatcher.ConnectFileChanged(func(path string) {
		fmt.Printf("file changed: %s\n", path)
		menu := setupMenu()
		systemTray.SetContextMenu(menu)
	})


	qApp.Exec()
}
