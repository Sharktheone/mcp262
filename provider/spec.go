package provider

type SpecProvider interface {
	GetSpec(specPath string) (string, error)
	SpecForIntrinsic(intrinsic string) (string, error)
	SearchSpec(query string) ([]string, error)
	SearchSections(query string) ([]string, error)
}

var Spec SpecProvider

func SetSpecProvider(s SpecProvider) {
	Spec = s
}
