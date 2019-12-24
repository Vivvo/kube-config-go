package models

/*Pod Struct */
type Pod struct {
	Items []struct {
		Metadata struct {
			GeneratedName string
			Name          string
			Namespace     string
			Labels        struct {
				App string
			}
		}
		Spec struct {
			Containers []struct {
				Name  string
				Image string
				Ports []struct {
					ContainerPort int
					Name          string
					Protocol      string
				}
			}
		}
	}
}
