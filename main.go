package main 

import(
	"maqola-backent/routes" // router fayilini chaqirib oldik
	"github.com/gin-gonic/gin" // gin ramkasini chaqirib oldik
)

func main(){
     //1. yangi server divigatelini yaratamiz 
	 router := gin.Default()
	 //2. mashrutlarni serverga ulaymiz 
     routes.SetupRoutes(router)
	 //3. serverni 8080 portda ishga tushuramiz 
	 router.Run(":8080")
}
