package iso8583

import (
	"GoSwitch/pkg/field"
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLSpec matches the structure of your .yaml file
type YAMLSpec struct {
	Fields map[int]struct {
		Length      int    `yaml:"length"`
		Description string `yaml:"description"`
		Encoder     string `yaml:"encoder"` // e.g., "FANumeric", "FBNumeric"
	} `yaml:"fields"`
}

// FieldSpec defines how a specific field should be packed/unpacked
type FieldSpec struct {
	Length      int
	Description string
	Encoder     field.ISOField
}

// Spec defines the configuration for the ISO8583 message
type Spec struct {
	MTIEncoder    field.ISOField
	BitmapEncoder field.BitMap
	Fields        map[int]FieldSpec
}

// LoadSpecFromFile reads a YAML file and returns a usable Spec
func LoadSpecFromFile(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var y YAMLSpec
	if err := yaml.Unmarshal(data, &y); err != nil {
		return nil, err
	}

	spec := &Spec{
		Fields:        make(map[int]FieldSpec),
		MTIEncoder:    &field.FANumeric{},
		BitmapEncoder: &field.FBBitmap{},
	}

	for id, f := range y.Fields {
		var encoder field.ISOField

		switch f.Encoder {
		case "FBNumeric":
			encoder = &field.FBNumeric{}
		case "FANumeric":
			encoder = &field.FANumeric{}
		default:
			// Fallback to ASCII Numeric
			encoder = &field.FANumeric{}
		}

		spec.Fields[id] = FieldSpec{
			Length:      f.Length,
			Description: f.Description,
			Encoder:     encoder,
		}
	}

	return spec, nil
}
