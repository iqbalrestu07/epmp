package config

// FieldType is the set of supported domain field types.
type FieldType string

const (
	FieldTypeUUID      FieldType = "uuid"
	FieldTypeString    FieldType = "string"
	FieldTypeInt       FieldType = "int"
	FieldTypeBool      FieldType = "bool"
	FieldTypeText      FieldType = "text"
	FieldTypeEnum      FieldType = "enum"
	FieldTypeTimestamp FieldType = "timestamp"
	FieldTypeDecimal   FieldType = "decimal"
)

// Operation is a REST API operation.
type Operation string

const (
	OperationCreate Operation = "create"
	OperationRead   Operation = "read"
	OperationUpdate Operation = "update"
	OperationDelete Operation = "delete"
	OperationList   Operation = "list"
)

// ArtifactType is a strongly-typed code-generation target.
type ArtifactType string

const (
	// Backend artifacts
	ArtifactDTO        ArtifactType = "dto"
	ArtifactRepository ArtifactType = "repository"
	ArtifactUseCase    ArtifactType = "usecase"
	ArtifactCRUD       ArtifactType = "crud"
	ArtifactREST       ArtifactType = "rest"
	ArtifactMigration  ArtifactType = "migration"
	ArtifactTest       ArtifactType = "test"
	ArtifactOpenAPI    ArtifactType = "openapi"

	// Frontend artifacts
	ArtifactFETypes      ArtifactType = "fe_types"
	ArtifactFESchema     ArtifactType = "fe_schema"
	ArtifactFEApi        ArtifactType = "fe_api"
	ArtifactFEHooks      ArtifactType = "fe_hooks"
	ArtifactFEComponents ArtifactType = "fe_components"
	ArtifactFEPages      ArtifactType = "fe_pages"
)

// GeneratorConfig is the root configuration object loaded from YAML.
// It is a unified configuration shared between both Backend and Frontend generators.
type GeneratorConfig struct {
	Version  string         `yaml:"version"`
	DryRun   bool           `yaml:"dry_run"`
	Backend  BackendConfig  `yaml:"backend"`
	Frontend FrontendConfig `yaml:"frontend"`
	Domain   DomainSpec     `yaml:"domain"`
}

// BackendConfig holds backend generation settings.
type BackendConfig struct {
	OutputRoot string         `yaml:"output_root"`
	ModulePath string         `yaml:"module_path"`
	Artifacts  []ArtifactType `yaml:"artifacts"`
}

// FrontendConfig holds frontend generation settings.
type FrontendConfig struct {
	OutputRoot string         `yaml:"output_root"`
	Artifacts  []ArtifactType `yaml:"artifacts"`
}

// DomainSpec describes one bounded-context entity to generate.
type DomainSpec struct {
	Name           string       `yaml:"name"`
	Package        string       `yaml:"package"`
	Table          string       `yaml:"table"`
	BoundedContext string       `yaml:"bounded_context"`
	Fields         []FieldSpec  `yaml:"fields"`
	Behaviors      BehaviorSpec `yaml:"behaviors"`
	REST           RESTSpec     `yaml:"rest"`
	Events         []string     `yaml:"events"`
}

// FieldSpec describes a single domain field.
type FieldSpec struct {
	Name       string      `yaml:"name"`
	Type       FieldType   `yaml:"type"`
	Nullable   bool        `yaml:"nullable"`
	PrimaryKey bool        `yaml:"primary_key"`
	Unique     bool        `yaml:"unique"`
	MaxLength  int         `yaml:"max_length"`
	Default    interface{} `yaml:"default"`
	Searchable bool        `yaml:"searchable"`
	SoftDelete bool        `yaml:"soft_delete"`
	EnumValues []string    `yaml:"enum_values"`
	ForeignKey *ForeignKey `yaml:"foreign_key"`
	Auto       bool        `yaml:"auto"`
}

// ForeignKey describes a reference to another domain entity.
type ForeignKey struct {
	Domain string `yaml:"domain"`
	Field  string `yaml:"field"`
}

// BehaviorSpec controls cross-cutting behaviors.
type BehaviorSpec struct {
	SoftDelete bool     `yaml:"soft_delete"`
	AuditTrail bool     `yaml:"audit_trail"`
	Pagination bool     `yaml:"pagination"`
	Search     bool     `yaml:"search"`
	FilterBy   []string `yaml:"filter_by"`
	SortBy     []string `yaml:"sort_by"`
}

// RESTSpec controls the REST API surface.
type RESTSpec struct {
	BasePath     string      `yaml:"base_path"`
	AuthRequired bool        `yaml:"auth_required"`
	Operations   []Operation `yaml:"operations"`
}
