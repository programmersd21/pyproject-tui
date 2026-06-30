package parser

// PyProject is the top-level structure of a pyproject.toml file.
// Raw holds the entire parsed file as a map[string]any for generic access.
// Known top-level sections are parsed into typed structs.
type PyProject struct {
	Project     *ProjectSection
	BuildSystem *BuildSystemSection
	Tool        map[string]any
	Raw         map[string]any
	FilePath    string
}

// ProjectSection represents the [project] table.
type ProjectSection struct {
	Name                 string
	Version              string
	Description          string
	Readme               any
	RequiresPython       string
	License              any
	Authors              []AuthorEntry
	Maintainers          []AuthorEntry
	Keywords             []string
	Classifiers          []string
	URLs                 map[string]string
	Dependencies         []string
	OptionalDependencies map[string][]string
	Scripts              map[string]string
	EntryPoints          map[string]map[string]string
	Dynamic              []string
}

// AuthorEntry is a project author or maintainer entry.
type AuthorEntry struct {
	Name  string `toml:"name,omitempty"`
	Email string `toml:"email,omitempty"`
}

// BuildSystemSection represents the [build-system] table.
type BuildSystemSection struct {
	Requires     []string
	BuildBackend string
	BackendPath  []string
}

// SectionID identifies a navigable section in the sidebar.
type SectionID struct {
	Kind    SectionKind
	ToolKey string
}

// SectionKind identifies the section family.
type SectionKind int

const (
	// KindProject identifies the [project] section.
	KindProject SectionKind = iota
	// KindBuildSystem identifies the [build-system] section.
	KindBuildSystem
	// KindTool identifies a [tool.*] section.
	KindTool
)

// FieldType describes how a field value should be rendered and edited.
type FieldType int

const (
	// FieldTypeString renders string-like values.
	FieldTypeString FieldType = iota
	// FieldTypeStringArray renders []string values.
	FieldTypeStringArray
	// FieldTypeStringMap renders map[string]string values.
	FieldTypeStringMap
	// FieldTypeAuthorArray renders []AuthorEntry values.
	FieldTypeAuthorArray
	// FieldTypeStringMapOfArray renders map[string][]string values.
	FieldTypeStringMapOfArray
	// FieldTypeAny renders generic TOML data.
	FieldTypeAny
	// FieldTypeBool renders booleans.
	FieldTypeBool
	// FieldTypeInt renders integers.
	FieldTypeInt
)

// Field is a single key-value entry in the editor pane.
type Field struct {
	Key   string
	Value any
	Type  FieldType
	Dirty bool
}
