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
// yaml struct tags. Unknown YAML keys and multiple YAML documents are
// rejected.
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
		return nil
	default:
		return fmt.Errorf("unsupported configuration format %s", format)
	}
}
