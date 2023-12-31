package controllers

import (
	"context"
	"go_basics/configs"
	"go_basics/responses"
	"go_basics/models"
	"net/http"
	"time"
    "go.mongodb.org/mongo-driver/bson"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection=configs.GetCollection(configs.DB,"users")
var validate=validator.New()

func CreateUser(c *fiber.Ctx)error{
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
    var user  models.User
	defer cancel()

	if err:=c.BodyParser(&user);err!=nil{
	return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status:http.StatusBadRequest,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})
	}

	if validationerr := validate.Struct(&user); validationerr!=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status:http.StatusBadRequest,Message:validationerr.Error(),Data:&fiber.Map{"data":validationerr.Error()}})
	}
    
	newUser:=models.User{
		Id :primitive.NewObjectID(),
		Name: user.Name,
		Location: user.Location,
		Title: user.Title,
	}
	result,err :=userCollection.InsertOne(ctx,newUser)
	if err!=nil{
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status:http.StatusInternalServerError,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})
	}
	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status:http.StatusCreated,Message:"success",Data:&fiber.Map{"data":result}})
	

}

func GetAUser(c *fiber.Ctx)error{
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	userId := c.Params("userid")
	var user models.User
	defer cancel()
	objId,_:=primitive.ObjectIDFromHex(userId)

	err:=userCollection.FindOne(ctx,bson.M{"_id":objId}).Decode(&user)
	if err!=nil{
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status:http.StatusInternalServerError,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status:http.StatusOK,Message:"success",Data:&fiber.Map{"data":user}})


}


func EditUser(c *fiber.Ctx)error{
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	userId := c.Params("userid")
	var user models.User
	defer cancel()

	objId,_:=primitive.ObjectIDFromHex(userId)
	if err:=c.BodyParser(&user);err!=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status:http.StatusBadRequest,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})
	}

	if validationerr := validate.Struct(&user); validationerr!=nil{
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status:http.StatusBadRequest,Message:validationerr.Error(),Data:&fiber.Map{"data":validationerr.Error()}})
	}

	update :=bson.M{"name":user.Name,"location":user.Location,"title":user.Title}
	result,err :=userCollection.UpdateOne(ctx,bson.M{"_id":objId},bson.M{"$set":update})
	if err !=nil{
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status:http.StatusInternalServerError,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})
	}
    var updateduser models.User
	if result.MatchedCount==1{
		err:=userCollection.FindOne(ctx,bson.M{"_id":objId}).Decode(&updateduser)
		if err!=nil{return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status:http.StatusInternalServerError,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})}
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status:http.StatusOK,Message:"success",Data:&fiber.Map{"data":updateduser}})
}

func DeleteUser(c *fiber.Ctx)error{
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	var userId= c.Params("userid")
	defer cancel()

	objId,_:=primitive.ObjectIDFromHex(userId)
	result,err :=userCollection.DeleteOne(ctx,bson.M{"_id":objId})
	if err!=nil{
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status:http.StatusInternalServerError,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})
	}
	if result.DeletedCount<1{
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{Status:http.StatusNotFound,Message:"user not found",Data:&fiber.Map{"data":"user not found"}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status:http.StatusOK,Message:"success",Data:&fiber.Map{"data":"user deleted successfully"}})
}

func GetAllUsers(c *fiber.Ctx)error{
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	var users []models.User
	defer cancel()

	results,err :=userCollection.Find(ctx,bson.M{})
	if err!=nil{
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status:http.StatusInternalServerError,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})

	}
	defer results.Close(ctx)

	for results.Next(ctx){
		var singleUser models.User
		if err=results.Decode(&singleUser);err!=nil{
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status:http.StatusInternalServerError,Message:err.Error(),Data:&fiber.Map{"data":err.Error()}})
		}
		users = append(users, singleUser)
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status:http.StatusOK,Message:"success",Data:&fiber.Map{"data":users}})
}

