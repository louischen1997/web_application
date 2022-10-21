package call

import (
	"Golangapi/pojo"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var hasAuth bool = false

func GetAllUsers(c *gin.Context) {
	users := pojo.GetAllUsers_db()
	c.JSON(200, users)
}

func GetAllUsersToCheck(email string) bool {
	users := pojo.GetAllUsers_db()
	for index, element := range users {
		//fmt.Println(element.Username)
		if element.Username == email {
			return false
		}
		index += 1
	}
	return true
}

func GetUsers(c *gin.Context) {
	if !hasAuth {
		return
	}
	user := pojo.GetUsers_db(c.Param("id"))
	c.JSON(http.StatusOK, gin.H{
		"first_name":      user.First_name,
		"last_name":       user.Last_name,
		"username":        user.Username,
		"id":              user.ID,
		"account_created": user.Account_created,
		"account_updated": user.Account_updated,
	})
	hasAuth = false
}

func PostUsers(c *gin.Context) {
	user := pojo.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.Status(400)
		return
	}

	if !GetAllUsersToCheck(user.Username) {
		c.Status(400)
		return
	}
	if user.ID != "" || user.Account_created != "" || user.Account_updated != "" {
		c.Status(400)
		return
	}

	hashPass := []byte(user.Password)
	user.Password = HashAndSalt(hashPass)

	id := uuid.New().String()
	now := time.Now().String()

	user.ID = id
	user.Account_created = now
	user.Account_updated = now

	pojo.PostUsers_db(user)
	c.JSON(201, gin.H{
		"first_name":      user.First_name,
		"last_name":       user.Last_name,
		"username":        user.Username,
		"id":              user.ID,
		"account_created": user.Account_created,
		"account_updated": user.Account_updated,
	})

}

func DeleteUser(c *gin.Context) {
	pojo.DeleteUser(c.Param("id"))
	c.Status(http.StatusOK)
}

func DeleteAllUsers(c *gin.Context) {
	users := pojo.GetAllUsers_db()
	for index, element := range users {
		pojo.DeleteUser(element.ID)
		index += 1
	}
	c.JSON(200, "DELETE ALL")
}

func PutUser(c *gin.Context) {
	if !hasAuth {
		return
	}
	user := pojo.User{}
	now := time.Now().String()

	err := c.BindJSON(&user)
	if err != nil {
		c.Status(400)
		return
	}
	tmp_user := pojo.GetUsers_db(c.Param("id"))
	if user.ID != "" || user.Account_created != "" || user.Account_updated != "" || user.Username != tmp_user.Username {
		c.Status(400)
		return
	}
	hashPass := []byte(user.Password)
	user.Password = HashAndSalt(hashPass)
	user.Account_updated = now
	pojo.UpdateUser(c.Param("id"), user)
	hasAuth = false
	c.Status(204)
}

func HandleBA(c *gin.Context) {
	var username_ba, password_ba string
	users := pojo.GetAllUsers_db()
	username_ba, password_ba, hasAuth = c.Request.BasicAuth()
	if len(users) == 0 {
		c.Status(401)
		hasAuth = false
		return
	}
	for index, element := range users {
		if index == len(users)-1 && element.ID != c.Param("id") {
			c.Status(401)
			hasAuth = false
			return
		}
		if element.ID == c.Param("id") {
			break
		}
	}
	tmpHash := pojo.GetUsers_db_Pass(c.Param("id"))
	if username_ba == pojo.GetUsers_db_Username(c.Param("id")) && CheckPasswordHash(password_ba, tmpHash) {
		return
	}
	hasAuth = false
	c.Status(403)
}

func HashAndSalt(pass []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
