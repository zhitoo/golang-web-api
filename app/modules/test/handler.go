package test

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/response"
)

var List []int = []int{}

type Handler struct{}

func NewHandler() *Handler {
	return new(Handler)
}

func (h *Handler) List(c *gin.Context) {
	for i := 0; i < 99999999999; i++ {
		log.Println(i)
	}
	c.String(http.StatusOK, "list")
}

func (h *Handler) Show(c *gin.Context) {
	c.String(http.StatusOK, "show: "+c.Param("id"))
}

func (h *Handler) Store(c *gin.Context) {
	//form := new(forms.CreateTestForm)
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "id must be an integer")
		return
	}

	List = append(List, id)
	log.Println(List)
	c.String(http.StatusOK, "store: "+strconv.Itoa(id))
}

func (h *Handler) HeaderBinder1(c *gin.Context) {
	userId := c.GetHeader("UserId")

	c.JSON(http.StatusOK, gin.H{
		"result": "HeaderBinder1",
		"UserId": userId,
	})

}

func (h *Handler) HeaderBinder2(c *gin.Context) {

	heads := new(struct {
		UserId  string
		Browser string
	})
	c.BindHeader(heads)

	c.JSON(http.StatusOK, gin.H{
		"result": "HeaderBinder2",
		"header": heads,
	})

}

func (t *Handler) QueryBinder1(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")

	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": name,
	})

}

func (t *Handler) QueryBinder2(c *gin.Context) {
	id := c.QueryArray("id")
	name := c.Query("name")

	c.JSON(http.StatusOK, gin.H{
		"ids":  id,
		"name": name,
	})

}

func (t *Handler) UriBinder(c *gin.Context) {
	id := c.Param("id")
	name := c.Param("name")

	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": name,
	})

}

type user struct {
	ID     int    `json:"id" binding:"required,numeric,min=1"`
	Name   string `json:"name" binding:"required,alpha,min=3"`
	Mobile string `json:"mobile" binding:"irmobile"`
}

// BodyBinder
// @Summary BodyBinder
// @Description BodyBinder
// @Tags test
// @Accept json
// @Produce json
// @Param user body user true "user data"
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Success 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/test/body-binder [post]
func (t *Handler) BodyBinder(c *gin.Context) {
	u := new(user)
	err := c.ShouldBindJSON(u)

	if err != nil {
		response.NewReponse().SetResult(nil).SetError(err).SetStatus(false).SetHttpStatusCode(422).Json(c)
		return
	}

	response.NewReponse().SetResult(gin.H{
		"id":     u.ID,
		"name":   u.Name,
		"mobile": u.Mobile,
	}).Json(c)
}

func (t *Handler) FileBinder(c *gin.Context) {
	file, _ := c.FormFile("file")

	err := c.SaveUploadedFile(file, "./assets/"+file.Filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"file": file.Filename,
	})
}
