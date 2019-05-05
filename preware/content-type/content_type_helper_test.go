package content_type_helper

import (
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

type widget struct {
	Name        string `json:"name" xml:"name"`
	ContentType string `json:"content_type" xml:"content_type"`
}

func app(name string) *buffalo.App {
	app := buffalo.New(buffalo.Options{
		PreWares: []buffalo.PreWare{AutoSetContentType},
	})

	r := render.New(render.Options{
		TemplatesBox: packr.New("HTML Templates", "./fixtures/test_templates"),
	})

	app.GET("/content", func(c buffalo.Context) error {
		ct, _ := c.Value("contentType").(string)
		w := &widget{Name: name, ContentType: ct}
		return c.Render(200, r.Auto(c, w))
	})
	return app
}

func Test_AutoSetContentType(t *testing.T) {
	name := "test_default"
	r := require.New(t)
	w := httptest.New(app(name))

	// Test normal html
	resHTML := w.HTML("/content").Get()
	r.Equal(200, resHTML.Code)
	ret := resHTML.Body.String()
	r.Equal("test html", ret)

	// Test extension html
	resHTML = w.HTML("/content.html").Get()
	r.Equal(200, resHTML.Code)
	ret = resHTML.Body.String()
	r.Equal("test html", ret)

	// Test normal json
	resJSON := w.JSON("/content").Get()
	r.Equal(200, resJSON.Code)
	ret = resJSON.Body.String()
	r.Equal("{\"name\":\"test_default\",\"content_type\":\"application/json\"}\n", ret)
	wid := &widget{}
	resJSON.Bind(wid)
	r.Equal(name, wid.Name)
	r.Equal("application/json", wid.ContentType)

	//Test extension json
	resHTML = w.HTML("/content.json").Get()
	r.Equal(200, resHTML.Code)
	ret = resHTML.Body.String()
	r.Equal("{\"name\":\"test_default\",\"content_type\":\"json\"}\n", ret)

	// Test normal xml
	resXML := w.XML("/content").Get()
	r.Equal(200, resXML.Code)
	ret = resXML.Body.String()
	r.Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<widget>\n  <name>test_default</name>\n  <content_type>application/xml</content_type>\n</widget>", ret)
	wid2 := &widget{}
	resXML.Bind(wid2)
	r.Equal(name, wid2.Name)
	r.Equal("application/xml", wid2.ContentType)

	//Test extension xml
	resHTML = w.HTML("/content.xml").Get()
	r.Equal(200, resHTML.Code)
	ret = resHTML.Body.String()
	r.Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<widget>\n  <name>test_default</name>\n  <content_type>xml</content_type>\n</widget>", ret)
}
