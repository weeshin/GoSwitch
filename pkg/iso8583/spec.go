package iso8583

import (
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLSpec matches the structure of your .yaml file
type YAMLSpec struct {
	Fields map[int]struct {
		Length      int    `yaml:"length"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
		Format      string `yaml:"format"`
	} `yaml:"fields"`
}

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
		var formatter Formatter

		// Map string format names to our Interface implementations
		switch f.Format {
		case "FIXED":
			formatter = &FixedFormatter{}
		case "LLVAR":
			formatter = &LLVarFormatter{}
		case "LLLVAR":
			formatter = &LLLVarFormatter{}
		default:
			formatter = &FixedFormatter{}
		}

		spec[id] = FieldSpec{
			Length:      f.Length,
			Description: f.Description,
			Type:        f.Type,
			Formatter:   formatter,
		}
	}

	return spec, nil
}

// FieldSpec defines how a specific field should be packed/unpacked
type FieldSpec struct {
	Length      int // Max length for variable, exact length for fixed
	Description string
	Type        string    // e.g., "N" (Numeric), "A" (Alpha), "AN" (Alphanumeric)
	Formatter   Formatter // e.g., "FIXED", "LLVAR", "LLLVAR"
}

// Spec is a map of field numbers to their definitions
type Spec map[int]FieldSpec

// GetDefaultSpec returns a basic ISO 8583 specification
func GetDefaultSpec() Spec {
	return Spec{
		// Field 1 is usually the Secondary Bitmap, handled by the engine
		2:  {Length: 19, Description: "Primary Account Number (PAN)", Type: "N", Formatter: &LLVarFormatter{}},
		3:  {Length: 6, Description: "Processing Code", Type: "N", Formatter: &FixedFormatter{}},
		4:  {Length: 12, Description: "Amount, Transaction", Type: "N", Formatter: &FixedFormatter{}},
		7:  {Length: 10, Description: "Transmission Date & Time", Type: "N", Formatter: &FixedFormatter{}},
		11: {Length: 6, Description: "System Trace Audit Number (STAN)", Type: "N", Formatter: &FixedFormatter{}},
		12: {Length: 6, Description: "Time, Local Transaction", Type: "N", Formatter: &FixedFormatter{}},
		13: {Length: 4, Description: "Date, Local Transaction", Type: "N", Formatter: &FixedFormatter{}},
		37: {Length: 12, Description: "Retrieval Reference Number", Type: "AN", Formatter: &FixedFormatter{}},
		41: {Length: 8, Description: "Card Acceptor Terminal Identification", Type: "ANS", Formatter: &FixedFormatter{}},
		70: {Length: 3, Description: "Network Management Information Code", Type: "N", Formatter: &FixedFormatter{}},
	}
}
