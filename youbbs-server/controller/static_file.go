package controller

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func (h *BaseHandler) StaticFile(c *gin.Context) {
	filePath := c.Param("filepath")
	buf, err := ioutil.ReadFile("static/" + filePath)
	if err != nil {
		c.Header("Content-Type", "text/plain; charset=utf-8")
		//w.WriteHeader(http.StatusNotFound)
		//_, _ = w.Write([]byte(err.Error()))
		return
	}
	fileType := http.DetectContentType(buf)
	c.Header("Content-Type", fileType)
	//_, _ = w.Write(buf)
}
