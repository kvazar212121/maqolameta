package controllers

import (
	
	"net/http"
	
	"github.com/gin-gonic/gin"
)


//ping funcksiyasini  -server ishlayotganini va  mashrut tog`rishidagi tekshiruvlar uchun 
func Ping(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message":"backnet tog`ri ishlayabdi akam",
	})

}

//GetArticles funksiyasi - bazadan maqolalarni olish uchun 

func GetArticles(c *gin.Context){
//hozricha maqolalar yoq shuning uchun  bosh royhat qaytaramiz 
       c.JSON(http.StatusOK, gin.H{
		"data": []string{}, //KEYINCHALIK BUYERGA MAQOLALAR YOZILADI
		"total": 0, //jami maqollar soni
	   })	
}

