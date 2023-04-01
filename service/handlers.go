package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"wd_users/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GetUser(c *gin.Context) {
	l := s.l.Named("GetUser")

	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error getting user", zap.Error(err)) //error message shouldn't contain single quote(') cause it might break. Spacebar is okay.
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUser(context.Background(), uid)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			l.Info("error getting user", zap.Error(err))
			c.Status(http.StatusNotFound)
			return
		}
		l.Error("error getting user", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbUser2ApiUser(user))
}

func (s *Service) CreateUser(c *gin.Context) {
	l := s.l.Named("CreateUser")

	var createUserRequest struct { //embedding User struct
		User
		Password string `json:"password,omitempty"` //you only need to use this field once at this one spot(CreateUser)
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&createUserRequest); err != nil {
		l.Info("error creating user", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidCreateUserRequest(createUserRequest.User, createUserRequest.Password) {
		l.Info("error creating user")
		c.Status(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), bcrypt.MinCost)
	if err != nil {
		l.Error("error generating hashed password", zap.Error(err)) //"unexpected error ..."는 테스트 할 때 써라.
		c.Status(http.StatusInternalServerError)
		return
	}

	u := apiUser2DBUser(createUserRequest.User)
	u.HashedPassword = string(hashedPassword)

	user, err := s.db.CreateUser(context.Background(), u) //이렇게 에러 처리해놓으면 굳이 db에서 *store.User를 리턴할 필요는 없지만, 이 에러처리를 까먹는 개발자도 있음..
	if err != nil {
		l.Error("error creating user", zap.Error(err)) //error creating user, user.UUID(create 실패한 user인데 이 user의 UUID를 에러 메시지로 반환하려고;;) 이렇게 zero value 필드에 접근하는 실수를 하는 개발자도 있다고..;; 그래서 위를 포함한 이러한 상황에 대비해 *store.User로 리턴하는 것임..
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbUser2ApiUser(user))
}

func (s *Service) UpdateUser(c *gin.Context) {
	l := s.l.Named("UpdateUser")

	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error updating user", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	var updateUserRequest User

	if err := json.NewDecoder(c.Request.Body).Decode(&updateUserRequest); err != nil {
		l.Info("error updating user", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidUpdateUserRequest(updateUserRequest, uid) {
		l.Info("error updating user")
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := s.db.UpdateUser(context.Background(), apiUser2DBUser(updateUserRequest)) //if I have two variables, I can still do combined if statement, like if user, err ... ; err != nil {}, but then user can only survive within the next 3 lines of if statement. So I can't return user variable at the bottom in c.JSON().
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			l.Info("error updating user", zap.Error(err))
			c.Status(http.StatusNotFound)
			return
		}
		l.Error("error updating user", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbUser2ApiUser(user))
}
