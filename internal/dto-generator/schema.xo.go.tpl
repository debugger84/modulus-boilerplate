{{ define "enum" }}
{{- $e := .Data -}}
// {{ $e.GoName }} is the '{{ $e.SQLName }}' enum type from schema '{{ schema }}'.
type {{ $e.GoName }} uint16

// {{ $e.GoName }} values.
const (
{{ range $e.Values -}}
	// {{ $e.GoName }}{{ .GoName }} is the '{{ .SQLName }}' {{ $e.SQLName }}.
	{{ $e.GoName }}{{ .GoName }} {{ $e.GoName }} = {{ .ConstValue }}
{{ end -}}
)

// String satisfies the fmt.Stringer interface.
func ({{ short $e.GoName }} {{ $e.GoName }}) String() string {
	switch {{ short $e.GoName }} {
{{ range $e.Values -}}
	case {{ $e.GoName }}{{ .GoName }}:
		return "{{ .SQLName }}"
{{ end -}}
	}
	return fmt.Sprintf("{{ $e.GoName }}(%d)", {{ short $e.GoName }})
}

// MarshalText marshals {{ $e.GoName }} into text.
func ({{ short $e.GoName }} {{ $e.GoName }}) MarshalText() ([]byte, error) {
	return []byte({{ short $e.GoName }}.String()), nil
}

// UnmarshalText unmarshals {{ $e.GoName }} from text.
func ({{ short $e.GoName }} *{{ $e.GoName }}) UnmarshalText(buf []byte) error {
	switch str := string(buf); str {
{{ range $e.Values -}}
	case "{{ .SQLName }}":
		*{{ short $e.GoName }} = {{ $e.GoName }}{{ .GoName }}
{{ end -}}
	default:
		return ErrInvalid{{ $e.GoName }}(str)
	}
	return nil
}

// Value satisfies the driver.Valuer interface.
func ({{ short $e.GoName }} {{ $e.GoName }}) Value() (driver.Value, error) {
	return {{ short $e.GoName }}.String(), nil
}

// Scan satisfies the sql.Scanner interface.
func ({{ short $e.GoName }} *{{ $e.GoName }}) Scan(v interface{}) error {
	if buf, ok := v.([]byte); ok {
		return {{ short $e.GoName }}.UnmarshalText(buf)
	}
	return ErrInvalid{{ $e.GoName }}(fmt.Sprintf("%T", v))
}

{{ $nullName := (printf "%s%s" "Null" $e.GoName) -}}
{{- $nullShort := (short $nullName) -}}
// {{ $nullName }} represents a null '{{ $e.SQLName }}' enum for schema '{{ schema }}'.
type {{ $nullName }} struct {
	{{ $e.GoName }} {{ $e.GoName }}
	// Valid is true if {{ $e.GoName }} is not null.
	Valid bool
}

// Value satisfies the driver.Valuer interface.
func ({{ $nullShort }} {{ $nullName }}) Value() (driver.Value, error) {
	if !{{ $nullShort }}.Valid {
		return nil, nil
	}
	return {{ $nullShort }}.{{ $e.GoName }}.Value()
}

// Scan satisfies the sql.Scanner interface.
func ({{ $nullShort }} *{{ $nullName }}) Scan(v interface{}) error {
	if v == nil {
		{{ $nullShort }}.{{ $e.GoName }}, {{ $nullShort }}.Valid = 0, false
		return nil
	}
	err := {{ $nullShort }}.{{ $e.GoName }}.Scan(v)
	{{ $nullShort }}.Valid = err == nil
	return err
}

// ErrInvalid{{ $e.GoName }} is the invalid {{ $e.GoName }} error.
type ErrInvalid{{ $e.GoName }} string

// Error satisfies the error interface.
func (err ErrInvalid{{ $e.GoName }}) Error() string {
	return fmt.Sprintf("invalid {{ $e.GoName }}(%s)", string(err))
}
{{ end }}

{{ define "foreignkey" }}
{{- $k := .Data -}}
// {{ func_name_context $k }} returns the {{ $k.RefTable }} associated with the {{ $k.Table.GoName }}'s ({{ names "" $k.Fields }}).
//
// Generated from foreign key '{{ $k.SQLName }}'.
{{ recv_context $k.Table $k }} {
	return {{ foreign_key_context $k }}
}
{{- if context_both }}

// {{ func_name $k }} returns the {{ $k.RefTable }} associated with the {{ $k.Table }}'s ({{ names "" $k.Fields }}).
//
// Generated from foreign key '{{ $k.SQLName }}'.
{{ recv $k.Table $k }} {
	return {{ foreign_key $k }}
}
{{- end }}
{{ end }}


