package bundleformat

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"testing"

	"github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestBytes(t *testing.T) {
	const (
		certPEM = `-----BEGIN CERTIFICATE-----
MIIBKjCB0aADAgECAgEBMAoGCCqGSM49BAMCMAAwIhgPMDAwMTAxMDEwMDAwMDBa
GA85OTk5MTIzMTIzNTk1OVowADBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHyv
sCk5yi+yhSzNu5aquQwvm8a1Wh+qw1fiHAkhDni+wq+g3TQWxYlV51TCPH030yXs
RxvujD4hUUaIQrXk4KKjODA2MA8GA1UdEwEB/wQFMAMBAf8wIwYDVR0RAQH/BBkw
F4YVc3BpZmZlOi8vZG9tYWluMS50ZXN0MAoGCCqGSM49BAMCA0gAMEUCIA2dO09X
makw2ekuHKWC4hBhCkpr5qY4bI8YUcXfxg/1AiEA67kMyH7bQnr7OVLUrL+b9ylA
dZglS5kKnYigmwDh+/U=
-----END CERTIFICATE-----
`
	)
	block, _ := pem.Decode([]byte(certPEM))
	require.NotNil(t, block, "unable to unmarshal certificate response: malformed PEM block")

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	keyPkix, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	require.NoError(t, err)

	testBundle := &types.Bundle{
		TrustDomain:     "example.org",
		X509Authorities: []*types.X509Certificate{{Asn1: cert.Raw}},
		JwtAuthorities: []*types.JWTKey{
			{
				KeyId:     "KID",
				PublicKey: keyPkix,
			},
		},
		RefreshHint:    1440,
		SequenceNumber: 100,
	}
	standardJWKS := `{"keys":[{%s"kty":"EC","crv":"P-256","x":"fK-wKTnKL7KFLM27lqq5DC-bxrVaH6rDV-IcCSEOeL4","y":"wq-g3TQWxYlV51TCPH030yXsRxvujD4hUUaIQrXk4KI","x5c":["MIIBKjCB0aADAgECAgEBMAoGCCqGSM49BAMCMAAwIhgPMDAwMTAxMDEwMDAwMDBaGA85OTk5MTIzMTIzNTk1OVowADBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHyvsCk5yi+yhSzNu5aquQwvm8a1Wh+qw1fiHAkhDni+wq+g3TQWxYlV51TCPH030yXsRxvujD4hUUaIQrXk4KKjODA2MA8GA1UdEwEB/wQFMAMBAf8wIwYDVR0RAQH/BBkwF4YVc3BpZmZlOi8vZG9tYWluMS50ZXN0MAoGCCqGSM49BAMCA0gAMEUCIA2dO09Xmakw2ekuHKWC4hBhCkpr5qY4bI8YUcXfxg/1AiEA67kMyH7bQnr7OVLUrL+b9ylAdZglS5kKnYigmwDh+/U="]},{%s"kty":"EC","kid":"KID","crv":"P-256","x":"fK-wKTnKL7KFLM27lqq5DC-bxrVaH6rDV-IcCSEOeL4","y":"wq-g3TQWxYlV51TCPH030yXsRxvujD4hUUaIQrXk4KI"}]%s}`
	expectedJWKS := fmt.Sprintf(standardJWKS, "", "", "")
	expectedSPIFFEBundle := fmt.Sprintf(standardJWKS, `"use":"x509-svid",`, `"use":"jwt-svid",`, `,"spiffe_sequence":100,"spiffe_refresh_hint":1440`)
	expectedOIDCJWKS := `{"keys":[{"use":"sig","kty":"EC","kid":"KID","crv":"P-256","alg":"ES256","x":"fK-wKTnKL7KFLM27lqq5DC-bxrVaH6rDV-IcCSEOeL4","y":"wq-g3TQWxYlV51TCPH030yXsRxvujD4hUUaIQrXk4KI"}]}`

	for _, tt := range []struct {
		name        string
		format      Format
		bundle      *types.Bundle
		expectBytes []byte
		expectError string
	}{
		{
			name:        "format not set",
			bundle:      testBundle,
			expectError: "no format specified",
		},
		{
			name:        "invalid format",
			format:      math.MaxInt,
			bundle:      testBundle,
			expectError: fmt.Sprintf("invalid format: \"UNKNOWN(%d)\"", math.MaxInt),
		},
		{
			name:        "no bundle",
			format:      SPIFFE,
			expectError: "missing bundle",
		},
		{
			name:        "jwks format",
			format:      JWKS,
			bundle:      testBundle,
			expectBytes: []byte(expectedJWKS),
		},
		{
			name:        "oidc_jwks format",
			format:      OIDCJWKS,
			bundle:      testBundle,
			expectBytes: []byte(expectedOIDCJWKS),
		},
		{
			name:        "pem format",
			format:      PEM,
			bundle:      testBundle,
			expectBytes: []byte(certPEM),
		},
		{
			name:        "spiffe format",
			format:      SPIFFE,
			bundle:      testBundle,
			expectBytes: []byte(expectedSPIFFEBundle),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			b := NewFormatter(tt.bundle)

			if !proto.Equal(tt.bundle, b.bundle) {
				require.Equal(t, tt.bundle, b.bundle)
			}

			// Test the Format function that's provided by the formatter and
			// also test the FormatBundle function that should have the same
			// result.
			formatResult, formatErr := b.Format(tt.format)
			formatBundleResult, formatBundleErr := FormatBundle(tt.bundle, tt.format)
			if tt.expectError != "" {
				require.EqualError(t, formatErr, tt.expectError)
				require.Nil(t, formatResult)

				require.EqualError(t, formatBundleErr, tt.expectError)
				require.Nil(t, formatBundleResult)
				return
			}
			require.NoError(t, formatErr)
			require.NoError(t, formatBundleErr)

			require.Equal(t, string(tt.expectBytes), string(formatResult))
			require.Equal(t, string(tt.expectBytes), string(formatBundleResult))
		})
	}
}

