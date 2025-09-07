package github

type GithubSpecProvider struct {
	Sections map[string]string
}

func (g GithubSpecProvider) GetSpec(specPath string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (g GithubSpecProvider) SpecForIntrinsic(intrinsic string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (g GithubSpecProvider) SearchSpec(query string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
