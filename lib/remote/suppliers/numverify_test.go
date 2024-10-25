package suppliers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/url"
	"os"
	"testing"
)

func TestNumverifySupplierSuccessCustomApiKey(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "4793068820"
	apikey := "e96f38cfe8f8dd6920e015a150859598"

	expectedResult := &NumverifyValidateResponse{
		Valid:               true,
		Number:              "79516566591",
		LocalFormat:         "9516566591",
		InternationalFormat: "+79516566591",
		CountryPrefix:       "+7",
		CountryCode:         "RU",
		CountryName:         "Russian Federation",
		Location:            "Saint Petersburg and Leningrad Oblast",
		Carrier:             "OJSC St. Petersburg Telecom (OJSC Tele2-Saint-Petersburg)",
		LineType:            "mobile",
	}

	gock.New("https://api.apilayer.com").
		Get("/number_verification/validate").
		MatchHeader("Apikey", apikey).
		MatchParam("number", number).
		Reply(200).
		JSON(expectedResult)

	s := NewNumverifySupplier()

	got, err := s.Request().SetApiKey(apikey).ValidateNumber(number)
	assert.Nil(t, err)

	assert.Equal(t, expectedResult, got)
}

func TestNumverifySupplierError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "4793068820"
	apikey := "e96f38cfe8f8dd6920e015a150859598"

	expectedResult := &NumverifyErrorResponse{
		Message: "You have exceeded your daily\\/monthly API rate limit. Please review and upgrade your subscription plan at https:\\/\\/apilayer.com\\/subscriptions to continue.",
	}

	gock.New("https://api.apilayer.com").
		Get("/number_verification/validate").
		MatchHeader("Apikey", apikey).
		MatchParam("number", number).
		Reply(429).
		JSON(expectedResult)

	s := NewNumverifySupplier()

	got, err := s.Request().SetApiKey(apikey).ValidateNumber(number)
	assert.Nil(t, got)
	assert.Equal(t, errors.New("You have exceeded your daily\\/monthly API rate limit. Please review and upgrade your subscription plan at https:\\/\\/apilayer.com\\/subscriptions to continue."), err)
}

func TestNumverifySupplierHTTPError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "4793068820"

	_ = os.Setenv("NUMVERIFY_API_KEY", "e96f38cfe8f8dd6920e015a150859598")
	defer os.Clearenv()

	dummyError := errors.New("test")

	gock.New("https://api.apilayer.com").
		Get("/number_verification/validate").
		ReplyError(dummyError)

	s := NewNumverifySupplier()

	got, err := s.Request().ValidateNumber(number)
	assert.Nil(t, got)
	assert.Equal(t, &url.Error{
		Op:  "Get",
		URL: "https://api.apilayer.com/number_verification/validate?number=4793068820",
		Err: dummyError,
	}, err)
}
