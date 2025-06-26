package services

import (
	"net/http"
	"runners/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRunnerInvalidFirstName(t *testing.T) {
	runner := &models.Runner{
		LastName: "Smith",
		Age:      30,
		Country:  "United States",
	}
	responseErr := validateRunner(runner)
	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid first name", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerInvalidLastName(t *testing.T) {
	runner := &models.Runner{
		FirstName: "Smith",
		Age:       30,
		Country:   "United States",
	}
	responseErr := validateRunner(runner)
	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid last name", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerInvalidAge(t *testing.T) {
	runner := &models.Runner{
		FirstName: "Smith",
		LastName:  "Smithon",
		Age:       -1,
		Country:   "United States",
	}
	responseErr := validateRunner(runner)
	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid age", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerInvalidCountry(t *testing.T) {
	runner := &models.Runner{
		FirstName: "Smith",
		LastName:  "Smithon",
		Age:       18,
		Country:   "",
	}
	responseErr := validateRunner(runner)
	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid country", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

// func TestValidateRunnerInvalidRunnerId(t *testing.T) {
// 	runner := &models.Runner{
// 		FirstName: "Smith",
// 		LastName:  "Smithon",
// 		Age:       -1,
// 		Country:   "",
// 	}
// 	responseErr := validateRunner(runner)
// 	assert.NotEmpty(t, responseErr)
// 	assert.Equal(t, "Invalid country", responseErr.Message)
// 	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
// }

func TestValidateRunner(t *testing.T) {
	type testCase struct {
		name   string
		runner models.Runner
		want   models.ResponseError
	}
	testCases := []testCase{
		{name: "Invalid first name",
			runner: models.Runner{
				LastName: "Smith",
				Age:      30,
				Country:  "United States"},
			want: models.ResponseError{
				Message: "Invalid first name",
				Status:  http.StatusBadRequest,
			},
		},
		{name: "Invalid last name",
			runner: models.Runner{
				FirstName: "Smith",
				Age:       30,
				Country:   "United States"},
			want: models.ResponseError{
				Message: "Invalid last name",
				Status:  http.StatusBadRequest,
			},
		},
	}

	// for _, testCase := range testCases {
	// 	t.Run(testCase.name, func(t *testing.T) {
	// 		responseError := validateRunner(&testCase.runner)
	// 		if responseError == nil {
	// 			t.Error("resonseError is empty")
	// 		} else {
	// 			if responseError.Message != testCase.want.Message {
	// 				t.Error("expected: ", testCase.want.Message, "\n", "getted: ", responseError.Message)
	// 			}
	// 			if responseError.Status != testCase.want.Status {
	// 				t.Error("expected: ", testCase.want.Status, "\n", "getted: ", responseError.Status)
	// 			}
	// 		}
	// 	})
	// }
	// using github.com/stretchr/testify/assert
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			responseError := validateRunner(&tc.runner)
			assert.NotEmpty(t, responseError)
			assert.Equal(t, tc.want.Message, responseError.Message)
			assert.Equal(t, tc.want.Status, responseError.Status)
		})
	}

}
