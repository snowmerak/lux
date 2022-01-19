package swagger

type Swagger struct {
	SwaggerVersion string `json:"swagger,omitempty"`
	Info           struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Version     string `json:"version,omitempty"`
		Contact     struct {
			Email string `json:"email,omitempty"`
		} `json:"contact,omitempty"`
		License struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"license,omitempty"`
	} `json:"info,omitempty"`
	Host     string `json:"host,omitempty"`
	BasePath string `json:"basePath,omitempty"`
	Tags     []struct {
		Name         string `json:"name,omitempty"`
		Description  string `json:"description,omitempty"`
		ExternalDocs struct {
			Description string `json:"description,omitempty"`
			URL         string `json:"url,omitempty"`
		} `json:"externalDocs,omitempty"`
	} `json:"tags,omitempty"`
	Schemes             []string                      `json:"schemes,omitempty"`
	Paths               map[Path]map[Method]Router    `json:"paths,omitempty"`
	SecurityDefinitions map[string]SecurityDefinition `json:"securityDefinitions,omitempty"`
	Definitions         map[string]Definition         `json:"definitions,omitempty"`
	ExternalDocs        struct {
		Description string `json:"description,omitempty"`
		URL         string `json:"url,omitempty"`
	} `json:"externalDocs,omitempty"`
}

type Path string
type Method string

const (
	GET     Method = "get"
	POST    Method = "post"
	PUT     Method = "put"
	DELETE  Method = "delete"
	OPTIONS Method = "options"
	PATCH   Method = "patch"
)

type Parameter struct {
	In          string            `json:"in,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Required    bool              `json:"required,omitempty"`
	Type        string            `json:"type,omitempty"`
	Minimum     float64           `json:"minimum,omitempty"`
	Maximum     float64           `json:"maximum,omitempty"`
	Format      string            `json:"format,omitempty"`
	Schema      map[string]string `json:"schema,omitempty"`
	Items       struct {
		Type   string `json:"type,omitempty"`
		Format string `json:"format,omitempty"`
	} `json:"items,omitempty"`
	CollectionFormat string `json:"collectionFormat,omitempty"`
}

type Router struct {
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Consumes    []string              `json:"consumes,omitempty"`
	Produces    []string              `json:"produces,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	Responses   map[string]Response   `json:"responses,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
}

type Response struct {
	Description string `json:"description,omitempty"`
	Schema      Schema `json:"schema,omitempty"`
}

type Schema struct {
	Type  string            `json:"type,omitempty"`
	Items map[string]string `json:"items,omitempty"`
}

type Definition struct {
	Type       string   `json:"type,omitempty"`
	Required   []string `json:"required,omitempty"`
	Properties map[string]struct {
		Type   string   `json:"type,omitempty"`
		Format string   `json:"format,omitempty"`
		Enum   []string `json:"enum,omitempty"`
	} `json:"properties,omitempty"`
}

type SecurityDefinition struct {
	Type string `json:"type,omitempty"`
	In   string `json:"in,omitempty"`
	Name string `json:"name,omitempty"`
}
