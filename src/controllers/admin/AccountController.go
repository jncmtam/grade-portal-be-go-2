package controller_admin

import (
	"Go2/models"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// HandleCreateAccount xử lý việc tạo tài khoản mới.
func HandleCreateAccount(c *gin.Context) {
	var newAccounts []InterfaceAccount
	if err := c.ShouldBindJSON(&newAccounts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu yêu cầu không hợp lệ"})
		return
	}
	accountCollection := models.AccountModel()

	var existingAccounts []models.InterfaceAccount
	cursor, err := accountCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản từ cơ sở dữ liệu"})
		return
	}
	if err := cursor.All(context.TODO(), &existingAccounts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi giải mã tài khoản"})
		return
	}
	createdBy, _ := c.Get("ID")

	existingMap := make(map[string]bool)
	var validAccounts []InterfaceAccount
	var invalidAccounts []InterfaceAccount
	for _, account := range existingAccounts {
		existingMap[account.Email] = true
		existingMap[account.Ms] = true
	}
	for _, newAccount := range newAccounts {
		if !existingMap[newAccount.Email] && !existingMap[newAccount.Ms] && strings.HasSuffix(newAccount.Email, "@hcmut.edu.vn") && (newAccount.Role == "student" || newAccount.Role == "teacher") {
			newAccount.CreatedBy = createdBy
			newAccount.ExpiredAt = time.Now().AddDate(5, 0, 0)
			validAccounts = append(validAccounts, newAccount)
		} else {
			invalidAccounts = append(invalidAccounts, newAccount)
		}
	}
	if len(validAccounts) > 0 {
		if _, err := accountCollection.InsertMany(context.TODO(), validAccounts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo tài khoản hợp lệ"})
			return
		}
	}
	c.JSON(200, gin.H{
		"code":            "success",
		"invalidAccounts": invalidAccounts,
		"validAccounts":   validAccounts,
	})
}

// HandleGetAccountByID xử lý việc lấy thông tin tài khoản theo ID.
func HandleGetAccountByID(c *gin.Context) {
	idParam := c.Param("id")
	accountID, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "ID không hợp lệ",
		})
		return
	}

	accountCollection := models.AccountModel()
	var account models.InterfaceAccount
	err = accountCollection.FindOne(context.TODO(), bson.M{"_id": accountID}).Decode(&account)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tài khoản"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Tìm thấy tài khoản thành công",
		"account": account,
	})
}

// HandleGetTeacherAccounts xử lý việc lấy thông tin tài khoản giáo viên.
func HandleGetTeacherAccounts(c *gin.Context) {
	accountCollection := models.AccountModel()
	query := c.Query("ms")

	if query == "" {
		var teachers []models.InterfaceAccount
		cursor, err := accountCollection.Find(context.TODO(), bson.M{"role": "teacher"})
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tài khoản"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản"})
			return
		}
		if err := cursor.All(context.TODO(), &teachers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":      "Tìm thấy tài khoản thành công",
			"foundedUser": teachers,
		})
	} else {
		var teacher models.InterfaceAccount
		err := accountCollection.FindOne(context.TODO(), bson.M{"ms": query}).Decode(&teacher)
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tài khoản"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":      "Tìm thấy tài khoản thành công",
			"foundedUser": teacher,
		})
	}
}

// HandleGetStudentAccounts xử lý việc lấy thông tin tài khoản sinh viên.
func HandleGetStudentAccounts(c *gin.Context) {
	accountCollection := models.AccountModel()
	query := c.Query("ms")

	if query == "" {
		var students []models.InterfaceAccount
		cursor, err := accountCollection.Find(context.TODO(), bson.M{"role": "student"})
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tài khoản"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản"})
			return
		}
		if err := cursor.All(context.TODO(), &students); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":      "Tìm thấy tài khoản thành công",
			"foundedUser": students,
		})
	} else {
		var student models.InterfaceAccount
		err := accountCollection.FindOne(context.TODO(), bson.M{"ms": query}).Decode(&student)
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tài khoản"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":      "Tìm thấy tài khoản thành công",
			"foundedUser": student,
		})
	}
}

// HandleDeleteAccount xử lý việc xóa tài khoản.
func HandleDeleteAccount(c *gin.Context) {
	idParam := c.Param("id")
	accountID, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "ID không hợp lệ",
		})
		return
	}
	accountCollection := models.AccountModel()
	result, err := accountCollection.DeleteOne(context.TODO(), bson.M{"_id": accountID})
	if err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "Lỗi khi xóa tài khoản",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    "success",
		"message": "Xóa tài khoản thành công",
		"user":    result,
	})
}

// HandleUpdateAccount xử lý việc cập nhật thông tin tài khoản.
func HandleUpdateAccount(c *gin.Context) {
	idParam := c.Param("id")
	createdBy, _ := c.Get("ID")
	accountID, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "ID không hợp lệ",
		})
		return
	}
	var updatedAccount InterfaceAccount
	if err := c.ShouldBindJSON(&updatedAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu yêu cầu không hợp lệ"})
		return
	}
	updatedAccount.CreatedBy = createdBy
	accountCollection := models.AccountModel()
	if _, err := accountCollection.UpdateOne(context.TODO(), bson.M{"_id": accountID}, bson.M{"$set": updatedAccount}); err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "Lỗi khi cập nhật tài khoản",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    "success",
		"message": "Cập nhật tài khoản thành công",
	})
}
