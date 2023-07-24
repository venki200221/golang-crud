package controllers

import (
"errors"
"net/http"
"go_basics/configs"
"go_basics/configs/postgres"
"github.com/gin-gonic/gin"
)


func CreateBook(c *gin.Context){
	var book *postgres.Book
	err:=c.ShouldBind(&book)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	res:=configs.DB1.Create(&book)
	if res.Error!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":res.Error})
		return
	}
	c.JSON(http.StatusOK,gin.H{"book":book})
	return 


}


func ReadBooks(c *gin.Context){
	var books []postgres.Book
	res :=configs.DB1.Find(&books)
	if res.Error!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":res.Error})
		return
	}
	c.JSON(http.StatusOK,gin.H{"books":books})
	return
}

func UpdateBook(c *gin.Context){
	var book *postgres.Book
	id :=c.Param("id")
	err :=c.ShouldBind(&book)

	if err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return 
	}
	var updateBook postgres.Book
	res:=postgres.DB1.Where("id=?",id).updates(book)
	if res.RowAffected==0{
        c.JSON(http.StatusNotFound,gin.H{"error":errors.New("book not found")})
		return
	}
	c.JSON(http.StatusOK,gin.H{"book":updateBook})
	return

}

func ReadBook(c *gin.Context){
	var book postgres.Book
	id :=c.Param("id")
	res :=postgres.DB1.Where("id=?",id).Find(&book)
	if res.RowAffected==0{
		c.JSON(http.StatusNotFound,gin.H{"error":errors.New("book not found")})
		return
	}
	c.JSON(http.StatusOK,gin.H{"book":book})
	return
}

func DeleteBook(c *gin.Context){
	var book postgres.Book
	id :=c.Param("id")
	res :=postgres.DB1.Find(&book,id)
	if res.RowAffected==0{
		c.JSON(http.StatusNotFound,gin.H{"error":errors.New("book not found")})
		return
	}
	postgres.DB1.Delete(&book,id)
	c.JSON(http.StatusOK,gin.H{"message":"book deleted suicessfully"})
	return
}