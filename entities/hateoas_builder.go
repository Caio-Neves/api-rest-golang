package entities

import "fmt"

type Hateoas struct {
	Links HateoasLinks `json:"_links"`
}

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
	Type   string `json:"type,omitempty"`
}

type HateoasLinks map[string]Link

type HateoasBuilder struct {
	links   HateoasLinks
	baseUrl string
}

func NewHateoasBuilder() *HateoasBuilder {
	return &HateoasBuilder{links: make(HateoasLinks)}
}

func (b *HateoasBuilder) AddBaseUrl(baseUrl string) *HateoasBuilder {
	b.baseUrl = baseUrl
	return b
}

func (b *HateoasBuilder) AddPost(name, href string) *HateoasBuilder {
	return b.add(name, href, "POST", "application/json")
}

func (b *HateoasBuilder) AddPatch(name, href string) *HateoasBuilder {
	return b.add(name, href, "PATCH", "application/json")
}

func (b *HateoasBuilder) AddGet(name, href string) *HateoasBuilder {
	return b.add(name, href, "GET", "")
}

func (b *HateoasBuilder) AddDelete(name, href string) *HateoasBuilder {
	return b.add(name, href, "DELETE", "")
}

func (b *HateoasBuilder) AddPut(name, href string) *HateoasBuilder {
	return b.add(name, href, "PUT", "application/json")
}

func (b *HateoasBuilder) add(name, href, method, contentType string) *HateoasBuilder {
	b.links[name] = Link{Href: fmt.Sprintf("%s%s", b.baseUrl, href), Method: method, Type: contentType}
	return b
}

func (b *HateoasBuilder) Build() Hateoas {
	return Hateoas{Links: b.links}
}
