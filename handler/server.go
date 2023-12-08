package handler

import "github.com/SawitProRecruitment/UserService/repository"

type Server struct {
	Repository repository.RepositoryInterface
	SecretKey  string
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	SecretKey  string
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		SecretKey:  opts.SecretKey,
	}
}
