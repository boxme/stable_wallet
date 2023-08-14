package main

import (
	"stable_wallet/main/login"
)

func (s *Server) startRouting() {
	// s.mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	loginHandler := login.CreateLoginHandler(s.App)
	s.Mux.HandleFunc("/user/login", loginHandler.HandleLogin())
}
