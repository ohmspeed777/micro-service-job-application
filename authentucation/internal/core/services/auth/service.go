package auth

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	refreshtoken "app/repository/refreshToken"
	"app/repository/user"
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/ohmspeed777/go-pkg/errorx"
	"github.com/ohmspeed777/go-pkg/logx"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	SignUp(ctx context.Context, req SignUpReq) (*domain.CommonResponse, error)
	SignIn(ctx context.Context, req SignInReq) (*SignInRes, error)
}

type Service struct {
	user ports.IUserRepo
	rt   ports.IRefreshTokenRepo
}

func NewService(db *mongo.Database) IAuthService {
	return &Service{
		user: user.NewRepo(db),
		rt:   refreshtoken.NewRepo(db),
	}
}

func (s *Service) SignUp(ctx context.Context, req SignUpReq) (*domain.CommonResponse, error) {
	entity := &domain.User{}
	copier.Copy(entity, &req)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not encrypt password", err)
	}

	entity.Password = string(hash)
	err = s.user.Create(entity)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not create user", err)
	}

	return domain.NewOkMessage(), nil
}

func (s *Service) SignIn(ctx context.Context, req SignInReq) (*SignInRes, error) {
	user, err := s.user.FindOneByEmail(req.Email)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "invalid email or password", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errorx.New(http.StatusBadRequest, "invalid email or password", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(viper.GetString("jwt.key")))
	if err != nil {
		logx.GetLog().Errorf("jwt getting key was error: %v", err)
		return nil, errorx.New(http.StatusBadRequest, "jwt getting key was error", err)
	}

	uuidToken := uuid.NewV4().String()
	uuidRTToken := uuid.NewV4().String()

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "access_token"
	claims["iss"] = "app"
	claims["jti"] = uuidToken
	claims["iat"] = time.Now().Local().Unix()
	claims["id"] = user.ID.Hex()
	claims["role"] = "users"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	accessToken, err := token.SignedString(key)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not generate access token", err)
	}

	// create refresh token
	rtToken := jwt.New(jwt.SigningMethodRS256)
	rtClaims := rtToken.Claims.(jwt.MapClaims)

	rtClaims["id"] = user.ID.Hex()
	rtClaims["sub"] = "refresh_token"
	rtClaims["iss"] = "app"
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rtClaims["jti"] = uuidRTToken

	refreshToken, err := rtToken.SignedString(key)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not generate refresh token", err)
	}

	err = s.rt.Create(&domain.RefreshToken{
		UUID:   refreshToken,
		Type:   domain.REFRESH_TOKEN,
		UserID: user.ID,
	})
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not store refresh token", err)
	}


	err = s.rt.Create(&domain.RefreshToken{
		UUID:   refreshToken,
		Type:   domain.ACCESS_TOKEN,
		UserID: user.ID,
	})
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not store access token", err)
	}

	return &SignInRes{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
