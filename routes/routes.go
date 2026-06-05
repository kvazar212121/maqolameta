package routes 

import (
	"maqola-backent/controllers" // controlersni chaqirib olyabmiz 
	"github.com/gin-gonic/gin"
)
// SetupRouter -barcha manziullarni (api larni ) yolga qoyuvchi funksiya 
func SetupRoutes(router *gin.Engine) {
	//apilar chalkashmasligi uchun api guruhini yaratamiz
	api := router.Group("/api")
    {
		// 1 . agar browserdan "api ping deb " kelishsa , ping funksiyasini ishlat 
		api.GET("/ping",controllers.Ping)
		// 2 . agar  "/api/articles " dep kelishsa GetArticles funksiyasini ishlatib ber 
		api.GET("/articles",controllers.GetArticles)
	}

}
         
