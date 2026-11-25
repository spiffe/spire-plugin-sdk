package pluginconfig

import (
	"testing"

	configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"
	"github.com/stretchr/testify/require"
)

type config struct {
	Name string `hcl:"name" yaml:"name"`
}

func TestText(t *testing.T) {
	for _, tt := range []struct {
		name       string
		req        *configv1.ConfigureRequest
		wantText   string
		wantFormat configv1.ConfigFormat
	}{
		{
			name:       "format not set",
			req:        &configv1.ConfigureRequest{HclConfiguration: `name = "a"`},
			wantText:   `name = "a"`,
			wantFormat: configv1.ConfigFormat_CONFIG_FORMAT_HCL,
		},
		{
			name: "HCL",
			req: &configv1.ConfigureRequest{
				HclConfiguration: `name = "a"`,
				Configuration:    `name = "a"`,
				ConfigFormat:     configv1.ConfigFormat_CONFIG_FORMAT_HCL,
			},
			wantText:   `name = "a"`,
			wantFormat: configv1.ConfigFormat_CONFIG_FORMAT_HCL,
		},
		{
			name: "empty YAML",
			req: &configv1.ConfigureRequest{
				HclConfiguration: `name = "a"`,
				ConfigFormat:     configv1.ConfigFormat_CONFIG_FORMAT_YAML,
			},
			wantFormat: configv1.ConfigFormat_CONFIG_FORMAT_YAML,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			text, format := Text(tt.req)
			require.Equal(t, tt.wantText, text)
			require.Equal(t, tt.wantFormat, format)
		})
	}
}

func TestDecode(t *testing.T) {
	for _, tt := range []struct {
		name    string
		req     *configv1.ValidateRequest
		want    config
		wantErr string
	}{
		{
			name: "format not set",
			req:  &configv1.ValidateRequest{HclConfiguration: `name = "a"`},
			want: config{Name: "a"},
		},
		{
			name: "HCL",
			req:  &configv1.ValidateRequest{Configuration: `name = "a"`, ConfigFormat: configv1.ConfigFormat_CONFIG_FORMAT_HCL},
			want: config{Name: "a"},
		},
		{
			name: "YAML",
			req:  &configv1.ValidateRequest{Configuration: "name: a\n", ConfigFormat: configv1.ConfigFormat_CONFIG_FORMAT_YAML},
			want: config{Name: "a"},
		},
		{
			name: "empty YAML",
			req:  &configv1.ValidateRequest{ConfigFormat: configv1.ConfigFormat_CONFIG_FORMAT_YAML},
		},
		{
			name:    "unknown YAML key",
			req:     &configv1.ValidateRequest{Configuration: "nme: a\n", ConfigFormat: configv1.ConfigFormat_CONFIG_FORMAT_YAML},
			wantErr: "field nme not found",
		},
		{
			name:    "unknown format",
			req:     &configv1.ValidateRequest{ConfigFormat: 99},
			wantErr: "unsupported configuration format 99",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got config
			err := Decode(tt.req, &got)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
