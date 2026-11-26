// Package pluginconfig reads the plugin configuration from Configure and
// Validate requests, which SPIRE sends as HCL or YAML.
package pluginconfig

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl"
	configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"
	"go.yaml.in/yaml/v3"
)

// Request is implemented by configv1.ConfigureRequest and
// configv1.ValidateRequest.
type Request interface {
	GetHclConfiguration() string
	GetConfiguration() string
	GetConfigFormat() configv1.ConfigFormat
}

// Text returns the plugin configuration of req and its format, which is never
// CONFIG_FORMAT_UNSPECIFIED.
func Text(req Request) (string, configv1.ConfigFormat) {
	format := req.GetConfigFormat()
	if format == configv1.ConfigFormat_CONFIG_FORMAT_UNSPECIFIED {
		return req.GetHclConfiguration(), configv1.ConfigFormat_CONFIG_FORMAT_HCL
	}
	return req.GetConfiguration(), format
}

// Decode decodes the plugin configuration of req into out, using its hcl or
// yaml struct tags. Unknown YAML keys, null map and list entries, and
// multiple YAML documents are rejected.
func Decode(req Request, out any) error {
	if v := reflect.ValueOf(out); v.Kind() != reflect.Pointer || v.IsNil() {
		return errors.New("result must be a non-nil pointer")
	}
	text, format := Text(req)
	switch format {
	case configv1.ConfigFormat_CONFIG_FORMAT_HCL:
		return hcl.Decode(out, text)
	case configv1.ConfigFormat_CONFIG_FORMAT_YAML:
		dec := yaml.NewDecoder(strings.NewReader(text))
		dec.KnownFields(true)
		if err := dec.Decode(out); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if err := dec.Decode(new(yaml.Node)); !errors.Is(err, io.EOF) {
			return errors.New("multiple YAML documents are not supported")
		}
		return checkNullEntries(reflect.ValueOf(out), "")
	default:
		return fmt.Errorf("unsupported configuration format %s", format)
	}
}

// checkNullEntries rejects map and list entries left nil by a null value,
// such as "prod" in "clusters: {prod: }".
func checkNullEntries(v reflect.Value, path string) error {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !v.IsNil() {
			return checkNullEntries(v.Elem(), path)
		}
	case reflect.Struct:
		t := v.Type()
		for i := range t.NumField() {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			name, opts, _ := strings.Cut(f.Tag.Get("yaml"), ",")
			if name == "-" {
				continue
			}
			fieldPath := path
			if opts != "inline" {
				if name == "" {
					name = strings.ToLower(f.Name)
				}
				fieldPath = joinPath(path, name)
			}
			if err := checkNullEntries(v.Field(i), fieldPath); err != nil {
				return err
			}
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			if err := checkEntry(iter.Value(), joinPath(path, fmt.Sprint(iter.Key()))); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := range v.Len() {
			if err := checkEntry(v.Index(i), fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkEntry(v reflect.Value, path string) error {
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return fmt.Errorf("%s must not be null; use {} for an entry without options", path)
	}
	return checkNullEntries(v, path)
}

func joinPath(path, name string) string {
	if path == "" {
		return name
	}
	return path + "." + name
}
