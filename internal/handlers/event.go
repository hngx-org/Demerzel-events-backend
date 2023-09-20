package handlers

import (
	//"demerzel-events/internal/models"
	//"demerzel-events/internal/db"
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"

	"github.com/gin-gonic/gin"
)

func GroupEvent(c *gin.Context) {

	id := c.Param("id")
	//group, err := models.Group.Get(id)
	// if err != nil{

	// 	c.JSON(404, gin.H{
	// 		"status":  "Failed",
	// 		"message": "Group Not Found",
	// 	})
	// }
	
    //result, err := group.GetGroupEvent(db.DB)

	//this is for testing
    group := models.Group{
        Id: id,	
	}

	result:= group.GetGroupEvent(db.DB)	
	
		c.JSON(200, gin.H{
			"status":  "success",
			"message": result,
		})	
}