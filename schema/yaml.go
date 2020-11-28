package schema

import (
	"io"

	"github.com/goccy/go-yaml"
)

// MarshalYAML return custom YAML byte
func (c Column) MarshalYAML() ([]byte, error) {
	if c.Default.Valid {
		return yaml.Marshal(&struct {
			Name            string      `yaml:"name"`
			Type            string      `yaml:"type"`
			Nullable        bool        `yaml:"nullable"`
			Default         string      `yaml:"default"`
			Comment         string      `yaml:"comment"`
			ParentRelations []*Relation `yaml:"-"`
			ChildRelations  []*Relation `yaml:"-"`
		}{
			Name:            c.Name,
			Type:            c.Type,
			Nullable:        c.Nullable,
			Default:         c.Default.String,
			Comment:         c.Comment,
			ParentRelations: c.ParentRelations,
			ChildRelations:  c.ChildRelations,
		})
	}
	return yaml.Marshal(&struct {
		Name            string      `yaml:"name"`
		Type            string      `yaml:"type"`
		Nullable        bool        `yaml:"nullable"`
		Default         *string     `yaml:"default"`
		Comment         string      `yaml:"comment"`
		ParentRelations []*Relation `yaml:"-"`
		ChildRelations  []*Relation `yaml:"-"`
	}{
		Name:            c.Name,
		Type:            c.Type,
		Nullable:        c.Nullable,
		Default:         nil,
		Comment:         c.Comment,
		ParentRelations: c.ParentRelations,
		ChildRelations:  c.ChildRelations,
	})
}

// MarshalYAML return custom YAML byte
func (r Relation) MarshalYAML() ([]byte, error) {
	columns := []string{}
	parentColumns := []string{}
	for _, c := range r.Columns {
		columns = append(columns, c.Name)
	}
	for _, c := range r.ParentColumns {
		parentColumns = append(parentColumns, c.Name)
	}

	return yaml.Marshal(&struct {
		Table         string   `yaml:"table"`
		Columns       []string `yaml:"columns"`
		ParentTable   string   `yaml:"parentTable"`
		ParentColumns []string `yaml:"parentColumns"`
		Def           string   `yaml:"def"`
		Virtual       bool     `yaml:"virtual"`
	}{
		Table:         r.Table.Name,
		Columns:       columns,
		ParentTable:   r.ParentTable.Name,
		ParentColumns: parentColumns,
		Def:           r.Def,
		Virtual:       r.Virtual,
	})
}

// UnmarshalYAML unmarshal YAML to schema.Column
func (c *Column) UnmarshalYAML(data []byte) error {
	s := struct {
		Name            string      `yaml:"name"`
		Type            string      `yaml:"type"`
		Nullable        bool        `yaml:"nullable"`
		Default         *string     `yaml:"default"`
		Comment         string      `yaml:"comment"`
		ParentRelations []*Relation `yaml:"-"`
		ChildRelations  []*Relation `yaml:"-"`
	}{}
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	c.Name = s.Name
	c.Type = s.Type
	c.Nullable = s.Nullable
	if s.Default != nil {
		c.Default.Valid = true
		c.Default.String = *s.Default
	} else {
		c.Default.Valid = false
		c.Default.String = ""
	}
	c.Comment = s.Comment
	return nil
}

// UnmarshalYAML unmarshal YAML to schema.Column
func (r *Relation) UnmarshalYAML(data []byte) error {
	s := struct {
		Table         string   `yaml:"table"`
		Columns       []string `yaml:"columns"`
		ParentTable   string   `yaml:"parentTable"`
		ParentColumns []string `yaml:"parentColumns"`
		Def           string   `yaml:"def"`
		Virtual       bool     `yaml:"virtual"`
	}{}
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	r.Table = &Table{
		Name: s.Table,
	}
	r.Columns = []*Column{}
	for _, c := range s.Columns {
		r.Columns = append(r.Columns, &Column{
			Name: c,
		})
	}
	r.ParentTable = &Table{
		Name: s.ParentTable,
	}
	r.ParentColumns = []*Column{}
	for _, c := range s.ParentColumns {
		r.ParentColumns = append(r.ParentColumns, &Column{
			Name: c,
		})
	}
	r.Def = s.Def
	r.Virtual = s.Virtual
	return nil
}

// YAML struct
type YAML struct{}

// OutputSchema output YAML format for full relation.
func (j *YAML) OutputSchema(wr io.Writer, s *Schema) error {
	encoder := yaml.NewEncoder(wr)
	err := encoder.Encode(s)
	if err != nil {
		return err
	}
	return nil
}

// OutputTable output YAML format for table.
func (j *YAML) OutputTable(wr io.Writer, t *Table) error {
	encoder := yaml.NewEncoder(wr)
	err := encoder.Encode(t)
	if err != nil {
		return err
	}
	return nil
}

func (s *Schema) SaveYaml(wr io.Writer) error {
	o := new(YAML)
	return o.OutputSchema(wr, s)
}
