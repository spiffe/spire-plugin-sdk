package plugintest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/plugintest"
	"github.com/spiffe/spire-plugin-sdk/private/proto/test"
	nodeattestorv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/agent/nodeattestor/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const panicDetails = "sensitive panic details"

func TestServeRecoversUnaryPanics(t *testing.T) {
	for _, tt := range []struct {
		name    string
		service bool
		value   any
		message string
	}{
		{name: "plugin", value: panicDetails, message: panicDetails},
		{name: "service", service: true, value: panicDetails, message: panicDetails},
		{name: "typed nil error", value: (*os.PathError)(nil), message: "<nil>"},
		{
			name:    "panicking error method",
			value:   panickingError{},
			message: "%!v(PANIC=Error method: " + panicDetails + ")",
		},
		{name: "panicking JSON marshaler", value: panickingJSON{}, message: "panicking JSON marshaler"},
		{name: "complex", value: complex(1, 2), message: "(1+2i)"},
		{name: "NaN", value: math.NaN(), message: "NaN"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := &panickingPlugin{panicValue: tt.value}
			pluginClient := new(test.SomePluginPluginClient)
			serviceClient := new(test.SomeServiceServiceClient)
			method := test.SomePlugin_PluginEcho_FullMethodName
			if tt.service {
				method = test.SomeService_ServiceEcho_FullMethodName
			}
			plugintest.ServeInBackground(t, plugintest.Config{
				PluginServer:   test.SomePluginPluginServer(server),
				PluginClient:   pluginClient,
				ServiceServers: []pluginsdk.ServiceServer{test.SomeServiceServiceServer(server)},
				ServiceClients: []pluginsdk.ServiceClient{serviceClient},
				Logger:         newRecoveryLogger(t, method, tt.message),
			})

			call := pluginClient.PluginEcho
			if tt.service {
				call = serviceClient.ServiceEcho
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := call(ctx, &test.EchoRequest{In: "panic"})
			if resp != nil {
				t.Fatalf("expected no response after a panic; got %v", resp)
			}
			if status.Code(err) != codes.Internal {
				t.Fatalf("expected Internal error; got %v", err)
			}
			assertStringEqual(t, "internal error", status.Convert(err).Message())

			// The same plugin must still serve successful calls and ordinary errors.
			resp, err = call(ctx, &test.EchoRequest{In: "hello"})
			if err != nil {
				t.Fatalf("echo after a panic failed unexpectedly with %v", err)
			}
			assertStringEqual(t, "hello", resp.Out)

			_, err = call(ctx, &test.EchoRequest{In: "error"})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("expected InvalidArgument error; got %v", err)
			}
			assertStringEqual(t, "invalid input", status.Convert(err).Message())
		})
	}
}

func TestServeRecoversStreamPanics(t *testing.T) {
	for _, tt := range []struct {
		name        string
		challenge   string
		wantPayload bool
	}{
		{name: "before response", challenge: "panic"},
		{name: "after response", challenge: "panic after send", wantPayload: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client := new(nodeattestorv1.NodeAttestorPluginClient)
			plugintest.ServeInBackground(t, plugintest.Config{
				PluginServer: nodeattestorv1.NodeAttestorPluginServer(new(panickingNodeAttestor)),
				PluginClient: client,
				Logger:       newRecoveryLogger(t, nodeattestorv1.NodeAttestor_AidAttestation_FullMethodName, panicDetails),
			})

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			stream, err := client.AidAttestation(ctx)
			if err != nil {
				t.Fatalf("AidAttestation failed unexpectedly with %v", err)
			}
			if err := stream.Send(&nodeattestorv1.Challenge{Challenge: []byte(tt.challenge)}); err != nil {
				t.Fatalf("sending challenge failed unexpectedly with %v", err)
			}
			if tt.wantPayload {
				resp, err := stream.Recv()
				if err != nil {
					t.Fatalf("receiving payload failed unexpectedly with %v", err)
				}
				assertStringEqual(t, "payload", string(resp.GetPayload()))
			}
			_, err = stream.Recv()
			if status.Code(err) != codes.Internal {
				t.Fatalf("expected Internal error; got %v", err)
			}
			assertStringEqual(t, "internal error", status.Convert(err).Message())

			// Recovery must leave the server able to handle another stream.
			stream, err = client.AidAttestation(ctx)
			if err != nil {
				t.Fatalf("AidAttestation after a panic failed unexpectedly with %v", err)
			}
			if err := stream.Send(&nodeattestorv1.Challenge{Challenge: []byte("hello")}); err != nil {
				t.Fatalf("sending challenge failed unexpectedly with %v", err)
			}
			resp, err := stream.Recv()
			if err != nil {
				t.Fatalf("receiving payload failed unexpectedly with %v", err)
			}
			assertStringEqual(t, "payload", string(resp.GetPayload()))
			_, err = stream.Recv()
			if !errors.Is(err, io.EOF) {
				t.Fatalf("expected end of stream; got %v", err)
			}

			stream, err = client.AidAttestation(ctx)
			if err != nil {
				t.Fatalf("AidAttestation failed unexpectedly with %v", err)
			}
			if err := stream.Send(&nodeattestorv1.Challenge{Challenge: []byte("error")}); err != nil {
				t.Fatalf("sending challenge failed unexpectedly with %v", err)
			}
			_, err = stream.Recv()
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("expected InvalidArgument error; got %v", err)
			}
			assertStringEqual(t, "invalid input", status.Convert(err).Message())
		})
	}
}

