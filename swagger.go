package lux

import (
	"github.com/snowmerak/lux/swagger"
)

type Swagger struct {
	inner *swagger.Swagger
}

func (s *Swagger) SetDescription(description string) {
	s.inner.Info.Description = description
}

func (s *Swagger) SetEmail(email string) {
	s.inner.Info.Contact.Email = email
}

func (s *Swagger) SetTitle(title string) {
	s.inner.Info.Title = title
}

func (s *Swagger) SetLicense(license string, url string) {
	s.inner.Info.License.Name = license
	s.inner.Info.License.URL = url
}

func (s *Swagger) SetVersion(version string) {
	s.inner.Info.Version = version
}

func (s *Swagger) AddDefinition(name string, definition swagger.Definition) {
	s.inner.Definitions[name] = definition
}

func (s *Swagger) AddSecurityDefinition(name string, definition swagger.SecurityDefinition) {
	s.inner.SecurityDefinitions[name] = definition
}
