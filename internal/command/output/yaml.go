package output

import "github.com/goccy/go-yaml"

var _ Formatter = (*YAMLFormatter)(nil)

type YAMLFormatter struct {
	Indent int
}

func (f *YAMLFormatter) Format(data Formattable) (string, error) {
	bytes, err := yaml.MarshalWithOptions(data, yaml.Indent(f.Indent))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
