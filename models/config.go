package models

/*Struct*/
type KubeConfigStruct struct {
	APIVersion     string `yaml:"apiVersion"`
	CurrentContext string `json:"current-context" yaml:"current-context"`

	Clusters []struct {
		Cluster struct {
			APIVersion string `json:"api-version" yaml:"api-version"`
			Server     string `json:"server" yaml:"server"`
		} `json:"cluster" yaml:"cluster"`

		Name string `json:"name" yaml:"name"`
	} `json:"clusters" yaml:"clusters"`

	Contexts []struct {
		Context struct {
			Cluster   string `json:"cluster" yaml:"cluster"`
			Namespace string `json:"namespace" yaml:"namespace"`
			User      string `json:"user" yaml:"user"`
		} `json:"context" yaml:"context"`
		Name string `json:"name" yaml:"name"`
	} `json:"contexts" yaml:"contexts"`

	Kind string `json:"kind" yaml:"kind"`

	Preferences struct {
		Colors bool `json:"colors" yaml:"colors"`
	} `json:"preferences" yaml:"preferences"`

	Users []struct {
		Name string `json:"name" yaml:"name"`

		User struct {
			Token string `json:"token" yaml:"token"`
		} `json:"user" yaml:"user"`
	} `json:"users" yaml:"users"`
}
