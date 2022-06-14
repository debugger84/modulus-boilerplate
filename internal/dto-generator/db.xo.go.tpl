{{ define "db" -}}
{{ if driver "postgres" -}}

type QueryBuilder interface {
	Build() *gorm.DB
}

type JSON json.RawMessage

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	b, err := json.RawMessage(j).MarshalJSON()
	return string(b), err
}

// StringSlice is a slice of strings.
type StringSlice []string

// Scan satisfies the sql.Scanner interface for StringSlice.
func (ss *StringSlice) Scan(v interface{}) error {
	buf, ok := v.([]byte)
	if !ok {
		return errors.New("wrong slice conversion")
	}
	// change quote escapes for csv parser
	str := strings.Replace(quoteEscRE.ReplaceAllString(string(buf), `$1""`), `\\`, `\`, -1)
	str = str[1 : len(str)-1]
	// bail if only one
	if len(str) == 0 {
		return nil
	}
	// parse with csv reader
	r := csv.NewReader(strings.NewReader(str))
	line, err := r.Read()
	if err != nil {
		return errors.New("wrong slice reading")
	}
	*ss = StringSlice(line)
	return nil
}

// quoteEscRE matches escaped characters in a string.
var quoteEscRE = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

// Value satisfies the sql/driver.Valuer interface.
func (ss StringSlice) Value() (driver.Value, error) {
	v := make([]string, len(ss))
	for i, s := range ss {
		v[i] = `"` + strings.Replace(strings.Replace(s, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(v, ",") + "}", nil
}
{{- end }}

{{ if driver "sqlite3" -}}

// Time is a SQLite3 Time that scans for the various timestamps values used by
// SQLite3 database drivers to store time.Time values.
type Time struct {
	time time.Time
}

// NewTime creates a time.
func NewTime(t time.Time) Time {
	return Time{time: t}
}

// String satisfies the fmt.Stringer interface.
func (t Time) String() string {
	return t.time.String()
}

// Format formats the time.
func (t Time) Format(layout string) string {
	return t.time.Format(layout)
}

// Time returns a time.Time.
func (t Time) Time() time.Time {
	return t.time
}

// Value satisfies the sql/driver.Valuer interface.
func (t Time) Value() (driver.Value, error) {
	return t.time, nil
}

// Scan satisfies the sql.Scanner interface.
func (t *Time) Scan(v interface{}) error {
	switch x := v.(type) {
	case time.Time:
		t.time = x
		return nil
	case []byte:
		return t.Parse(string(x))
	case string:
		return t.Parse(x)
	}
	return ErrInvalidTime(fmt.Sprintf("%T", v))
}

// Parse attempts to Parse string s to t.
func (t *Time) Parse(s string) error {
	if s == "" {
		return nil
	}
	for _, f := range TimestampFormats {
		if z, err := time.Parse(f, s); err == nil {
			t.time = z
			return nil
		}
	}
	return ErrInvalidTime(s)
}

// TimestampFormats are the timestamp formats used by SQLite3 database drivers
// to store a time.Time in SQLite3.
//
// The first format in the slice will be used when saving time values into the
// database.  When parsing a string from a timestamp or datetime column, the
// formats are tried in order.
var TimestampFormats = []string{
	// By default, use timestamps with the timezone they have. When parsed,
	// they will be returned with the same timezone.
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02T15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04",
	"2006-01-02T15:04",
	"2006-01-02",
}
{{- end }}
{{- end }}
