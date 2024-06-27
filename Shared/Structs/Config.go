package Structs

type DatabaseRow struct {
	RowId        string `yaml:"row_id"`
	DataType     string `yaml:"data_type"`
	IsPrimaryKey bool   `yaml:"primary_key"`
	IsNull       bool   `yaml:"is_null"`
}

type TableDef struct {
	TableName string        `yaml:"table_name"`
	TableRows []DatabaseRow `yaml:"table_rows"`
}
type Config struct {
	Server struct {
		Port    string `yaml:"port"`
		Host    string `yaml:"host"`
		SSLCert string `yaml:"ssl_cert"`
		SSLKey  string `yaml:"ssl_key"`
	} `yaml:"server"`
	Database struct {
		TableDef []TableDef `yaml:"database_tables"`
		HomeDir  string     `yaml:"home_dir"`
		FileName string     `yaml:"file_name"`
	} `yaml:"database"`
}
