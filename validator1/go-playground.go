// -----------------------------------------
// go-playground/validator
// -----------------------------------------
package validator1

import (
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

var (
	regexpState = regexp.MustCompile("^[A-Z]{2}$")
	regexpZip   = regexp.MustCompile("^[0-9]{5}$")
)

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("custom_state_regexp", state)
	_ = validate.RegisterValidation("custom_zip_regexp", zip)
	_ = validate.RegisterValidation("custom_gt_datetime", gtDateTime)
	_ = validate.RegisterValidation("custom_lt_datetime", ltDateTime)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	})
}

func state(fl validator.FieldLevel) bool {
	return regexpState.MatchString(fl.Field().String())
}

func zip(fl validator.FieldLevel) bool {
	return regexpZip.MatchString(fl.Field().String())
}

func ltDateTime(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, _, ok := fl.GetStructFieldOK2()
	if !ok || currentKind != kind {
		return true
	}

	t, err := time.ParseInLocation(time.RFC3339, currentField.Interface().(string), time.Local)
	if err != nil {
		return true
	}

	fieldTime, err := time.ParseInLocation(time.RFC3339, field.Interface().(string), time.Local)
	if err != nil {
		return true
	}

	return fieldTime.Before(t)
}

func gtDateTime(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, _, ok := fl.GetStructFieldOK2()
	if !ok || currentKind != kind {
		return true
	}

	t, err := time.ParseInLocation(time.RFC3339, currentField.Interface().(string), time.Local)
	if err != nil {
		return true
	}

	fieldTime, err := time.ParseInLocation(time.RFC3339, field.Interface().(string), time.Local)
	if err != nil {
		return true
	}

	return fieldTime.After(t)
}

type User struct {
	FirstName string     `json:"first_name" validate:"required"`
	LastName  string     `json:"last_name" validate:"required"`
	Age       int32      `json:"age" validate:"gte=0,lte=130"`
	Email     string     `json:"email" validate:"required,email"`
	Addresses []*Address `json:"addresses" validate:"required,dive"`
	Gender    string     `json:"gender" validate:"oneof=female male"`
}

type Address struct {
	Street string `json:"street" validate:"required,min=5,max=50"`
	City   string `json:"city" validate:"required,min=5,max=50"`
	State  string `json:"state" validate:"required,custom_state_regexp"`
	Zip    string `json:"zip" validate:"required,custom_zip_regexp"`
}

type Period struct {
	Start     string `json:"start" validate:"datetime=2006-01-02T15:04:05Z07:00,custom_lt_datetime=End"`
	End       string `json:"end" validate:"datetime=2006-01-02T15:04:05Z07:00,custom_gt_datetime=Start"`
	UpdatedAt string `json:"updated_at" validate:"datetime=2006-01-02T15:04:05Z07:00"`
}

func Validate() {
	var (
		a = Address{
			Street: "123",
			City:   "Unknown",
			State:  "Virginia",
			Zip:    "12345",
		}
		u = User{
			FirstName: "moge",
			LastName:  "hoge",
			Age:       135,
			Email:     "moge.hoge.com",
			Addresses: []*Address{&a},
			Gender:    "mae",
		}
		p = Period{
			Start:     "2017-07-01T03:56:40+09:00",
			End:       "2017-06-01T03:56:40+09:00",
			UpdatedAt: "2017-07-01T03:56:40+09:00a",
		}
	)

	log.Println("============ go-playground/validator ============")
	log.Println("======= Address =======")
	err := validate.Struct(a)
	if err == nil {
		return
	}

	for _, err := range err.(validator.ValidationErrors) {
		log.Printf("[\"%s\", \"%s\", \"%s\"]\n", err.Field(), err, err.ActualTag())
	}

	log.Println("======= User =======")
	err = validate.Struct(u)
	if err == nil {
		return
	}

	for _, err := range err.(validator.ValidationErrors) {
		log.Printf("[\"%s\", \"%s\", \"%s\"]\n", err.Field(), err, err.ActualTag())
	}

	log.Println("======= Period =======")
	err = validate.Struct(p)
	if err == nil {
		return
	}

	for _, err := range err.(validator.ValidationErrors) {
		log.Printf("[\"%s\", \"%s\", \"%s\"]\n", err.Field(), err, err.ActualTag())
	}
}
