// Code generated by protoc-gen-go-spire. DO NOT EDIT.

package v1

import (
	pluginsdk "github.com/spiffe/spire-plugin-sdk/pluginsdk"
	grpc "google.golang.org/grpc"
)

func WorkloadAttestorPluginServer(server WorkloadAttestorServer) pluginsdk.PluginServer {
	return workloadAttestorPluginServer{WorkloadAttestorServer: server}
}

type workloadAttestorPluginServer struct {
	WorkloadAttestorServer
}

func (s workloadAttestorPluginServer) Type() string {
	return "WorkloadAttestor"
}

func (s workloadAttestorPluginServer) GRPCServiceName() string {
	return "spire.plugin.agent.workloadattestor.v1.WorkloadAttestor"
}

func (s workloadAttestorPluginServer) RegisterServer(server *grpc.Server) interface{} {
	RegisterWorkloadAttestorServer(server, s.WorkloadAttestorServer)
	return s.WorkloadAttestorServer
}

type WorkloadAttestorPluginClient struct {
	WorkloadAttestorClient
}

func (s WorkloadAttestorPluginClient) Type() string {
	return "WorkloadAttestor"
}

func (c *WorkloadAttestorPluginClient) IsInitialized() bool {
	return c.WorkloadAttestorClient != nil
}

func (c *WorkloadAttestorPluginClient) GRPCServiceName() string {
	return "spire.plugin.agent.workloadattestor.v1.WorkloadAttestor"
}

func (c *WorkloadAttestorPluginClient) InitClient(conn grpc.ClientConnInterface) interface{} {
	c.WorkloadAttestorClient = NewWorkloadAttestorClient(conn)
	return c.WorkloadAttestorClient
}
