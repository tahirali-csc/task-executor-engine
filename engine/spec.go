package engine

type Metadata struct {
	UID       string `json:"uid,omitempty"`
	Namespace string
}

type VolumeMount struct {
	Name      string
	MountPath string
	ClaimName string
}

type Spec struct {
	Metadata Metadata
	Image    string
	Command  []string
	Args     []string
	Volumes  []VolumeMount
}
