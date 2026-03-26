package scorer

import (
	"github.com/gin-gonic/gin"
)

func PostScoreHanlder(c *gin.Context) {
	model := c.Param("model")

	_, exists := GetConfig().Models[model]

	if !exists {
		c.JSON(400, gin.H{
			"error": "Invalid input model",
		})
		return
	}

	// score, err := ValueAlignmentTest(model, "W IL COMUINISMO")
	// if err != nil {
	// 	panic(err)
	// }

	score := 0
	c.JSON(200, gin.H{
		"model": model,
		"score": score,
	})
}
