package domains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alserov/rently/api/internal/log"
	"github.com/alserov/rently/api/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func handleResponseError(err error) {
	if err != nil {
		log.GetLogger().Error("failed to send response", slog.String("error", err.Error()))
	}
}

func marshal(in interface{}) []byte {
	b, err := json.Marshal(in)
	if err != nil {
		log.GetLogger().Error("failed to marshal response", slog.String("error", err.Error()))
		return nil
	}

	return b
}

func decode(b []byte, target any, valid *validator.Validate) error {
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(target); err != nil {
		return fmt.Errorf("failed to decode req body: %v", err)
	}

	if err := valid.Struct(target); err != nil {
		return fmt.Errorf("invalid data: %s", err.Error())
	}

	return nil
}

func parseForm(c *fiber.Ctx, target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Struct {
		panic("unexpected type: " + val.Kind().String())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return fmt.Errorf("failed to open form: %v", err)
	}

	for key, v := range form.File {
		field, ok := val.Type().FieldByName(key)
		if !ok {
			return fmt.Errorf("invalid form parameter: %s", key)
		}

		if t := field.Tag.Get("type"); t == "main" {
			file, err := v[0].Open()
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}

			b, err := io.ReadAll(file)
			if err != nil {
				return fmt.Errorf("failed to read file: %v", err)
			}

			val.FieldByName(key).Set(reflect.ValueOf(b))
			continue
		}

		files := reflect.MakeSlice(reflect.TypeOf([][]uint8{}), len(v), len(v))
		for idx, file := range v {
			file, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}

			b, err := io.ReadAll(file)
			if err != nil {
				return fmt.Errorf("failed to read file: %v", err)
			}

			files.Index(idx).Set(reflect.ValueOf(b))
		}
		val.FieldByName(key).Set(files)
	}

	for key, v := range form.Value {
		f := val.FieldByName(key)
		if f == (reflect.Value{}) {
			return fmt.Errorf("field %s does not exist", key)
		}

		switch val.FieldByName(key).Interface().(type) {
		case int32:
			value, err := strconv.Atoi(v[0])
			if err != nil {
				panic("invalid type cast: " + err.Error())
			}
			f.SetInt(int64(value))
		case string:
			f.SetString(v[0])
		case float32:
			value, err := strconv.ParseFloat(v[0], 32)
			if err != nil {
				panic("invalid type cast: " + err.Error())
			}
			f.SetFloat(value)
		}
	}

	return nil
}

// parses query params form path and maps them to structure fields by json tag
func parseQueryParams(c *fiber.Ctx, target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Struct {
		panic("unexpected type: " + val.Kind().String())
	}

	params := c.Queries()
	if len(params) == 0 {
		return fmt.Errorf("no parameters provided")
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		tag := strings.Split(val.Type().Field(i).Tag.Get("json"), ",")[0]

		if _, ok := params[tag]; ok {
			switch f.Interface().(type) {
			case int:
				value, err := strconv.Atoi(params[tag])
				if err != nil {
					return fmt.Errorf("invalid parameter type: %s", tag)
				}
				f.SetInt(int64(value))
			case string:
				f.SetString(params[tag])
			case float32:
				value, err := strconv.ParseFloat(params[tag], 32)
				if err != nil {
					return fmt.Errorf("invalid parameter type: %s", tag)
				}
				f.SetFloat(value)
			}
		}
	}

	return nil
}

func transformImageInfoToLink(bucket string, id string) string {
	return fmt.Sprintf("%s/%s/%s", path, bucket, id)
}

// Handles error from gRPC user
func handleServiceError(w *fasthttp.Response, err error) {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Internal:
			w.SetStatusCode(http.StatusInternalServerError)
		case codes.NotFound:
			w.SetStatusCode(http.StatusNotFound)
			w.SetBody(marshal(models.Error{Err: st.Message()}))
		case codes.InvalidArgument:
			w.SetStatusCode(http.StatusBadRequest)
			w.SetBody(marshal(models.Error{
				Err: st.Message(),
			}))
		default:
			log.GetLogger().Error("unknown service error", slog.String("error", err.Error()))
			w.SetStatusCode(http.StatusInternalServerError)
		}
	}
}
