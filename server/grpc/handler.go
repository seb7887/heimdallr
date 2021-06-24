package grpc

import (
	"context"

	hrpc "github.com/seb7887/heimdallr/proto/heimdallr"
	"github.com/seb7887/heimdallr/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type heimdallrGRPCHandler struct {
	grpcService service.Service
}

func NewHeimdallrGRPCServer(grpcService service.Service) hrpc.HeimdallrServiceServer {
	return &heimdallrGRPCHandler{
		grpcService: grpcService,
	}
}

func (h heimdallrGRPCHandler) CreateClient(ctx context.Context, req *hrpc.ClientIdRequest) (*hrpc.KeyPairResponse, error) {
	privateKey, err := h.grpcService.Create(ctx, req.ClientId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &hrpc.KeyPairResponse{
		PrivateKey: *privateKey,
	}, nil
}

func (h heimdallrGRPCHandler) Authenticate(ctx context.Context, req *hrpc.AuthRequest) (*hrpc.ResultResponse, error) {
	authenticated := h.grpcService.Authenticate(ctx, req.ClientId, req.JwtToken)
	return &hrpc.ResultResponse{Success: authenticated}, nil
}

func (h heimdallrGRPCHandler) RegenerateKeys(ctx context.Context, req *hrpc.ClientIdRequest) (*hrpc.KeyPairResponse, error) {
	privateKey, err := h.grpcService.Create(ctx, req.ClientId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &hrpc.KeyPairResponse{
		PrivateKey: *privateKey,
	}, nil
}

func (h heimdallrGRPCHandler) AddToBlacklist(ctx context.Context, req *hrpc.ClientIdRequest) (*hrpc.ResultResponse, error) {
	if err := h.grpcService.UpdateBlacklist(ctx, req.ClientId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &hrpc.ResultResponse{Success: true}, nil
}

func (h heimdallrGRPCHandler) GetBlacklist(ctx context.Context, req *hrpc.EmptyReq) (*hrpc.ClientIdsResponse, error) {
	blacklist, err := h.grpcService.ReadBlacklist(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &hrpc.ClientIdsResponse{Ids: blacklist}, nil
}

func (h heimdallrGRPCHandler) DeleteClient(ctx context.Context, req *hrpc.ClientIdRequest) (*hrpc.ResultResponse, error) {
	if err := h.grpcService.Delete(ctx, req.ClientId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &hrpc.ResultResponse{Success: true}, nil
}
