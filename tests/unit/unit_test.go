package unit

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/l1qwie/JWTAuth/app"
	"github.com/l1qwie/JWTAuth/app/database"
	"github.com/l1qwie/JWTAuth/app/logs"
	"github.com/l1qwie/JWTAuth/app/types"
	"github.com/l1qwie/JWTAuth/tests"
	"golang.org/x/crypto/bcrypt"
)

const testip = "192.168.1.101"
const testguid = "123e4567-e89b-12d3-a456-426614174000"
const tguid = "123e4567-e89b-12d3-a456-426614174001"
const refreshtoken = "a-refresh-token-for-test-with-the-real-ip::192.168.1.101"
const justrefreshtoken = "a-refresh-token-for-test-with-the-real-ip"

type test struct {
	name           string
	errMsg         []string
	input1, input2 []string
	isWrong        []bool
	checkedf       func(string, string) ([]byte, error)
	checkdbf       func(t *testing.T)
}

func testAction(test *test, t *testing.T) {
	for i := 0; i < len(test.errMsg); i++ {
		t.Logf("Function under test: %s\n\tInput(1): %v\n\tInput(2): %v\n\tExpect error: %v\n\tError message: %s",
			test.name, test.input1[i], test.input2[i], test.isWrong[i], test.errMsg[i])
		body, err := test.checkedf(test.input1[i], test.input2[i])
		if test.isWrong[i] {
			if err.Error() != test.errMsg[i] {
				t.Fatal(err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			tokens := new(types.Tokens)
			if err := json.Unmarshal(body, tokens); err != nil {
				t.Fatal(err)
			}
			t.Logf("Access-Token: %s", tokens.Access)
			t.Logf("Refresh-Token: %s", tokens.Refresh)
			test.checkdbf(t)
		}
	}
}

func isRenewed(t *testing.T) {
	bryptHash, err := database.Conn.GetRefreshToken("192.168.1.102")
	if err != nil {
		t.Fatal(err)
	}
	if err = bcrypt.CompareHashAndPassword(bryptHash, []byte(justrefreshtoken)); err.Error() != "crypto/bcrypt: hashedPassword is not the hash of the given password" {
		t.Fatal(err)
	}
}

func isRefreshSaved(t *testing.T) {
	if ok, err := database.Conn.IsThereRefreshToken(tguid, testip); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Fatal("there isn't the refresh token in database")
	}
}

func createEnv(t *testing.T) {
	var err error
	if database.Conn, err = database.Connect(); err != nil {
		t.Fatal(err)
	}
	if err = database.Conn.CreateMokData(tguid, testip); err != nil {
		t.Fatal(err)
	}
}

func createRefreshToken(token1, token2 string, t *testing.T) string {
	var s bool
	var bcryptHash []byte
	var err error
	encodedString := base64.StdEncoding.EncodeToString([]byte(token1))
	if bcryptHash, err = bcrypt.GenerateFromPassword([]byte(token2), bcrypt.DefaultCost); err != nil {
		t.Fatal(err)
	}
	if err = database.Conn.SaveRefreshToken(bcryptHash, tguid, &s); err != nil {
		t.Fatal(err)
	}
	return encodedString
}

func TestNewAccessAndRefreshTokens(t *testing.T) {
	createEnv(t)
	test := new(test)
	test.name = "app.NewAccessAndRefreshTokens"
	test.checkedf = app.NewAccessAndRefreshTokens
	test.errMsg = []string{"[ERROR:12] invalid guid", "[ERROR:11] invalid ip", "[ERROR:10] the guid is unknown", ""}
	test.input1 = []string{"12346", testguid, testguid, tguid}
	test.input2 = []string{testip, "12i03918094sx", testip, testip}
	test.isWrong = []bool{true, true, true, false}
	test.checkdbf = isRefreshSaved
	testAction(test, t)
	database.Conn.DeleteUsers()
}

func TestRefreshAction(t *testing.T) {
	createEnv(t)
	test := new(test)
	test.name = "app.RefreshAction"
	test.checkedf = app.RefreshAction
	test.errMsg = []string{"[ERROR:11] invalid ip", "[ERROR:14] a refresh token is required", "[ERROR:13] an ip in a refresh token is required", ""}
	test.input1 = []string{"023013=-542", testip, testip, "192.168.1.102"}
	test.input2 = []string{refreshtoken, "", createRefreshToken(justrefreshtoken, "", t), createRefreshToken(refreshtoken, justrefreshtoken, t)}
	test.isWrong = []bool{true, true, true, false}
	test.checkdbf = isRenewed
	testAction(test, t)
	database.Conn.DeleteUsers()
}

func init() {
	logs.SetDebug()
	tests.PutEnvVal()
}
