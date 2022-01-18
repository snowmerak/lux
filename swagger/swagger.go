package swagger

type Swagger struct {
	SwaggerVersion string `json:"swagger"`
	Info           struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Version     string `json:"version"`
		Contact     struct {
			Email string `json:"email"`
		} `json:"contact"`
		License struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license"`
	} `json:"info"`
	Host     string `json:"host"`
	BasePath string `json:"basePath"`
	Tags     []struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		ExternalDocs struct {
			Description string `json:"description"`
			URL         string `json:"url"`
		} `json:"externalDocs"`
	} `json:"tags"`
	Schemes             []string                   `json:"schemes"`
	Paths               map[Path]map[Method]Router `json:"paths"`
	SecurityDefinitions map[string]struct {
		Type string `json:"type"`
		In   string `json:"in"`
		Name string `json:"name"`
	} `json:"securityDefinitions"`
	Definitions map[string]struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties map[string]struct {
			Type   string   `json:"type"`
			Format string   `json:"format"`
			Enum   []string `json:"enum"`
		} `json:"properties"`
	} `json:"definitions"`
	ExternalDocs struct {
		Description string `json:"description"`
		URL         string `json:"url"`
	} `json:"externalDocs"`
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
	In          string            `json:"in"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Required    bool              `json:"required"`
	Type        string            `json:"type"`
	Minimum     float64           `json:"minimum"`
	Maximum     float64           `json:"maximum"`
	Format      string            `json:"format"`
	Schema      map[string]string `json:"schema"`
	Items       struct {
		Type   string `json:"type"`
		Format string `json:"format"`
	} `json:"items"`
	CollectionFormat string `json:"collectionFormat"`
}

type Router struct {
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Tags        []string    `json:"tags"`
	Comsumes    []string    `json:"consumes"`
	Produces    []string    `json:"produces"`
	Parameters  []Parameter `json:"parameters"`
	Responses   map[string]struct {
		Description string `json:"description"`
		Schema      struct {
			Type  string            `json:"type"`
			Items map[string]string `json:"items"`
		} `json:"schema"`
	} `json:"responses"`
	Security []map[string][]string `json:"security"`
}
