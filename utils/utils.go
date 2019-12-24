package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"../models"
	"github.com/therecipe/qt/core"
)

func GetKubeCtl() string {
	if strings.Contains(core.NewQSysInfoFromPointer(nil).ProductType(), "win") {
		return "kubectl.exe"
	}
	return "kubectl"
}

func GetSettingsPath() string {
	var path = ""
	if strings.Contains(core.NewQSysInfoFromPointer(nil).ProductType(), "win") {
		path = fmt.Sprintf("%s\\.kube-config-go\\config", core.NewQStandardPathsFromPointer(nil).WritableLocation(core.QStandardPaths__ConfigLocation))
	} else {
		path = fmt.Sprintf("%s/.kube-config-go/config", core.NewQStandardPathsFromPointer(nil).WritableLocation(core.QStandardPaths__HomeLocation))
	}
	return path
}
func GetKubeConfigPath() string {
	var path = ""
	if strings.Contains(core.NewQSysInfoFromPointer(nil).ProductType(), "win") {
		path = fmt.Sprintf("%s\\.kube\\config", core.NewQStandardPathsFromPointer(nil).WritableLocation(core.QStandardPaths__ConfigLocation))
	} else {
		path = fmt.Sprintf("%s/.kube/config", core.NewQStandardPathsFromPointer(nil).WritableLocation(core.QStandardPaths__HomeLocation))
	}
	return path
}

func SetContext(ctx string) error {
	args := []string{
		"config",
		"use-context",
		ctx,
	}

	cmd := exec.Command(GetKubeCtl(), args...)
	_, err := cmd.Output()
	return err
}

func SetNamespace(ctx string, namespace string) error {
	args := []string{
		"config",
		"set-context",
		ctx,
		"--namespace",
		namespace,
	}

	cmd := exec.Command(GetKubeCtl(), args...)
	_, err := cmd.Output()
	return err
}

//RGC: https://github.com/ahmetb/kubectx/blob/master/kubens
func GetCurrentContext() string {
	args := []string{
		"config",
		"current-context",
	}
	cmd := exec.Command(GetKubeCtl(), args...)
	out, err := cmd.Output()
	if err == nil {
		fmt.Printf("current-context: %s", out)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func GetCurrentNamespace() string {
	args := []string{
		"config",
		"view",
		"--minify",
		"--output",
		"jsonpath={..namespace}",
	}
	cmd := exec.Command(GetKubeCtl(), args...)
	out, err := cmd.Output()
	if err == nil {
		fmt.Printf("current-namespace: %s\n", out)
	}
	return strings.TrimSuffix(string(out), "\n")
}

//RGC: https://github.com/ahmetb/kubectx/blob/master/kubens
func GetNamespaces(cluster string) models.Namespace {
	args := []string{
		"get",
		"namespaces",
		fmt.Sprintf("--context=%s", cluster),
		"--output",
		"json",
	}
	cmd := exec.Command(GetKubeCtl(), args...)
	out, err := cmd.Output()
	namespaces := models.Namespace{}
	if err == nil {
		json.Unmarshal(out, &namespaces)
	}
	return namespaces
}

func GetPods(ctx string, namespace string) models.Pod {
	args := []string{
		"get", "pods",
		fmt.Sprintf("--context=%s", ctx),
		fmt.Sprintf("--namespace=%s", namespace),
		"--output",
		"json",
	}
	cmd := exec.Command(GetKubeCtl(), args...)
	out, _ := cmd.Output()
	pod := models.Pod{}
	json.Unmarshal(out, &pod)
	return pod
}
