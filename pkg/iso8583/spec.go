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

// Spec is a map of field numbers to their definitions
type Spec map[int]FieldSpec

// LoadSpecFromFile reads a YAML file and returns a usable Spec
func LoadSpecFromFile(path string) (Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var y YAMLSpec
	if err := yaml.Unmarshal(data, &y); err != nil {
		return nil, err
	}

	spec := make(Spec)
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

		spec[id] = FieldSpec{
			Length:      f.Length,
			Description: f.Description,
			Encoder:     encoder,
		}
	}

	return spec, nil
}

// GetDefaultSpec provides a hardcoded fallback using the new field package
func GetDefaultSpec() Spec {
	return Spec{
		2:  {Length: 19, Description: "PAN", Encoder: &field.FANumeric{}}, // Assuming ASCII for default
		3:  {Length: 6, Description: "Proc Code", Encoder: &field.FANumeric{}},
		11: {Length: 6, Description: "STAN", Encoder: &field.FBNumeric{}}, // Example of Binary
		70: {Length: 3, Description: "Net Code", Encoder: &field.FANumeric{}},
	}
}
