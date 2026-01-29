package handlers

import (
	"net/http"
	"real-time-forum/backend/domain"
	"real-time-forum/backend/middleware"
)

func Router(cs domain.ClientService, ps domain.PostsService, cts domain.CommentsService) http.Handler {
	clientService = cs
	postsService = ps
	commentsService = cts
	middleware.SetClientService(cs)

	mux := http.NewServeMux()

	// Routes publiques (pas d'authentification requise)
	mux.Handle("/api/login", middleware.CORS(http.HandlerFunc(LoginHandler)))
	mux.Handle("/api/register", middleware.CORS(http.HandlerFunc(RegisterHandler)))
	mux.Handle("/api/check-auth", middleware.CORS(http.HandlerFunc(CheckAuthHandler)))
	mux.Handle("/api/logout", middleware.CORS(http.HandlerFunc(LogoutHandler)))
	mux.Handle("/api/posts", middleware.CORS(http.HandlerFunc(GetAllPostsHandler)))

	// Routes protégées (authentification requise)
	mux.Handle("/api/posts/{postId}/comments", middleware.CORS(middleware.HandleAuth(http.HandlerFunc(NewCommentHandler))))
	mux.Handle("/api/posts/{postId}/reaction", middleware.CORS(middleware.HandleAuth(http.HandlerFunc(ReactionHandler))))

	// Route pour un post spécifique (doit être après les routes plus spécifiques)
	mux.Handle("/api/posts/", middleware.CORS(http.HandlerFunc(GetPostByIDHandler)))

	// Routes:
	// mux.Handle("/thread", middleware.Handle(http.HandlerFunc(ThreadHandler)))
	// mux.HandleFunc("/login", AuthenticateHandler)

	return middleware.CORS(mux)
}
