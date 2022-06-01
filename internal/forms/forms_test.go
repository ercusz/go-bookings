package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestValid(t *testing.T) {
	r := httptest.NewRequest("POST", "/abc", nil)
	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("got an invalid form when it should be valid")
	}
}

func TestRequired(t *testing.T) {
	r := httptest.NewRequest("POST", "/abc", nil)
	form := New(r.PostForm)

	// add required fields and check the form 
	// that has required fields is contains values
	// in required fields or not
	form.Required("first_name", "last_name", "email")
	if form.Valid() {
		t.Error("got a valid form even if required fields are missing")
	}

	// create new request
	r, _ = http.NewRequest("POST", "/abc", nil)
	// added the fields with values into request form
	data := url.Values{}
	data.Add("username", "John")
	data.Add("email", "John@Doe.com")
	r.PostForm = data
	// create new form with request form
	form = New(r.PostForm)
	// added required fields
	form.Required("username","email")
	// check the form that contains values of required fields is invalid or not
	if !form.Valid() {
		t.Error("got an invalid form even though required fields were added")
	}
}

func TestHas(t *testing.T) {
	r := httptest.NewRequest("POST", "/abc", nil)
	form := New(r.PostForm)
	
	// check if the form without any fields has "first_name" field or not
	if form.Has("first_name") {
		t.Error("got field from the form that has no field")
	}
	
	data := url.Values{}
	data.Add("username", "John")
	// create new form with data
	form = New(data)
	// check if the form doesn't have a "username" field
	if !form.Has("username") {
		t.Error("field not found but that field should be exist in the form")
	}
}

func TestMinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/abc", nil)
	// create form without any fields
	form := New(r.PostForm)
	// check if "password" field value is reached the minimum specified length
	form.MinLength("password", 8)
	if form.Valid() {
		t.Error("got valid when found the min length of non-exist field")
	}

	isError := form.Errors.Get("password")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}
	
	data := url.Values{}
	data.Add("password", "123456")
	// create new form with data
	form = New(data)
	// check if "password" field value is reached the minimum specified length
	form.MinLength("password", 8)
	if form.Valid() {
		t.Error("got valid when value length shorter than min length")
	}

	data.Set("password", "john12345678")
	// create new form with data
	form = New(data)
	// check if "password" field value is reached the minimum specified length
	form.MinLength("password", 8)
	if !form.Valid() {
		t.Error("got invalid even though field value length greater than min length")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestIsEmail(t *testing.T) {

	data := url.Values{}
	form := New(data)
	
	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid when found non-exist field")
	}
	
	// add data and assign to new form
	data.Add("email", "123456")
	form = New(data)
	// check if the field value is valid email format
	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid when field value is not valid email format")
	}
	
	// update data and assign to new form
	data.Set("email", "John@Doe.com")
	form = New(data)
	// check if the field value is valid email format
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got invalid even though the field value is in email format")
	}
}
