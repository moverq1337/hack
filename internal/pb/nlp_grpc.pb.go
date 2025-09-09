package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	NLPService_ParseResume_FullMethodName        = "/pb.NLPService/ParseResume"
	NLPService_MatchResumeVacancy_FullMethodName = "/pb.NLPService/MatchResumeVacancy"
)

type NLPServiceClient interface {
	ParseResume(ctx context.Context, in *ParseRequest, opts ...grpc.CallOption) (*ParseResponse, error)
	MatchResumeVacancy(ctx context.Context, in *MatchRequest, opts ...grpc.CallOption) (*MatchResponse, error)
}

type nLPServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNLPServiceClient(cc grpc.ClientConnInterface) NLPServiceClient {
	return &nLPServiceClient{cc}
}

func (c *nLPServiceClient) ParseResume(ctx context.Context, in *ParseRequest, opts ...grpc.CallOption) (*ParseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ParseResponse)
	err := c.cc.Invoke(ctx, NLPService_ParseResume_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nLPServiceClient) MatchResumeVacancy(ctx context.Context, in *MatchRequest, opts ...grpc.CallOption) (*MatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MatchResponse)
	err := c.cc.Invoke(ctx, NLPService_MatchResumeVacancy_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type NLPServiceServer interface {
	ParseResume(context.Context, *ParseRequest) (*ParseResponse, error)
	MatchResumeVacancy(context.Context, *MatchRequest) (*MatchResponse, error)
	mustEmbedUnimplementedNLPServiceServer()
}

type UnimplementedNLPServiceServer struct{}

func (UnimplementedNLPServiceServer) ParseResume(context.Context, *ParseRequest) (*ParseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseResume not implemented")
}
func (UnimplementedNLPServiceServer) MatchResumeVacancy(context.Context, *MatchRequest) (*MatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MatchResumeVacancy not implemented")
}
func (UnimplementedNLPServiceServer) mustEmbedUnimplementedNLPServiceServer() {}
func (UnimplementedNLPServiceServer) testEmbeddedByValue()                    {}

type UnsafeNLPServiceServer interface {
	mustEmbedUnimplementedNLPServiceServer()
}

func RegisterNLPServiceServer(s grpc.ServiceRegistrar, srv NLPServiceServer) {
	
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NLPService_ServiceDesc, srv)
}

func _NLPService_ParseResume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NLPServiceServer).ParseResume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NLPService_ParseResume_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NLPServiceServer).ParseResume(ctx, req.(*ParseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NLPService_MatchResumeVacancy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NLPServiceServer).MatchResumeVacancy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NLPService_MatchResumeVacancy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NLPServiceServer).MatchResumeVacancy(ctx, req.(*MatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var NLPService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.NLPService",
	HandlerType: (*NLPServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ParseResume",
			Handler:    _NLPService_ParseResume_Handler,
		},
		{
			MethodName: "MatchResumeVacancy",
			Handler:    _NLPService_MatchResumeVacancy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/nlp.proto",
}
