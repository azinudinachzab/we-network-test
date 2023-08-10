package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) UserRegistration(ctx echo.Context) error {
	var (
		req generated.RegisterRequest
		res generated.RegisterResponse
	)
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	if req.Password == "" || req.FullName == "" || req.PhoneNumber == "" {
		errMsg := "field full_name, password and phone_number cannot be empty"
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	err = validateRegistrationRequest(req)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	saltedAndHashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	ID := uint64(s.Snowflake.Generate())
	if err := s.Repository.InsertUser(ctx.Request().Context(), repository.User{
		ID:          ID,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Password:    string(saltedAndHashed),
	}); err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	IDint64 := int64(ID)
	res.Id = &IDint64
	return ctx.JSON(http.StatusOK, res)
}

func validateRegistrationRequest(req generated.RegisterRequest) error {
	err := validatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return err
	}
	err = validateFullName(req.FullName)
	if err != nil {
		return err
	}
	err = validatePassword(req.Password)
	if err != nil {
		return err
	}
	return nil
}

func validateUpdateProfileRequest(req generated.UpdateProfileRequest) error {
	err := validatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return err
	}
	err = validateFullName(req.FullName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) UserLogin(ctx echo.Context) error {
	var (
		req generated.LoginRequest
		res generated.LoginResponse
	)
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	if req.Password == "" || req.PhoneNumber == "" {
		errMsg := "field password and phone_number cannot be empty"
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	usr, err := s.Repository.GetUserByPhone(ctx.Request().Context(), req.PhoneNumber)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password)); err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	oneDay := 24*time.Hour
	token, err := s.JWTCreate(oneDay, usr.ID)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	if err := s.Repository.UpdateUserRecord(ctx.Request().Context(), usr.ID); err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	IDint64 := int64(usr.ID)
	res.Id = &IDint64
	res.Token = &token
	return ctx.JSON(http.StatusOK, res)
}

func (s *Server) UserGetProfile(ctx echo.Context, params generated.UserGetProfileParams) error {
	data, err := s.JWTVal(params.Authorization)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusForbidden, errResp)
	}
	id, err := convert(data)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusForbidden, errResp)
	}

	usr, err := s.Repository.GetUserByID(ctx.Request().Context(), uint64(id))
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	return ctx.JSON(http.StatusOK, generated.GetProfileResponse{
		FullName:    &usr.FullName,
		PhoneNumber: &usr.PhoneNumber,
	})
}

func (s *Server) UserUpdateProfile(ctx echo.Context, params generated.UserUpdateProfileParams) error {
	data, err := s.JWTVal(params.Authorization)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusForbidden, errResp)
	}
	id, err := convert(data)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusForbidden, errResp)
	}

	var (
		req generated.UpdateProfileRequest
	)
	err = json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	if req.FullName == "" || req.PhoneNumber == "" {
		errMsg := "field full_name and phone_number cannot be empty"
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	err = validateUpdateProfileRequest(req)
	if err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	_, err = s.Repository.GetUserByPhone(ctx.Request().Context(), req.PhoneNumber)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			errMsg := err.Error()
			errResp := generated.ErrorResponse{
				Message: &errMsg,
			}
			return ctx.JSON(http.StatusBadRequest, errResp)
		}
	}

	if err := s.Repository.UpdateProfileByID(ctx.Request().Context(), uint64(id), req.FullName, req.PhoneNumber); err != nil {
		errMsg := err.Error()
		errResp := generated.ErrorResponse{
			Message: &errMsg,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	msg := "Update profile success"
	return ctx.JSON(http.StatusOK, generated.UpdateProfileResponse{
		Message: &msg,
	})
}

func validatePhoneNumber(p string) error {
	if len(p) < 10 || len(p) > 13 {
		return errors.New("phone number cannot less than 10 char or more than 13 char")
	}

	firstThreeChar := p[0:3]
	if firstThreeChar != "+62" {
		return errors.New("phone number should start with '+62'")
	}

	return nil
}

func validateFullName(fn string) error {
	if len(fn) < 3 || len(fn) > 60 {
		return errors.New("full name cannot less than 3 char or more than 60 char")
	}

	return nil
}

func validatePassword(p string) error {
	if len(p) < 6 || len(p) > 64 {
		return errors.New("password cannot less than 6 char or more than 64 char")
	}

	var (
		isUpperExist   = false
		isNumberExist  = false
		isSpecialExist = false
	)

	for _, v := range p {
		if unicode.IsUpper(v) {
			isUpperExist = true
			continue
		}

		if unicode.IsNumber(v) {
			isNumberExist = true
			continue
		}

		if !unicode.IsLetter(v) && !unicode.IsNumber(v) {
			isSpecialExist = true
			continue
		}
	}

	if !isSpecialExist || !isNumberExist || !isUpperExist {
		return errors.New("password wrong format")
	}

	return nil
}

func convert(t interface{}) (int64, error) {
	switch t := t.(type) {   // This is a type switch.
	case int64:
		return t, nil        // All done if we got an int64.
	case int:
		return int64(t), nil // This uses a conversion from int to int64
	case string:
		return strconv.ParseInt(t, 10, 64)
	case float64:
		return int64(t), nil
	default:
		return 0, fmt.Errorf("type %T not supported", t)
	}
}
