package handler

import (
	"context"
	"time"

	"github.com/lecex/socialite-api/config"
	pb "github.com/lecex/socialite-api/proto/health"
)

// Health 用户结构
type Health struct {
}

// Health 用户是否存在
func (srv *Health) Health(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	res.Valid = true
	res.Time = time.Now().Format("2006-01-02 15:04:05")
	res.Service = config.Conf.Name
	return nil
}
