// Package bundleformat provides helper functions related with bundle formatting
// for plugins implementing the BundlePublisher interface.
// BundlePublisher plugins should use this package as a way to have a
// standarized name for bundle formats in their configuration, and avoid the
// re-implementation of bundle parsing logic of formats supported in this
// package.
package bundleformat

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/spiffe/go-spiffe/v2/bundle/spiffebundle"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/types"
)

const (
	FormatUnset Format = iota
	SPIFFE
	PEM
	JWKS
	// OIDCJWKS is a JWK Set for OIDC relying parties with the JWT authorities
	// only. Every key has its "kid", "use" set to "sig", and "alg" when the
	// JWS algorithm of the key is known: RS256 for RSA keys of at least 2048
	// bits, ES256 for P-256 and ES384 for P-384 ECDSA keys. Other keys have
	// no "alg".
	OIDCJWKS
)

// oidcKeyUse is the "use" parameter of the keys in the OIDCJWKS format.
const oidcKeyUse = "sig"

// Formatter formats a bundle in different formats.
type Formatter struct {
	bundle *types.Bundle

	jwksBytes     []byte
	oidcJWKSBytes []byte
	pemBytes      []byte
	spiffeBytes   []byte
}

// Format represents the bundle formats that are supported by the Formatter.
type Format int

// String returns the string name for the bundle format.
func (bundleFormat Format) String() string {
	switch bundleFormat {
	case FormatUnset:
		return "UNSET"
	case SPIFFE:
		return "spiffe"
	case PEM:
		return "pem"
	case JWKS:
		return "jwks"
	case OIDCJWKS:
		return "oidc_jwks"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(bundleFormat))
	}
}

// FromString returns the Format corresponding to the provided string.
func FromString(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "spiffe":
		return SPIFFE, nil
	case "jwks":
		return JWKS, nil
	case "pem":
		return PEM, nil
	case "oidc_jwks":
		return OIDCJWKS, nil
	default:
		return FormatUnset, fmt.Errorf("unknown bundle format: %q", s)
	}
}

// NewFormatter return a new *Formatter with the *types.Bundle provided.
// Use the Bytes() function to get a slice of bytes with the bundle formatted in
// the format specified.
func NewFormatter(pluginBundle *types.Bundle) *Formatter {
	return &Formatter{
		bundle: pluginBundle,
	}
}

// Format returns the bundle in the form of a slice of bytes in
// the chosen format.
func (b *Formatter) Format(format Format) ([]byte, error) {
	if b.bundle == nil {
		return nil, errors.New("missing bundle")
	}

	switch format {
	case FormatUnset:
		return nil, errors.New("no format specified")
	case JWKS:
		if b.jwksBytes != nil {
			return b.jwksBytes, nil
		}
		jwksBytes, err := b.toJWKS()
		if err != nil {
			return nil, fmt.Errorf("could not convert bundle to jwks format: %w", err)
		}
		b.jwksBytes = jwksBytes
		return jwksBytes, nil
	case OIDCJWKS:
		if b.oidcJWKSBytes != nil {
			return b.oidcJWKSBytes, nil
		}
		oidcJWKSBytes, err := b.toOIDCJWKS()
		if err != nil {
			return nil, fmt.Errorf("could not convert bundle to oidc_jwks format: %w", err)
		}
		b.oidcJWKSBytes = oidcJWKSBytes
		return oidcJWKSBytes, nil
	case PEM:
		if b.pemBytes != nil {
			return b.pemBytes, nil
		}
		pemBytes, err := b.toPEM()
		if err != nil {
			return nil, fmt.Errorf("could not convert bundle to pem format: %w", err)
		}
		b.pemBytes = pemBytes
		return pemBytes, nil
	case SPIFFE:
		if b.spiffeBytes != nil {
			return b.spiffeBytes, nil
		}
		spiffeBytes, err := b.toSPIFFEBundle()
		if err != nil {
			return nil, fmt.Errorf("could not convert bundle to spiffe format: %w", err)
		}
		b.spiffeBytes = spiffeBytes
		return spiffeBytes, nil
	default:
		return nil, fmt.Errorf("invalid format: %q", format)
	}
}

// toJWKS converts to JWKS the current bundle.
func (b *Formatter) toJWKS() ([]byte, error) {
	var jwks jose.JSONWebKeySet

	x509Authorities, jwtAuthorities, err := getAuthorities(b.bundle)
	if err != nil {
		return nil, err
	}

	for _, rootCA := range x509Authorities {
		jwks.Keys = append(jwks.Keys, jose.JSONWebKey{
			Key:          rootCA.PublicKey,
			Certificates: []*x509.Certificate{rootCA},
		})
	}

	for keyID, jwtSigningKey := range jwtAuthorities {
		jwks.Keys = append(jwks.Keys, jose.JSONWebKey{
			Key:   jwtSigningKey,
			KeyID: keyID,
		})
	}

	return json.Marshal(jwks)
}

