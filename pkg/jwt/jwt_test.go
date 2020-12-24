package jwt

import (
	"fmt"
	"os"
	"testing"

	"github.com/mhdiiilham/oauth2-auth-server-implementation/entity/user"
	"github.com/mhdiiilham/oauth2-auth-server-implementation/pkg/password"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var entity user.User
var service TokenService
var token string
var jwtToken *jwt.Token
var rsa256 string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJSUzI1NmluT1RBIiwibmFtZSI6IkpvaG4gRG9lIn0.ICV6gy7CDKPHMGJxV80nDZ7Vxe0ciqyzXD_Hr4mTDrdTyi6fNleYAyhEZq2J29HSI5bhWnJyOBzg2bssBUKMYlC2Sr8WFUas5MAKIr2Uh_tZHDsrCxggQuaHpF4aGCFZ1Qc0rrDXvKLuk1Kzrfw1bQbqH6xTmg2kWQuSGuTlbTbDhyhRfu1WDs-Ju9XnZV-FBRgHJDdTARq1b4kuONgBP430wJmJ6s9yl3POkHIdgV-Bwlo6aZluophoo5XWPEHQIpCCgDm3-kTN_uIZMOHs2KRdb6Px-VN19A5BYDXlUBFOo-GvkCBZCgmGGTlHF_cWlDnoA9XTWWcIYNyUI4PXNw"
var anotherIssuer string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2MDQwNjU5MDEsImV4cCI6MTYzNTYwMTkwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.wR0L7T9T2hzfOlODGG8K4Anr4YX_E7kgPOe8so8ejLw2N4sF-jhsnf3crP20CH4YwvGJgnFAutalvbiWrkVUqg"

func init() {
	user := user.User{
		Fullname: "Muhammad Ilham",
		Email:    "example@mail.com",
		Password: password.Hash("HelloWorld"),
	}
	entity = user

	secret := fmt.Sprintf("%s", os.Getenv("JWT_SECRET"))
	issuer := fmt.Sprintf("%s", os.Getenv("APP_NAME"))
	service = NewJWTService(secret, issuer)
}

func TestGenerateToken(t *testing.T) {
	token = service.Generate(&entity)
	valueType := fmt.Sprintf("%T", token)
	assert.NotEqual(t, "", token)
	assert.Equal(t, "string", valueType)
}

func TestValidateToken(t *testing.T) {
	verifiedToken, err := service.Verify(token)
	jwtToken = verifiedToken
	assert.Equal(t, nil, err)
}

func TestVerifyNonHMACSigned(t *testing.T) {
	_, err := service.Verify(rsa256)
	assert.NotEqual(t, nil, err)
}

func TestVerifyTokenGeneratedByAnotherIssuer(t *testing.T) {
	_, err := service.Verify(anotherIssuer)
	assert.NotEqual(t, nil, err)
}

func TestIfTokenStillValid(t *testing.T) {
	err := service.Validate(token)
	assert.Equal(t, nil, err)
}

func TestValidateTokenSignedByAnotherAlgorithm(t *testing.T) {
	err := service.Validate(rsa256)
	assert.NotEqual(t, nil, err)
}

func testGetJWTIssuer(t *testing.T) {
	issuer := service.GetIssuer()
	assert.Equal(t, os.Getenv("APP_NAME"), issuer)
}

func TestExtractTokenMetaData(t *testing.T) {
	metaData, err := service.Extract(token)

	assert.Equal(t, nil, err)
	assert.Equal(t, entity.Email, metaData.Email)
	assert.Equal(t, service.GetIssuer(), metaData.Issuer)
}

func TestExtractWrongToken(t *testing.T) {
	_, err := service.Extract(rsa256)
	assert.NotEqual(t, nil, err)
}
