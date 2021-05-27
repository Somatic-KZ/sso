package resources

import (
	"context"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/Somatic-KZ/sso-client/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Verify struct {
	db drivers.DataStore
}

func NewVerify(db drivers.DataStore) *Verify {
	return &Verify{db: db}
}

func (v *Verify) UserToken(ctx context.Context, req *protobuf.UserTokenRequest) (*emptypb.Empty, error) {
	err := v.db.VerifyToken(ctx, req.Id)
	switch err {
	case drivers.ErrTokenNotFound:
		return nil, status.Error(codes.NotFound, "token expired or not exists")
	case nil:
		return &empty.Empty{}, nil
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}
