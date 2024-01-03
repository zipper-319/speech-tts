package jwtUtil

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Identifier struct{}

type TraceId struct{}

type Option func(*options)

var NotVerify = []string{
	"Register",
	"GetTtsConfig",
	"GetVersion",
	"GetSpeaker",
}

const (

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized(reason, "jwt token is missing")
	ErrMissingKey             = errors.Unauthorized(reason, "jwt key is missing")
	ErrTokenInvalid           = errors.Unauthorized(reason, "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized(reason, "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized(reason, "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized(reason, "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized(reason, "Wrong context for middleware")
	ErrNeedTokenProvider      = errors.Unauthorized(reason, "Token provider is missing")
	ErrSignToken              = errors.Unauthorized(reason, "Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.Unauthorized(reason, "Can not get key while signing token")
)

var keyFunc jwt.Keyfunc

var isOpenJwt bool

// Parser is a jwt parser
type options struct {
	signingMethod jwt.SigningMethod
	claims        jwt.Claims
	tokenHeader   map[string]interface{}
}

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithClaims with customer claim
// If you use it in Server, f needs to return a new jwt.Claims object each time to avoid concurrent write problems
// If you use it in Client, f only needs to return a single object to provide performance
func WithClaims(claims jwt.Claims) Option {
	return func(o *options) {
		o.claims = claims
	}
}

// WithTokenHeader withe customer tokenHeader for client side
func WithTokenHeader(header map[string]interface{}) Option {
	return func(o *options) {
		o.tokenHeader = header
	}
}

func Server(logger log.Logger, jwtKey string, isOpen bool, opts ...Option) middleware.Middleware {
	log.NewHelper(logger).Infof("jwtKey:%s, isOpenJwt:%t", jwtKey, isOpen)
	if jwtKey == "" {
		panic(ErrMissingKey)
	}
	keyFunc = KeyProvider(jwtKey)
	isOpenJwt = isOpen
	o := &options{
		signingMethod: jwt.SigningMethodHS256,
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {

				isNotVerify := false
				for _, operation := range NotVerify {
					if strings.Contains(header.Operation(), operation) {
						isNotVerify = true
						break
					}
				}
				if !isNotVerify {

					jwtToken := header.RequestHeader().Get(authorizationKey)
					log.NewHelper(logger).Infof("jwtToken:%s", jwtToken)
					if jwtToken == "" {
						return nil, ErrMissingJwtToken
					}

					var (
						tokenInfo *jwt.Token
						err       error
					)

					if o.claims != nil {
						tokenInfo, err = jwt.ParseWithClaims(jwtToken, o.claims, keyFunc)
					} else {
						tokenInfo, err = jwt.Parse(jwtToken, keyFunc)
					}
					if err != nil {
						ve, ok := err.(*jwt.ValidationError)
						if !ok {
							return nil, errors.Unauthorized(reason, err.Error())
						}
						if ve.Errors&jwt.ValidationErrorMalformed != 0 {
							return nil, ErrTokenInvalid
						}
						if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
							return nil, ErrTokenExpired
						}
						if ve.Inner != nil {
							return nil, ve.Inner
						}
						return nil, ErrTokenParseFail
					}
					if !tokenInfo.Valid {
						return nil, ErrTokenInvalid
					}
					if tokenInfo.Method != o.signingMethod {
						return nil, ErrUnSupportSigningMethod
					}
					ctx = context.WithValue(ctx, Identifier{}, tokenInfo)

				}
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}

	}
}

func IsValidity(logger log.Logger, header transport.Transporter) (*IdentityClaims, error) {

	for _, urlPath := range NotVerify {
		if strings.Contains(header.Operation(), urlPath) {
			return nil, nil
		}
	}
	jwtToken := header.RequestHeader().Get(authorizationKey)
	log.NewHelper(logger).Infof("jwtToken:%s; operation:%s", jwtToken, header.Operation())
	if !isOpenJwt && jwtToken == "" {
		return nil, nil
	}
	if jwtToken == "" {
		return nil, ErrMissingJwtToken
	}
	if keyFunc == nil {
		return nil, ErrMissingKey
	}
	identifier := IdentityClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, &identifier, keyFunc)
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if !ok {
			return nil, errors.Unauthorized(reason, err.Error())
		}
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, ErrTokenInvalid
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, ErrTokenExpired
		}
		if ve.Inner != nil {
			return nil, ve.Inner
		}
		return nil, ErrTokenParseFail
	}
	log.NewHelper(logger).Infof("account:%s", identifier.Account)
	return &identifier, nil

}

func GenerateJwtToken(kf jwt.Keyfunc, opts ...Option) (string, error) {
	o := &options{
		signingMethod: jwt.SigningMethodHS256,
	}
	for _, opt := range opts {
		opt(o)
	}
	token := jwt.NewWithClaims(o.signingMethod, o.claims)
	if o.tokenHeader != nil {
		for k, v := range o.tokenHeader {
			token.Header[k] = v
		}
	}
	key, err := kf(token)
	if err != nil {
		return "", err
	}

	if tokenStr, err := token.SignedString(key); err != nil {
		return "", err
	} else {
		return tokenStr, nil
	}
}

func KeyProvider(key string) jwt.Keyfunc {
	return func(*jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}
}

type IdentityClaims struct {
	NameId  int
	Account string
	Role    string
	jwt.RegisteredClaims
}

func GetToken(identifier string, expire int, key string) (string, error) {
	now := time.Now()
	id, _ := uuid.NewUUID()
	claims := IdentityClaims{
		Account: identifier,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{identifier},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        id.String(),
		},
	}
	return GenerateJwtToken(KeyProvider(key), WithClaims(claims))
}

func ParseToken(token, key string) (*jwt.Token, error) {
	identifier := IdentityClaims{}
	tokenInfo, err := jwt.ParseWithClaims(token, &identifier, KeyProvider(key))
	if err != nil {
		return nil, err
	}
	fmt.Println(identifier.Account)
	return tokenInfo, nil
}
