package github

import (
	"bytes"
	"errors"
	"net/http"

	"golang.org/x/net/html"
)

const SpecLocation = "https://raw.githubusercontent.com/tc39/ecma262/refs/heads/main/spec.html"

type GithubSpecProvider struct {
	Content  map[string]string
	Sections []string
}

func NewGithubSpecProvider() *GithubSpecProvider {
	return &GithubSpecProvider{}
}

func (g *GithubSpecProvider) Initialize() error {
	resp, err := http.Get(SpecLocation)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("failed to fetch spec")
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	if g.Content == nil {
		g.Content = make(map[string]string)
	}
	g.Sections = g.Sections[:0]

	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val != "" {
					var buf bytes.Buffer
					_ = html.Render(&buf, n)
					id := a.Val
					g.Sections = append(g.Sections, id)
					g.Content[id] = buf.String()
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)

	return nil
}

func (g *GithubSpecProvider) GetSpec(specPath string) (string, error) {
	if g.Content == nil {
		if err := g.Initialize(); err != nil {
			return "", errors.Join(errors.New("spec provider not initialized"), err)
		}
	}

	if content, exists := g.Content["sec-"+specPath]; exists {
		return content, nil
	}

	return "", errors.New("spec section not found")
}

func (g *GithubSpecProvider) SpecForIntrinsic(intrinsic string) (string, error) {
	if g.Content == nil {
		if err := g.Initialize(); err != nil {
			return "", errors.Join(errors.New("spec provider not initialized"), err)
		}
	}

	return "", errors.New("not implemented")
}

func (g *GithubSpecProvider) SearchSpec(query string) ([]string, error) {
	if g.Content == nil {
		if err := g.Initialize(); err != nil {
			return nil, errors.Join(errors.New("spec provider not initialized"), err)
		}
	}

	return nil, errors.New("not implemented")
}
