package main

import (
	"os"
	"real-time-forum/backend/config"
	"real-time-forum/backend/handlers"
	"real-time-forum/backend/repositories"
	"real-time-forum/backend/services"

	// "real-time-forum/backend/repository"
	// "real-time-forum/backend/service"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("❌ error loading .env file")
	}
	db := config.InitDB()
	defer db.Close()

	clientRepository := repositories.NewLogRepository(db)
	clientService := services.NewAuthService(clientRepository)

	commentsRepository := repositories.NewCommentsRepository(db)
	commentsService := services.NewCommentsService(commentsRepository)

	postsRepository := repositories.NewPostsRepository(db)
	postsService := services.NewPostsService(postsRepository)

	router := handlers.Router(clientService, postsService, commentsService)
	addr := os.Getenv("SERVER_PORT")
	//router := SetupRouter(db)

	////////////////////////////////////////////////////////////////////
	////////////////// lancement du serveur ////////////////////////////
	////////////////////////////////////////////////////////////////////

	log.Printf("Server start → http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("❌ error trying to run the server: ", err)
	}
}
