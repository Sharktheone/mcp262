package github

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const SpecLocation = "https://raw.githubusercontent.com/tc39/ecma262/refs/heads/main/spec.html"
const IntlSpecLocation = "https://tc39.es/ecma402/"

type GithubSpecProvider struct {
	Content  map[string]string
	Sections []string
}

func NewGithubSpecProvider() *GithubSpecProvider {
	return &GithubSpecProvider{}
}

func (g *GithubSpecProvider) initializeSpec(url string) error {
	resp, err := http.Get(url)
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

	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val != "" {
					var buf bytes.Buffer
					_ = html.Render(&buf, n)
					id := a.Val

					id = strings.TrimPrefix(id, "sec-")

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

func (g *GithubSpecProvider) Initialize() error {
	if g.Content == nil {
		g.Content = make(map[string]string)
	}
	g.Sections = g.Sections[:0]

	if err := g.initializeSpec(IntlSpecLocation); err != nil {
		return err
	}

	if err := g.initializeSpec(SpecLocation); err != nil {
		return err
	}

	return nil
}

func (g *GithubSpecProvider) GetSpec(specPath string) (string, error) {
	if g.Content == nil {
		if err := g.Initialize(); err != nil {
			return "", errors.Join(errors.New("spec provider not initialized"), err)
		}
	}

	if content, exists := g.Content[specPath]; exists {
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

	intrinsic = strings.TrimPrefix(intrinsic, "%")
	intrinsic = strings.TrimSuffix(intrinsic, "%")
	intrinsic = strings.ReplaceAll(intrinsic, "[[", "")
	intrinsic = strings.ReplaceAll(intrinsic, "]]", "")
	intrinsic = strings.ReplaceAll(intrinsic, "(", "")
	intrinsic = strings.ReplaceAll(intrinsic, ")", "")
	intrinsic = strings.ToLower(intrinsic)

	if content, exists := g.Content[intrinsic]; exists {
		return content, nil
	}

	return "", errors.New("spec section not found")
}

func (g *GithubSpecProvider) SearchSpec(query string) ([]string, error) {
	if g.Content == nil {
		if err := g.Initialize(); err != nil {
			return nil, errors.Join(errors.New("spec provider not initialized"), err)
		}
	}

	matches := make([]string, 0)
	for id, content := range g.Content {
		if strings.Contains(strings.ToLower(content), strings.ToLower(query)) {
			matches = append(matches, id)
		}
	}

	return matches, nil
}

func (g *GithubSpecProvider) SearchSections(query string) ([]string, error) {
	if g.Content == nil {
		if err := g.Initialize(); err != nil {
			return nil, errors.Join(errors.New("spec provider not initialized"), err)
		}
	}

	matches := make([]string, 0)
	for id, content := range g.Content {
		if strings.Contains(strings.ToLower(content), strings.ToLower(query)) {
			matches = append(matches, id)
		}
	}

	return matches, nil
}
