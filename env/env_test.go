package env_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lcnascimento/go-kit/env"
)

func TestGetString(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "fake-env-1",
		"FAKE_ENV_2": "fake-env-2",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc          string
		env           string
		defaultValues []string
		expectedValue string
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			defaultValues: []string{"default-value-1", "default-value-2"},
			expectedValue: "fake-env-1",
		},
		{
			desc:          "should access an unknown environment variable with default value successfully",
			env:           "FAKE_ENV_3",
			defaultValues: []string{"default-value-1", "default-value-2"},
			expectedValue: "default-value-1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expectedValue, env.GetString(tc.env, tc.defaultValues...))
		})
	}
}

func TestMustGetString(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "fake-env-1",
		"FAKE_ENV_2": "fake-env-2",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc             string
		env              string
		expectedPanicMsg string
		expectedValue    string
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			expectedValue: "fake-env-1",
		},
		{
			desc:             "should panic when accessing an unknown environment variable",
			env:              "FAKE_ENV_3",
			expectedPanicMsg: "FAKE_ENV_3 can't be empty",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil && tc.expectedPanicMsg == "" {
					return
				}

				assert.Equal(t, tc.expectedPanicMsg, e)
			}()

			assert.Equal(t, tc.expectedValue, env.MustGetString(tc.env))
		})
	}
}

func TestGetInt(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "1",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc          string
		env           string
		defaultValues []int
		expectedValue int
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			defaultValues: []int{1, 2},
			expectedValue: 1,
		},
		{
			desc:          "should access an unknown environment variable with default value successfully",
			env:           "FAKE_ENV_3",
			defaultValues: []int{1, 2},
			expectedValue: 1,
		},
		{
			desc:          "should use default value when environment variable is not an Int",
			env:           "FAKE_ENV_2",
			defaultValues: []int{1, 2},
			expectedValue: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expectedValue, env.GetInt(tc.env, tc.defaultValues...))
		})
	}
}

func TestMustGetInt(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "1",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc             string
		env              string
		expectedPanicMsg string
		expectedValue    int
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			expectedValue: 1,
		},
		{
			desc:             "should panic when environment variable value is not and int",
			env:              "FAKE_ENV_2",
			expectedPanicMsg: "FAKE_ENV_2 must contain an int value",
		},
		{
			desc:             "should panic when accessing an unknown environment variable",
			env:              "FAKE_ENV_3",
			expectedPanicMsg: "FAKE_ENV_3 must contain an int value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil && tc.expectedPanicMsg == "" {
					return
				}

				assert.Equal(t, tc.expectedPanicMsg, e)
			}()

			assert.Equal(t, tc.expectedValue, env.MustGetInt(tc.env))
		})
	}
}

func TestGetFloat(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "1.7",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc          string
		env           string
		defaultValues []float64
		expectedValue float64
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			defaultValues: []float64{1.5, 2.3},
			expectedValue: 1.7,
		},
		{
			desc:          "should access an unknown environment variable with default value successfully",
			env:           "FAKE_ENV_3",
			defaultValues: []float64{1.5, 2.3},
			expectedValue: 1.5,
		},
		{
			desc:          "should use default value when environment variable is not a Float",
			env:           "FAKE_ENV_2",
			defaultValues: []float64{1.5, 2.3},
			expectedValue: 1.5,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expectedValue, env.GetFloat(tc.env, tc.defaultValues...))
		})
	}
}

func TestMustGetFloat(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "1.7",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc             string
		env              string
		expectedPanicMsg string
		expectedValue    float64
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			expectedValue: 1.7,
		},
		{
			desc:             "should panic when accessing an environment variable that isn't a Float",
			env:              "FAKE_ENV_2",
			expectedPanicMsg: "FAKE_ENV_2 must contain a float value",
		},
		{
			desc:             "should panic when accessing an unknown environment variable",
			env:              "FAKE_ENV_3",
			expectedPanicMsg: "FAKE_ENV_3 must contain a float value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil && tc.expectedPanicMsg == "" {
					return
				}

				assert.Equal(t, tc.expectedPanicMsg, e)
			}()

			assert.Equal(t, tc.expectedValue, env.MustGetFloat(tc.env))
		})
	}
}

func TestGetBool(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "true",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc          string
		env           string
		defaultValues []bool
		expectedValue bool
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			defaultValues: []bool{true, false},
			expectedValue: true,
		},
		{
			desc:          "should use default value when environment variable is not a Float",
			env:           "FAKE_ENV_2",
			defaultValues: []bool{true, false},
			expectedValue: true,
		},
		{
			desc:          "should access an unknown environment variable with default value successfully",
			env:           "FAKE_ENV_3",
			defaultValues: []bool{true, false},
			expectedValue: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expectedValue, env.GetBool(tc.env, tc.defaultValues...))
		})
	}
}

func TestMustGetBool(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "true",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc             string
		env              string
		expectedPanicMsg string
		expectedValue    bool
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			expectedValue: true,
		},
		{
			desc:             "should panic when accessing an environment variable that isn't a Boolean",
			env:              "FAKE_ENV_2",
			expectedPanicMsg: "FAKE_ENV_2 must contain a boolean value",
		},
		{
			desc:             "should panic when accessing an unknown environment variable",
			env:              "FAKE_ENV_3",
			expectedPanicMsg: "FAKE_ENV_3 must contain a boolean value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil && tc.expectedPanicMsg == "" {
					return
				}

				assert.Equal(t, tc.expectedPanicMsg, e)
			}()

			assert.Equal(t, tc.expectedValue, env.MustGetBool(tc.env))
		})
	}
}