// toOIDCJWKS converts the current bundle to a JWKS for OIDC relying parties.
// Only the JWT authorities are included, in the order of the bundle. Every key
// has its "kid" and "use" set to "sig", and "alg" when the JWS algorithm of
// the key is known; a key with no known algorithm is included without "alg".
func (b *Formatter) toOIDCJWKS() ([]byte, error) {
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{},
	}

	for i, jwtKey := range b.bundle.JwtAuthorities {
		publicKey, err := x509.ParsePKIXPublicKey(jwtKey.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("unable to parse JWT signing key %d: %w", i, err)
		}
		jwks.Keys = append(jwks.Keys, jose.JSONWebKey{
			Key:       publicKey,
			KeyID:     jwtKey.KeyId,
			Algorithm: string(joseAlgFromPublicKey(publicKey)),
			Use:       oidcKeyUse,
		})
	}

	return json.Marshal(jwks)
}

// joseAlgFromPublicKey returns the JWS algorithm used with the public key, as
// SPIRE does to sign JWT-SVIDs, or an empty string if it is not known: RSA
// keys of at least 2048 bits use RS256, and ECDSA keys on the P-256 and P-384
// curves use ES256 and ES384.
func joseAlgFromPublicKey(publicKey crypto.PublicKey) jose.SignatureAlgorithm {
	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		if publicKey.Size() < 256 {
			return ""
		}
		return jose.RS256
	case *ecdsa.PublicKey:
		switch publicKey.Params().BitSize {
		case 256:
			return jose.ES256
		case 384:
			return jose.ES384
		}
	}
	return ""
}

// toPEM converts to PEM the current bundle.
func (b *Formatter) toPEM() ([]byte, error) {
	bundleData := new(bytes.Buffer)
	for _, x509Authority := range b.bundle.X509Authorities {
		if err := pem.Encode(bundleData, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: x509Authority.Asn1,
		}); err != nil {
			return nil, fmt.Errorf("could not perform PEM encoding: %w", err)
		}
	}

	return bundleData.Bytes(), nil
}

// toSPIFFEBundle converts to a SPIFFE bundle the current bundle.
func (b *Formatter) toSPIFFEBundle() ([]byte, error) {
	sb, err := spiffeBundleFromPluginProto(b.bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bundle: %w", err)
	}
	docBytes, err := sb.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bundle: %w", err)
	}

	return docBytes, nil
}

// FormatBundle returns the bundle in the form of a slice of bytes in
// the chosen format.
func FormatBundle(bundle *types.Bundle, format Format) ([]byte, error) {
	return NewFormatter(bundle).Format(format)
}

// getAuthorities gets the X.509 authorities and JWT authorities from the
// provided *types.Bundle.
func getAuthorities(bundleProto *types.Bundle) ([]*x509.Certificate, map[string]crypto.PublicKey, error) {
	x509Authorities, err := x509CertificatesFromProto(bundleProto.X509Authorities)
	if err != nil {
		return nil, nil, err
	}
	jwtAuthorities, err := jwtKeysFromProto(bundleProto.JwtAuthorities)
	if err != nil {
		return nil, nil, err
	}

	return x509Authorities, jwtAuthorities, nil
}

// jwtKeysFromProto converts JWT keys from the given []*types.JWTKey to
// map[string]crypto.PublicKey.
// The key ID of the public key is used as the key in the returned map.
func jwtKeysFromProto(proto []*types.JWTKey) (map[string]crypto.PublicKey, error) {
	keys := make(map[string]crypto.PublicKey)
	for i, publicKey := range proto {
		jwtSigningKey, err := x509.ParsePKIXPublicKey(publicKey.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("unable to parse JWT signing key %d: %w", i, err)
		}
		keys[publicKey.KeyId] = jwtSigningKey
	}
	return keys, nil
}

// spiffeBundleFromPluginProto converts a bundle from the given *types.Bundle to
// *spiffebundle.Bundle.
func spiffeBundleFromPluginProto(bundleProto *types.Bundle) (*spiffebundle.Bundle, error) {
	td, err := spiffeid.TrustDomainFromString(bundleProto.TrustDomain)
	if err != nil {
		return nil, err
	}
	x509Authorities, jwtAuthorities, err := getAuthorities(bundleProto)
	if err != nil {
		return nil, err
	}

	bundle := spiffebundle.New(td)
	bundle.SetX509Authorities(x509Authorities)
	bundle.SetJWTAuthorities(jwtAuthorities)
	if bundleProto.RefreshHint > 0 {
		bundle.SetRefreshHint(time.Duration(bundleProto.RefreshHint) * time.Second)
	}
	if bundleProto.SequenceNumber > 0 {
		bundle.SetSequenceNumber(bundleProto.SequenceNumber)
	}
	return bundle, nil
}

// x509CertificatesFromProto converts X.509 certificates from the given
// []*types.X509Certificate to []*x509.Certificate.
func x509CertificatesFromProto(proto []*types.X509Certificate) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for i, auth := range proto {
		cert, err := x509.ParseCertificate(auth.Asn1)
		if err != nil {
			return nil, fmt.Errorf("unable to parse root CA %d: %w", i, err)
		}
		certs = append(certs, cert)
	}
	return certs, nil
}
