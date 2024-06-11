package server

import (
	"context"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"
)

type JobetServer struct {
	proto.UnimplementedJobetServer

	q *db.Queries
}

func NewJobetServer(q *db.Queries) *JobetServer {
	return &JobetServer{q: q}
}

func (s *JobetServer) Probe(ctx context.Context, req *proto.ProbeRequest) (*proto.ProbeReply, error) {
	return &proto.ProbeReply{}, nil
}