func newRecoveryLogger(t *testing.T, method, panicMessage string) hclog.Logger {
	t.Helper()
	var log bytes.Buffer
	// Registered before ServeInBackground's cleanups so the server is stopped
	// before its logs are read.
	t.Cleanup(func() {
		decoder := json.NewDecoder(&log)
		var entry struct {
			Level   string `json:"@level"`
			Message string `json:"@message"`
			Method  string `json:"method"`
			Panic   string `json:"panic"`
			Stack   string `json:"stack"`
		}
		if err := decoder.Decode(&entry); err != nil {
			t.Fatalf("decoding recovery log failed unexpectedly with %v", err)
		}
		assertStringEqual(t, "error", entry.Level)
		assertStringEqual(t, "Recovered a panic in a plugin call", entry.Message)
		assertStringEqual(t, method, entry.Method)
		assertStringEqual(t, panicMessage, entry.Panic)
		if !strings.Contains(entry.Stack, "goroutine ") {
			t.Fatalf("expected a goroutine stack trace; got %q", entry.Stack)
		}
		if !strings.Contains(entry.Stack, "plugintest/recover_test.go:") {
			t.Fatalf("expected stack trace to include the panicking handler; got %q", entry.Stack)
		}
		var extra any
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			t.Fatalf("expected no more log entries; got entry %v and error %v", extra, err)
		}
	})
	return hclog.New(&hclog.LoggerOptions{
		JSONFormat:  true,
		DisableTime: true,
		Output:      &log,
	})
}

type panickingPlugin struct {
	test.UnimplementedSomePluginServer
	test.UnimplementedSomeServiceServer

	panicValue any
}

func (p *panickingPlugin) PluginEcho(_ context.Context, req *test.EchoRequest) (*test.EchoResponse, error) {
	return echoOrPanic(req, p.panicValue)
}

func (p *panickingPlugin) ServiceEcho(_ context.Context, req *test.EchoRequest) (*test.EchoResponse, error) {
	return echoOrPanic(req, p.panicValue)
}

func echoOrPanic(req *test.EchoRequest, panicValue any) (*test.EchoResponse, error) {
	switch req.In {
	case "panic":
		panic(panicValue)
	case "error":
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	default:
		return &test.EchoResponse{Out: req.In}, nil
	}
}

type panickingNodeAttestor struct {
	nodeattestorv1.UnimplementedNodeAttestorServer
}

func (*panickingNodeAttestor) AidAttestation(stream nodeattestorv1.NodeAttestor_AidAttestationServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	switch string(req.Challenge) {
	case "panic":
		panic(panicDetails)
	case "error":
		return status.Error(codes.InvalidArgument, "invalid input")
	}
	if err := stream.Send(&nodeattestorv1.PayloadOrChallengeResponse{
		Data: &nodeattestorv1.PayloadOrChallengeResponse_Payload{Payload: []byte("payload")},
	}); err != nil {
		return err
	}
	if string(req.Challenge) == "panic after send" {
		panic(panicDetails)
	}
	return nil
}

type panickingError struct{}

func (panickingError) Error() string {
	panic(panicDetails)
}

type panickingJSON struct{}

func (panickingJSON) MarshalJSON() ([]byte, error) {
	panic(panicDetails)
}

func (panickingJSON) String() string {
	return "panicking JSON marshaler"
}
