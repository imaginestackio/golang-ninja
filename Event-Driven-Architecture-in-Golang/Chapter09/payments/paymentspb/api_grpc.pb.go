// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: paymentspb/api.proto

package paymentspb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PaymentsServiceClient is the client API for PaymentsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentsServiceClient interface {
	AuthorizePayment(ctx context.Context, in *AuthorizePaymentRequest, opts ...grpc.CallOption) (*AuthorizePaymentResponse, error)
	ConfirmPayment(ctx context.Context, in *ConfirmPaymentRequest, opts ...grpc.CallOption) (*ConfirmPaymentResponse, error)
	CreateInvoice(ctx context.Context, in *CreateInvoiceRequest, opts ...grpc.CallOption) (*CreateInvoiceResponse, error)
	AdjustInvoice(ctx context.Context, in *AdjustInvoiceRequest, opts ...grpc.CallOption) (*AdjustInvoiceResponse, error)
	PayInvoice(ctx context.Context, in *PayInvoiceRequest, opts ...grpc.CallOption) (*PayInvoiceResponse, error)
	CancelInvoice(ctx context.Context, in *CancelInvoiceRequest, opts ...grpc.CallOption) (*CancelInvoiceResponse, error)
}

type paymentsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentsServiceClient(cc grpc.ClientConnInterface) PaymentsServiceClient {
	return &paymentsServiceClient{cc}
}

