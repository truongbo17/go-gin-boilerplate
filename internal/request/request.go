package request

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-playground/validator/v10"
	json "github.com/json-iterator/go"
	"io"
	"net/url"
	"reflect"
	"strconv"
)

type BaseRequestPaginate struct {
	Page    int `form:"page" json:"page"`
	PerPage int `form:"per_page" json:"per_page"`
}

func SetDefaults(obj interface{}) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if value, ok := fld.Tag.Lookup("default"); ok {
			return value
		}
		return ""
	})

	return validate.Struct(obj)
}

type FlexibleString string

func (fs *FlexibleString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*fs = FlexibleString(str)
		return nil
	}

	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*fs = FlexibleString(strconv.Itoa(num))
		return nil
	}

	return fmt.Errorf("must be a string")
}

func ExtractAllData(c *gin.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, values := range c.Request.URL.Query() {
		if len(values) > 1 {
			result[key] = values
		} else {
			result[key] = values[0]
		}
	}

	contentType := c.ContentType()
	if contentType == gin.MIMEPOSTForm {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return result, fmt.Errorf("error reading JSON body: %v", err)
		}
		//defer func(Body io.ReadCloser) {
		//	_ = Body.Close()
		//}(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		parsed, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		for key, values := range parsed {
			if len(values) > 1 {
				result[key] = values
			} else {
				result[key] = values[0]
			}
		}
	}

	if contentType == gin.MIMEMultipartPOSTForm {
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			return result, fmt.Errorf("error parsing form data: %v", err)
		}
		for key, values := range c.Request.PostForm {
			if len(values) > 1 {
				result[key] = values
			} else {
				result[key] = values[0]
			}
		}

		form, err := c.MultipartForm()
		if err != nil {
			return result, fmt.Errorf("error parsing multipart form: %v", err)
		}
		for key, files := range form.File {
			fileNames := make([]string, len(files))
			for i, file := range files {
				fileNames[i] = file.Filename
			}
			result[key] = fileNames
		}
	}

	if contentType == gin.MIMEJSON {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return result, fmt.Errorf("error reading JSON body: %v", err)
		}
		//defer func(Body io.ReadCloser) {
		//	_ = Body.Close()
		//}(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if len(body) == 0 {
			return result, nil
		}

		var jsonData map[string]interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			return result, fmt.Errorf("error unmarshalling JSON: %v", err)
		}

		for key, value := range jsonData {
			result[key] = value
		}
	}

	return result, nil
}

func Nested(target interface{}, fieldRules ...*validation.FieldRules) *validation.FieldRules {
	return validation.Field(target, validation.By(func(value interface{}) error {
		valueV := reflect.Indirect(reflect.ValueOf(value))
		if valueV.CanAddr() {
			addr := valueV.Addr().Interface()
			return validation.ValidateStruct(addr, fieldRules...)
		}
		return validation.ValidateStruct(target, fieldRules...)
	}))
}
