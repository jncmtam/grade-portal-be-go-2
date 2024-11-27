package controller_admin

import (
	"Go2/models"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// HandleCreateResult xử lý việc tạo bảng điểm mới.
func HandleCreateResult(c *gin.Context) {
	var data InterfaceResultScoreController
	// Lấy dữ liệu từ front end
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "Dữ liệu không hợp lệ",
		})
		return
	}
	classID, err := bson.ObjectIDFromHex(data.ClassID)
	if err != nil {
		c.JSON(204, gin.H{
			"code":    "error",
			"message": "Lớp chưa có giáo viên",
		})
		return
	}
	createdBy, _ := c.Get("ID")
	collection := models.ResultScoreModel()
	var result models.InterfaceResult
	err = collection.FindOne(context.TODO(), bson.M{"class_id": classID}).Decode(&result)
	// Có bản ghi result trước đó
	if err == nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "Bảng ghi của lớp học này đã được lưu trong database trước đó",
		})
		return
	}

	var classDetail models.InterfaceClass
	collectionClass := models.ClassModel()
	if err = collectionClass.FindOne(context.TODO(), bson.M{"_id": classID}).Decode(&classDetail); err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "Không tìm thấy lớp học đó",
		})
		return
	}
	collection.InsertOne(context.TODO(), bson.M{
		"semester":  classDetail.Semester,
		"course_id": classDetail.CourseId,
		"score":     data.SCORE,
		"class_id":  classID,
		"expiredAt": time.Now().AddDate(0, 6, 0),
		"createdBy": createdBy,
		"updatedBy": createdBy,
	})
	c.JSON(200, gin.H{
		"code":    "success",
		"message": "Cập nhật bảng điểm thành công",
	})
}

// HandleGetResult xử lý việc lấy bảng điểm theo ID lớp học.
func HandleGetResult(c *gin.Context) {
	param := c.Param("id")
	classID, err := bson.ObjectIDFromHex(param)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "ID không hợp lệ",
		})
		return
	}
	collection := models.ResultScoreModel()
	var data models.InterfaceResult
	if err = collection.FindOne(context.TODO(), bson.M{"class_id": classID}).Decode(&data); err != nil {
		c.JSON(400, gin.H{
			"code":    "error",
			"message": "Không tìm thấy bảng điểm",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    "success",
		"message": "Lấy bảng điểm thành công",
		"score":   data,
	})
}