{{ define "procs" }}
{{- $ps := .Data -}}
{{- range $p := $ps -}}
// {{ func_name_context $p }} calls the stored {{ $p.Type }} '{{ $p.Signature }}' on db.
{{ func_context $p }} {
{{- if and (driver "mysql") (eq $p.Type "procedure") (not $p.Void) }}
	// At the moment, the Go MySQL driver does not support stored procedures
	// with out parameters
	return {{ zero $p.Returns }}, fmt.Errorf("unsupported")
{{- else }}
	// call {{ schema $p.SQLName }}
	{{ sqlstr "proc" $p }}
	// run
{{- if not $p.Void }}
{{- range $p.Returns }}
	var {{ check_name .GoName }} {{ type .Type }}
{{- end }}
	logf(sqlstr, {{ params $p.Params false }})
{{- if and (driver "sqlserver" "oracle") (eq $p.Type "procedure")}}
	if _, err := {{ db_named "Exec" $p }}; err != nil {
{{- else }}
	if err := {{ db "QueryRow" $p }}.Scan({{ names "&" $p.Returns }}); err != nil {
{{- end }}
		return {{ zero $p.Returns }}, logerror(err)
	}
	return {{ range $p.Returns }}{{ check_name .GoName }}, {{ end }}nil
{{- else }}
	logf(sqlstr)
{{- if driver "sqlserver" "oracle" }}
	if _, err := {{ db_named "Exec" $p }}; err != nil {
{{- else }}
	if _, err := {{ db "Exec" $p }}; err != nil {
{{- end }}
		return logerror(err)
	}
	return nil
{{- end }}
{{- end }}
}

{{ if context_both -}}
// {{ func_name $p }} calls the {{ $p.Type }} '{{ $p.Signature }}' on db.
{{ func $p }} {
	return {{ func_name_context $p }}({{ names_all "" "context.Background()" "db" $p.Params }})
}
{{- end -}}
{{- end }}
{{ end }}

{{ define "typedef" }}
{{- $t := .Data -}}
{{- if $t.Comment -}}
// {{ $t.Comment | eval $t.GoName }}
{{- else -}}
// {{ $t.GoName }} represents a row from '{{ schema $t.SQLName }}'.
{{- end }}
type {{ $t.GoName }} struct {
{{ range $t.Fields -}}
	{{ field . $t.GoName }}
{{ end }}
}

func ({{ $t.GoName }}) TableName() string {
	return {{ $t.GoName }}Table
}
type {{ $t.GoName }}Finder struct {
    db *gorm.DB
}

func New{{ $t.GoName }}Finder(gormDb *gorm.DB) *UserFinder {
	return &UserFinder{db: gormDb}
}

func (f *{{ $t.GoName }}Finder) OneByQuery(query QueryBuilder) *{{ $t.GoName }} {
	var user *{{ $t.GoName }}
	query.Build().Limit(1).Scan(&user)

	return user
}

func (f *{{ $t.GoName }}Finder) ListByQuery(query QueryBuilder, count int) []*{{ $t.GoName }} {
	var users []*{{ $t.GoName }}
	query.Build().Limit(count).Scan(&users)

	return users
}

func (f *{{ $t.GoName }}Finder) CreateQuery(ctx context.Context) *{{ $t.GoName }}Query {
	return New{{ $t.GoName }}Query(ctx, f.db)
}

type {{ $t.GoName }}Saver struct {
	db *gorm.DB
}

func New{{ $t.GoName }}Saver(db *gorm.DB) *{{ $t.GoName }}Saver {
	return &UserSaver{db: db}
}

func (f *{{ $t.GoName }}Saver) Create(ctx context.Context, entity {{ $t.GoName }}) error {
	result := f.db.Table({{ $t.GoName }}Table).WithContext(ctx).Create(&entity)

	return result.Error
}

func (f *UserSaver) Update(ctx context.Context, entity {{ $t.GoName }}) error {
	result := f.db.Table({{ $t.GoName }}Table).WithContext(ctx).Save(&entity)

	return result.Error
}


const {{ $t.GoName }}Table = `"{{ schema }}"."{{ $t.SQLName }}"`

type {{ $t.GoName }}Query struct {
	db *gorm.DB
}

func New{{ $t.GoName }}Query(ctx context.Context, db *gorm.DB) *UserQuery {
	localCopy := db.Table({{ $t.GoName }}Table).WithContext(ctx)
	query := &{{ $t.GoName }}Query{
		db: localCopy,
	}
	return query
}

{{ range $t.Fields -}}
    {{- if hasEqQuery .Type }}
func (p *{{ $t.GoName }}Query) {{.GoName}}({{.GoName}} {{.Type}}) *{{ $t.GoName }}Query {
	p.db = p.db.Where({{ $t.GoName }}Table+".{{.SQLName}} = ?", {{.GoName}})

	return p
}
    {{ end -}}

    {{- if hasLikeQuery .Type }}
    func (p *{{ $t.GoName }}Query) {{.GoName}}Like(pattern string) *{{ $t.GoName }}Query {
    	p.db = p.db.Where({{ $t.GoName }}Table+".{{.SQLName}} ilike ?", pattern)

    	return p
    }
    {{ end -}}

    {{- if hasInQuery .Type }}
        func (p *{{ $t.GoName }}Query) {{.GoName}}In(items []{{.Type}}) *{{ $t.GoName }}Query {
        	p.db = p.db.Where({{ $t.GoName }}Table+".{{.SQLName}} IN (?)", items)

        	return p
        }
    {{ end -}}
{{ end -}}

func (p *{{ $t.GoName }}Query) Build() *gorm.DB {
	return p.db
}

{{ end }}


{{ define "index" }}
{{end}}