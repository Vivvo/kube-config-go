package models

type Namespace struct {
	Items []struct {
		Metadata struct {
			Name string
			UID  string
		}
	}
}
