// -----------------------------------------
// go-ozzo/ozzo-validation
// -----------------------------------------
package validator2

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	regexpState = regexp.MustCompile("^[A-Z]{2}$")
	regexpZip   = regexp.MustCompile("^[0-9]{5}$")
)

func isEmail(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("must be email")
	}

	if err := is.Email.Validate(s); err != nil {
		return validation.NewError("validation_email_invalid", "must be email")
	}

	return nil
}

type User struct {
	FirstName string     `json:"first_name" validate:"required"`
	LastName  string     `json:"last_name" validate:"required"`
	Age       int32      `json:"age" validate:"gte=0,lte=130"`
	Email     string     `json:"email" validate:"required,email"`
	Addresses []*Address `json:"addresses" validate:"required,dive"`
	Gender    string     `json:"gender" validate:"oneof=female male"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Age, validation.Min(0), validation.Max(130)),
		// validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Email, validation.Required, validation.By(isEmail)),
		validation.Field(&u.Addresses, validation.Required),
		validation.Field(&u.Gender, validation.Required, validation.In("female", "male")),
	)
}

type Address struct {
	Street string `json:"street" validate:"required,min=5,max=50"`
	City   string `json:"city" validate:"required,min=5,max=50"`
	State  string `json:"state" validate:"required,custom_state_regexp"`
	Zip    string `json:"zip" validate:"required,custom_zip_regexp"`
}

func (a Address) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Street, validation.Required, validation.Length(5, 50)),
		validation.Field(&a.City, validation.Required, validation.Length(5, 50)),
		validation.Field(&a.State, validation.Required, validation.Match(regexpState)),
		validation.Field(&a.Zip, validation.Required, validation.Match(regexpZip)),
	)
}

type Period struct {
	Start     string `json:"start" validate:"datetime=2006-01-02T15:04:05Z07:00,custom_lt_datetime=End"`
	End       string `json:"end" validate:"datetime=2006-01-02T15:04:05Z07:00,custom_gt_datetime=Start"`
	UpdatedAt string `json:"updated_at" validate:"datetime=2006-01-02T15:04:05Z07:00"`
}

func (p Period) Validate() error {
	end, err := time.ParseInLocation(time.RFC3339, p.End, time.Local)
	if err != nil {
		return errors.New("invalid end")
	}

	start, err := time.ParseInLocation(time.RFC3339, p.Start, time.Local)
	if err != nil {
		return errors.New("invalid start")
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.Start, validation.Date(time.RFC3339).Max(end)),
		validation.Field(&p.End, validation.Date(time.RFC3339).Min(start)),
		validation.Field(&p.UpdatedAt, validation.Date(time.RFC3339)),
	)
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
			Email:     "moge.hogegmail.com",
			Addresses: []*Address{&a},
			Gender:    "mae",
		}
		p = Period{
			Start:     "2017-07-01T03:56:40+09:00",
			End:       "2017-06-01T03:56:40+09:00",
			UpdatedAt: "2017-07-01T03:56:40+09:00a",
		}
	)

	log.Println("============ go-ozzo/ozzo-validation ============")
	log.Println("======= Address =======")
	err := a.Validate()
	if err == nil {
		return
	}

	for key, err := range err.(validation.Errors) {
		log.Printf("[\"%s\", \"%s\", \"%s\"]\n", key, err, err.(validation.Error).Code())
	}

	log.Println("======= User =======")
	err = u.Validate()
	if err == nil {
		return
	}

	for key, errs := range err.(validation.Errors) {
		if _, ok := errs.(validation.Errors); !ok {
			if _, ok := errs.(validation.Error); ok {
				log.Printf("[\"%s\", \"%s\", \"%s\"]\n", key, errs, errs.(validation.Error).Code())
			}
			continue
		}
		// for nested struct
		for _, err := range errs.(validation.Errors) {
			for key, errs := range err.(validation.Errors) {
				log.Printf("[\"%s\", \"%s\", \"%s\"]\n", key, errs, errs.(validation.Error).Code())
			}
		}
	}

	log.Println("======= Period =======")
	err = p.Validate()
	if err == nil {
		return
	}

	for key, err := range err.(validation.Errors) {
		log.Printf("[\"%s\", \"%s\", \"%s\"]\n", key, err, err.(validation.Error).Code())
	}
}