func TestStringConversion(t *testing.T) {
	for _, tt := range []struct {
		name         string
		formatString string
		expectError  string
		expectFormat Format
	}{
		{
			name:         "invalid format",
			formatString: "INVALID",
			expectError:  `unknown bundle format: "INVALID"`,
		},
		{
			name:         "jwks format",
			formatString: "jwks",
			expectFormat: JWKS,
		},
		{
			name:         "oidc_jwks format",
			formatString: "oidc_jwks",
			expectFormat: OIDCJWKS,
		},
		{
			name:         "pem format",
			formatString: "pem",
			expectFormat: PEM,
		},
		{
			name:         "spiffe format",
			formatString: "spiffe",
			expectFormat: SPIFFE,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			format, err := FromString(tt.formatString)
			if tt.expectError != "" {
				require.EqualError(t, err, tt.expectError)
				require.Equal(t, FormatUnset, format)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectFormat, format)
			require.Equal(t, tt.formatString, format.String())
		})
	}
}

func TestOIDCJWKS(t *testing.T) {
	const (
		// Public keys whose JWS algorithm is known.
		rsa2048PEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnbr0TRmoINIvlIWOq1o2
flmf7PFmG4nRh0GdotR2Eo2FW1Wxx4mQ3qPrJq8zAtro+5+qOlIRkODsMusvrS9j
pjFoikbbrbqZoHHWNwWrxXiYtOpfUj+jFyv+tD2NY4cJcdU6tNLHQP4Fv22b5WMT
eqwLFZuvGrLrpaBcQGxJDoakpFs0FP9OYrZcbtPV9YF71fYUiw2FAN0zXPlLBci5
5BtXYWIDBJt4eGCFkkVnmUQJlrbrR/0gyzrztB0Ly3lOmQguNDkEmXUik7qKjv9Y
0uzBGIRhIW5htdcIZUW2vMRJHPMetlx07rHou9IGaJUBDzMYZxih7BpdeLlUglTW
5wIDAQAB
-----END PUBLIC KEY-----`
		ec256PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEIHOr9UXslszDD+f59Al99MFluwv/
6W7xJxC9OSQgsAE7kQL5igSoUx59pAwgrySQiTNd6fykghdHiSlW65bp1g==
-----END PUBLIC KEY-----`
		ec384PEM = `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEZQpnFVLKVJk/ftXilslYwLJeD9WD3NGM
EkkK5YmJ4c7ybAknvOZmduaDTyb8VCGl3vfdQDyFx6rdZlUSedU6gouqTMvRvnJr
1t9TIs5LEepPH9Xgls2Un2yKYeig1ipa
-----END PUBLIC KEY-----`

		// Public keys with no JWS algorithm.
		rsa1024PEM = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCqjONNfz4cTg31u5U4GSIDmy6H
cqjWEiLnfrQNAaVlXTFqY4DFOEp5Zk6lbPb/Q2bwCf56246PoZqMSJ+SWuIOGurN
AOWZM4L9dYf1m5sqGPHfIByPhPTQYgH2muCKT8kLBzNzsDiVxzRWXvfc3owRwlHt
Je5F+vs5Dzh2RZi+jQIDAQAB
-----END PUBLIC KEY-----`
		ec521PEM = `-----BEGIN PUBLIC KEY-----
MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQA05mB+oMFGPSjb2DYsnQmMMAx7yUO
Q5RJ+HHlULTMYDpXzODVkuuQJeE/FGhPVZE3Rjpq2ECRiEDEG1p8b6DIg8kBmspC
T9pdEVVNx9ZxAUYYkCd/XAe7WvfQl1VKDh1aWjR7RB+5nwlfYWiwkc/Q5ZVK+xif
TyNAgTnpZUeUW2r6FsM=
-----END PUBLIC KEY-----`

		rsa2048JWK = `{"use":"sig","kty":"RSA","kid":"rsa2048","alg":"RS256","n":"nbr0TRmoINIvlIWOq1o2flmf7PFmG4nRh0GdotR2Eo2FW1Wxx4mQ3qPrJq8zAtro-5-qOlIRkODsMusvrS9jpjFoikbbrbqZoHHWNwWrxXiYtOpfUj-jFyv-tD2NY4cJcdU6tNLHQP4Fv22b5WMTeqwLFZuvGrLrpaBcQGxJDoakpFs0FP9OYrZcbtPV9YF71fYUiw2FAN0zXPlLBci55BtXYWIDBJt4eGCFkkVnmUQJlrbrR_0gyzrztB0Ly3lOmQguNDkEmXUik7qKjv9Y0uzBGIRhIW5htdcIZUW2vMRJHPMetlx07rHou9IGaJUBDzMYZxih7BpdeLlUglTW5w","e":"AQAB"}`
		ec256JWK   = `{"use":"sig","kty":"EC","kid":"ec256","crv":"P-256","alg":"ES256","x":"IHOr9UXslszDD-f59Al99MFluwv_6W7xJxC9OSQgsAE","y":"O5EC-YoEqFMefaQMIK8kkIkzXen8pIIXR4kpVuuW6dY"}`
		ec384JWK   = `{"use":"sig","kty":"EC","kid":"ec384","crv":"P-384","alg":"ES384","x":"ZQpnFVLKVJk_ftXilslYwLJeD9WD3NGMEkkK5YmJ4c7ybAknvOZmduaDTyb8VCGl","y":"3vfdQDyFx6rdZlUSedU6gouqTMvRvnJr1t9TIs5LEepPH9Xgls2Un2yKYeig1ipa"}`
		rsa1024JWK = `{"use":"sig","kty":"RSA","kid":"rsa1024","n":"qozjTX8-HE4N9buVOBkiA5suh3Ko1hIi5360DQGlZV0xamOAxThKeWZOpWz2_0Nm8An-etuOj6GajEifklriDhrqzQDlmTOC_XWH9ZubKhjx3yAcj4T00GIB9prgik_JCwczc7A4lcc0Vl733N6MEcJR7SXuRfr7OQ84dkWYvo0","e":"AQAB"}`
		ec521JWK   = `{"use":"sig","kty":"EC","kid":"ec521","crv":"P-521","x":"ANOZgfqDBRj0o29g2LJ0JjDAMe8lDkOUSfhx5VC0zGA6V8zg1ZLrkCXhPxRoT1WRN0Y6athAkYhAxBtafG-gyIPJ","y":"AZrKQk_aXRFVTcfWcQFGGJAnf1wHu1r30JdVSg4dWlo0e0QfuZ8JX2FosJHP0OWVSvsYn08jQIE56WVHlFtq-hbD"}`

		// A self-signed X.509 authority.
		certPEM = `-----BEGIN CERTIFICATE-----
MIIBKjCB0aADAgECAgEBMAoGCCqGSM49BAMCMAAwIhgPMDAwMTAxMDEwMDAwMDBa
GA85OTk5MTIzMTIzNTk1OVowADBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHyv
sCk5yi+yhSzNu5aquQwvm8a1Wh+qw1fiHAkhDni+wq+g3TQWxYlV51TCPH030yXs
RxvujD4hUUaIQrXk4KKjODA2MA8GA1UdEwEB/wQFMAMBAf8wIwYDVR0RAQH/BBkw
F4YVc3BpZmZlOi8vZG9tYWluMS50ZXN0MAoGCCqGSM49BAMCA0gAMEUCIA2dO09X
makw2ekuHKWC4hBhCkpr5qY4bI8YUcXfxg/1AiEA67kMyH7bQnr7OVLUrL+b9ylA
dZglS5kKnYigmwDh+/U=
-----END CERTIFICATE-----`
	)

	pemBytes := func(s string) []byte {
		block, _ := pem.Decode([]byte(s))
		require.NotNil(t, block)
		return block.Bytes
	}
	jwtKey := func(keyID, s string) *types.JWTKey {
		return &types.JWTKey{KeyId: keyID, PublicKey: pemBytes(s)}
	}
	x509Authorities := []*types.X509Certificate{{Asn1: pemBytes(certPEM)}}

	for _, tt := range []struct {
		name        string
		bundle      *types.Bundle
		expectJWKS  string
		expectError string
	}{
		{
			name: "keys with an algorithm",
			bundle: &types.Bundle{
				TrustDomain: "example.org",
				JwtAuthorities: []*types.JWTKey{
					jwtKey("rsa2048", rsa2048PEM),
					jwtKey("ec256", ec256PEM),
					jwtKey("ec384", ec384PEM),
				},
			},
			expectJWKS: `{"keys":[` + rsa2048JWK + `,` + ec256JWK + `,` + ec384JWK + `]}`,
		},
		{
			name: "keys with no algorithm are included without alg",
			bundle: &types.Bundle{
				TrustDomain: "example.org",
				JwtAuthorities: []*types.JWTKey{
					jwtKey("rsa2048", rsa2048PEM),
					jwtKey("rsa1024", rsa1024PEM),
					jwtKey("ec256", ec256PEM),
					jwtKey("ec521", ec521PEM),
					jwtKey("ec384", ec384PEM),
				},
			},
			expectJWKS: `{"keys":[` + rsa2048JWK + `,` + rsa1024JWK + `,` + ec256JWK + `,` + ec521JWK + `,` + ec384JWK + `]}`,
		},
		{
			name: "X.509 authorities are not included",
			bundle: &types.Bundle{
				TrustDomain:     "example.org",
				X509Authorities: x509Authorities,
				JwtAuthorities: []*types.JWTKey{
					jwtKey("ec256", ec256PEM),
				},
				RefreshHint:    1440,
				SequenceNumber: 100,
			},
			expectJWKS: `{"keys":[` + ec256JWK + `]}`,
		},
		{
			name: "no JWT authorities",
			bundle: &types.Bundle{
				TrustDomain:     "example.org",
				X509Authorities: x509Authorities,
			},
			expectJWKS: `{"keys":[]}`,
		},
		{
			name: "empty bundle",
			bundle: &types.Bundle{
				TrustDomain: "example.org",
			},
			expectJWKS: `{"keys":[]}`,
		},
		{
			name: "invalid JWT authority",
			bundle: &types.Bundle{
				TrustDomain: "example.org",
				JwtAuthorities: []*types.JWTKey{
					jwtKey("ec256", ec256PEM),
					{KeyId: "invalid", PublicKey: []byte("invalid")},
				},
			},
			expectError: "could not convert bundle to oidc_jwks format: unable to parse JWT signing key 1: ",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			jwks, err := FormatBundle(tt.bundle, OIDCJWKS)
			if tt.expectError != "" {
				require.ErrorContains(t, err, tt.expectError)
				require.Nil(t, jwks)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectJWKS, string(jwks))
		})
	}
}