func (c *paymentsServiceClient) AuthorizePayment(ctx context.Context, in *AuthorizePaymentRequest, opts ...grpc.CallOption) (*AuthorizePaymentResponse, error) {
	out := new(AuthorizePaymentResponse)
	err := c.cc.Invoke(ctx, "/paymentspb.PaymentsService/AuthorizePayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsServiceClient) ConfirmPayment(ctx context.Context, in *ConfirmPaymentRequest, opts ...grpc.CallOption) (*ConfirmPaymentResponse, error) {
	out := new(ConfirmPaymentResponse)
	err := c.cc.Invoke(ctx, "/paymentspb.PaymentsService/ConfirmPayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsServiceClient) CreateInvoice(ctx context.Context, in *CreateInvoiceRequest, opts ...grpc.CallOption) (*CreateInvoiceResponse, error) {
	out := new(CreateInvoiceResponse)
	err := c.cc.Invoke(ctx, "/paymentspb.PaymentsService/CreateInvoice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsServiceClient) AdjustInvoice(ctx context.Context, in *AdjustInvoiceRequest, opts ...grpc.CallOption) (*AdjustInvoiceResponse, error) {
	out := new(AdjustInvoiceResponse)
	err := c.cc.Invoke(ctx, "/paymentspb.PaymentsService/AdjustInvoice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsServiceClient) PayInvoice(ctx context.Context, in *PayInvoiceRequest, opts ...grpc.CallOption) (*PayInvoiceResponse, error) {
	out := new(PayInvoiceResponse)
	err := c.cc.Invoke(ctx, "/paymentspb.PaymentsService/PayInvoice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsServiceClient) CancelInvoice(ctx context.Context, in *CancelInvoiceRequest, opts ...grpc.CallOption) (*CancelInvoiceResponse, error) {
	out := new(CancelInvoiceResponse)
	err := c.cc.Invoke(ctx, "/paymentspb.PaymentsService/CancelInvoice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentsServiceServer is the server API for PaymentsService service.
// All implementations must embed UnimplementedPaymentsServiceServer
// for forward compatibility
type PaymentsServiceServer interface {
	AuthorizePayment(context.Context, *AuthorizePaymentRequest) (*AuthorizePaymentResponse, error)
	ConfirmPayment(context.Context, *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error)
	CreateInvoice(context.Context, *CreateInvoiceRequest) (*CreateInvoiceResponse, error)
	AdjustInvoice(context.Context, *AdjustInvoiceRequest) (*AdjustInvoiceResponse, error)
	PayInvoice(context.Context, *PayInvoiceRequest) (*PayInvoiceResponse, error)
	CancelInvoice(context.Context, *CancelInvoiceRequest) (*CancelInvoiceResponse, error)
	mustEmbedUnimplementedPaymentsServiceServer()
}

// UnimplementedPaymentsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentsServiceServer struct {
}

func (UnimplementedPaymentsServiceServer) AuthorizePayment(context.Context, *AuthorizePaymentRequest) (*AuthorizePaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthorizePayment not implemented")
}
func (UnimplementedPaymentsServiceServer) ConfirmPayment(context.Context, *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmPayment not implemented")
}
func (UnimplementedPaymentsServiceServer) CreateInvoice(context.Context, *CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateInvoice not implemented")
}
func (UnimplementedPaymentsServiceServer) AdjustInvoice(context.Context, *AdjustInvoiceRequest) (*AdjustInvoiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AdjustInvoice not implemented")
}
func (UnimplementedPaymentsServiceServer) PayInvoice(context.Context, *PayInvoiceRequest) (*PayInvoiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PayInvoice not implemented")
}
func (UnimplementedPaymentsServiceServer) CancelInvoice(context.Context, *CancelInvoiceRequest) (*CancelInvoiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelInvoice not implemented")
}
func (UnimplementedPaymentsServiceServer) mustEmbedUnimplementedPaymentsServiceServer() {}

// UnsafePaymentsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentsServiceServer will
// result in compilation errors.
type UnsafePaymentsServiceServer interface {
	mustEmbedUnimplementedPaymentsServiceServer()
}

func RegisterPaymentsServiceServer(s grpc.ServiceRegistrar, srv PaymentsServiceServer) {
	s.RegisterService(&PaymentsService_ServiceDesc, srv)
}

func _PaymentsService_AuthorizePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizePaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServiceServer).AuthorizePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paymentspb.PaymentsService/AuthorizePayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServiceServer).AuthorizePayment(ctx, req.(*AuthorizePaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentsService_ConfirmPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServiceServer).ConfirmPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paymentspb.PaymentsService/ConfirmPayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServiceServer).ConfirmPayment(ctx, req.(*ConfirmPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentsService_CreateInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateInvoiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServiceServer).CreateInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paymentspb.PaymentsService/CreateInvoice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServiceServer).CreateInvoice(ctx, req.(*CreateInvoiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentsService_AdjustInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdjustInvoiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServiceServer).AdjustInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paymentspb.PaymentsService/AdjustInvoice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServiceServer).AdjustInvoice(ctx, req.(*AdjustInvoiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentsService_PayInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PayInvoiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServiceServer).PayInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paymentspb.PaymentsService/PayInvoice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServiceServer).PayInvoice(ctx, req.(*PayInvoiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentsService_CancelInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelInvoiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServiceServer).CancelInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paymentspb.PaymentsService/CancelInvoice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServiceServer).CancelInvoice(ctx, req.(*CancelInvoiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PaymentsService_ServiceDesc is the grpc.ServiceDesc for PaymentsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PaymentsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "paymentspb.PaymentsService",
	HandlerType: (*PaymentsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthorizePayment",
			Handler:    _PaymentsService_AuthorizePayment_Handler,
		},
		{
			MethodName: "ConfirmPayment",
			Handler:    _PaymentsService_ConfirmPayment_Handler,
		},
		{
			MethodName: "CreateInvoice",
			Handler:    _PaymentsService_CreateInvoice_Handler,
		},
		{
			MethodName: "AdjustInvoice",
			Handler:    _PaymentsService_AdjustInvoice_Handler,
		},
		{
			MethodName: "PayInvoice",
			Handler:    _PaymentsService_PayInvoice_Handler,
		},
		{
			MethodName: "CancelInvoice",
			Handler:    _PaymentsService_CancelInvoice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "paymentspb/api.proto",
}
