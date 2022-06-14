package dto_generator

type OverrodeType struct {
	GoType string
	DbType string
}

var types = map[string]map[string]OverrodeType{
	"User": map[string]OverrodeType{
		"Settings": {
			GoType: "*Settings",
			DbType: "",
		},
		"Contacts": {
			GoType: "pq.StringArray",
			DbType: "text[]",
		},
	},
}
