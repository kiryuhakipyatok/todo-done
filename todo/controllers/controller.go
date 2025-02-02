package controllers

import (
	"todo/db"
	"todo/models"
	"time"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secret string = "secret"

func Index(c *fiber.Ctx) error{
	return c.Render("index",fiber.Map{})
}

func RegisterTMPL(c *fiber.Ctx) error{
	return c.Render("register",fiber.Map{})
}

func LoginTMPL(c *fiber.Ctx) error{
	return c.Render("login",fiber.Map{})
}

func Home(c *fiber.Ctx) error {
    cookie := c.Cookies("jwt")
    if cookie == "" {
        return c.Redirect("/log")
    }
    return c.Render("home",fiber.Map{})
}

func Register(c *fiber.Ctx) error{
	
	data:=map[string]string{}
	
	if err:=c.BodyParser(&data); err!=nil{
		return err
	}
	
	password,err:=bcrypt.GenerateFromPassword([]byte(data["password"]),14)
	if err!=nil{
		return err
	}
	user:=models.User{
		Id: uuid.New(),
		Nickname:data["nickname"],
		Email:data["email"],
		Login:data["login"],
		Password: password,
	}

	if err:=db.DB.Create(&user).Error;err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating user: " + err.Error(),
		})
	}

	return c.JSON(user)
}	



func AddTask(c *fiber.Ctx) error{
	taskdata:=map[string]string{}
	
	if err:=c.BodyParser(&taskdata); err!=nil{
		return err
	}

	userID, err := uuid.Parse(taskdata["user_id"])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Неверный формат user_id",
		})
	}

	task:=models.Tasks{
		Task_id: uuid.New(),
		User_id: userID,
		Time: taskdata["time"],
		Task: taskdata["task"],
	}

	if err:=db.DB.Create(&task).Error;err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating task: " + err.Error(),
		})
	}
	return c.JSON(task)
}

func GetTasks(c *fiber.Ctx) error {
    userId := c.Params("id") 

    var tasks []models.Tasks

   
    if err := db.DB.Where("user_id = ?", userId).Find(&tasks).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error fetching tasks: " + err.Error(),
        })
    }

    return c.JSON(tasks) 
}

func DeleteTask(c *fiber.Ctx) error{
	taskId:=c.Params("taskid")
	uuidTaskID, err := uuid.Parse(taskId)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Неверный формат task_id",
        })
    }
	if err:=db.DB.Where("task_id = ?", uuidTaskID).Delete(&models.Tasks{}).Error;err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Ошибка при удалении задачи",
        })
	}
	return c.JSON(fiber.Map{
        "message": "Задача успешно удалена",
    })
}	

func Login(c *fiber.Ctx) error{
	data:=map[string]string{}
	if err:=c.BodyParser(&data); err!=nil{
		return err
	}
	user :=models.User{}
	
	if errors.Is(db.DB.Where("login=?", data["login"]).First(&user).Error,gorm.ErrRecordNotFound){
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message":"такого пользователя нет",
		})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password,[]byte(data["password"]));err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message":"неправильный пароль",
		})
		
	}

	claims:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
		Issuer: user.Id.String(),
		ExpiresAt: time.Now().Add(time.Hour*24).Unix(),
	})

	token,err:=claims.SignedString([]byte(secret))
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message":"could not login",
		})
		
	}
	cookie:=fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}	



func User(c *fiber.Ctx) error{
	cookie :=c.Cookies("jwt")
	token,err:=jwt.ParseWithClaims(cookie,&jwt.StandardClaims{},func(t *jwt.Token) (interface{}, error) {
		return []byte(secret),nil
	})
	if err!=nil{
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message":"unauthenticated",
		})
		
	}
	claims:= token.Claims.(*jwt.StandardClaims)
	user :=models.User{}
	db.DB.Where("id=?",claims.Issuer).First(&user)
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error{
	cookie:=fiber.Cookie{
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"succes",
	})
}