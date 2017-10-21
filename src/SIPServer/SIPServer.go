package SIPServer

import (
	"net/http"

	"SIPCore"
)

type APIServer struct {
	Context *SIPCore.Context
}

type GetGroupIntersectionRequest struct {
	Other string `json:"other"`
}

func NewAPIServer(ctx *SIPCore.Context) *APIServer {
	return &APIServer {
		Context: ctx,
	}
}

func (s *APIServer) Run() {
	server := http.NewServeMux()
	server.HandleFunc("/group/intersection", s.wrapWithIdentityLoader(s.onGetGroupIntersection))

	hs := &http.Server {
		Addr: s.Context.Config.ListenAddr,
	}
	hs.Handler = server
	hs.ListenAndServe()
}

func (s *APIServer) onGetGroupIntersection(w http.ResponseWriter, r *http.Request, id *SIPCore.Identity) {
	req := &GetGroupIntersectionRequest{}
	err := ReadAPIRequest(r, &req)
	if err != nil {
		BuildAPIResult(false, "Invalid request").Write(w)
		return
	}

	other, err := s.Context.GetIdentityByName(req.Other)
	if err != nil {
		BuildAPIResult(false, "Requested target not found").Write(w)
		return
	}

	inter := id.GroupIntersection(other)
	BuildAPIResult(true, inter).Write(w)
}

func (s *APIServer) wrapWithIdentityLoader(
	target func (w http.ResponseWriter, r *http.Request, id *SIPCore.Identity),
) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Service-Key")
		id, err := s.Context.GetIdentityByKey(key)

		if err != nil {
			BuildAPIResult(false, "Invalid service key provided in X-Service-Key header").Write(w)
			return
		}

		target(w, r, id)
	}
}
