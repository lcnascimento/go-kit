package env_test

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/lcnascimento/go-kit/env"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
	require *require.Assertions
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (s *ParserTestSuite) SetupTest() {
	s.require = require.New(s.T())
}

func (s *ParserTestSuite) TestParser() {

	s.Run("Success Parse", func() {
		name := "Han-Solo"
		age := 35

		s.require.NoError(os.Setenv("NAME", name))
		s.require.NoError(os.Setenv("AGE", strconv.Itoa(age)))

		config := MySimpleTestConfig{}
		err := env.Parse(&config)
		s.require.NoError(err)
		s.require.Equal(name, config.Name)
		s.require.Equal(age, config.Age)

		s.require.NoError(os.Unsetenv("NAME"))
		s.require.NoError(os.Unsetenv("AGE"))
	})

	s.Run("Success ParseWithConfig CustomType", func() {
		name := "Han-Solo"
		age := 35

		s.require.NoError(os.Setenv("NAME", name))
		s.require.NoError(os.Setenv("AGE", strconv.Itoa(age)))
		s.require.NoError(os.Setenv("EMAIL", "han-solo@rebels.org"))

		parseCfg := env.ParseConfig{}
		parseCfg.AddCustomTypeParser(env.CustomTypeParser(NewMyEmail))

		config := MyCustomTestConfig{}
		err := env.ParseWithConfig(&config, &parseCfg)
		s.require.NoError(err)
		s.require.Equal(name, config.Name)
		s.require.Equal(age, config.Age)
		s.require.Equal(MyEmail{Name: "han-solo", Domain: "rebels.org"}, config.Email)

		s.require.NoError(os.Unsetenv("NAME"))
		s.require.NoError(os.Unsetenv("AGE"))
	})

	s.Run("Success Parse Options CustomType", func() {
		name := "Han-Solo"
		age := 35

		s.require.NoError(os.Setenv("NAME", name))
		s.require.NoError(os.Setenv("AGE", strconv.Itoa(age)))
		s.require.NoError(os.Setenv("EMAIL", "han-solo@rebels.org"))

		config := MyCustomTestConfig{}
		err := env.Parse(&config, env.WithCustomTypeParser(NewMyEmail))
		s.require.NoError(err)
		s.require.Equal(name, config.Name)
		s.require.Equal(age, config.Age)
		s.require.Equal(MyEmail{Name: "han-solo", Domain: "rebels.org"}, config.Email)

		s.require.NoError(os.Unsetenv("NAME"))
		s.require.NoError(os.Unsetenv("AGE"))
	})

	s.Run("Success Parse From File", func() {
		name := "Han-Solo"
		age := 35

		config := MyCustomTestConfig{}
		err := env.Parse(&config, env.WithCustomTypeParser(NewMyEmail), env.FromFiles("sample/.env"))
		s.require.NoError(err)
		s.require.Equal(name, config.Name)
		s.require.Equal(age, config.Age)
		s.require.Equal(MyEmail{Name: "han-solo", Domain: "rebels.org"}, config.Email)
	})

}

type MySimpleTestConfig struct {
	Name string `env:"NAME"`
	Age  int    `env:"AGE"`
}

type MyCustomTestConfig struct {
	Name  string  `env:"NAME"`
	Age   int     `env:"AGE"`
	Email MyEmail `env:"EMAIL"`
}

type MyEmail struct {
	Name   string
	Domain string
}

func NewMyEmail(email string) (MyEmail, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return MyEmail{}, errors.New("invalid email")
	}
	return MyEmail{
		Name:   parts[0],
		Domain: parts[1],
	}, nil
}

func (s *ParserTestSuite) TestDefaultParser() {

	s.Run("Success Parse URL", func() {
		name := "Han-Solo"
		social := "https://instarebel.org/han-solo?weapon=blaster"
		socialURL, err := url.Parse(social)
		s.require.NoError(err)

		s.require.NoError(os.Setenv("NAME", name))
		s.require.NoError(os.Setenv("SOCIAL_MEDIA", social))

		config := MySimpleWithCustomConfig{}
		err = env.Parse(&config)
		s.require.NoError(err)
		s.require.Equal(name, config.Name)
		s.require.Equal(*socialURL, config.SocialMedia)

		s.require.NoError(os.Unsetenv("NAME"))
		s.require.NoError(os.Unsetenv("SOCIAL_MEDIA"))
	})

	s.Run("Error Invalid URL", func() {
		name := "Han-Solo"
		social := "https://%inv^&#"
		_, err := url.Parse(social)
		s.require.Error(err)

		s.require.NoError(os.Setenv("NAME", name))
		s.require.NoError(os.Setenv("SOCIAL_MEDIA", social))

		config := MySimpleWithCustomConfig{}
		err = env.Parse(&config)
		s.require.Error(err)

		s.require.NoError(os.Unsetenv("NAME"))
		s.require.NoError(os.Unsetenv("SOCIAL_MEDIA"))
	})

	s.Run("With prefix", func() {
		familyName := "Solo"
		s.T().Setenv("FAMILY_NAME", familyName)

		config := MySimpleWithCustomConfig{}
		err := env.Parse(&config, env.WithPrefix("FAMILY_"))
		s.require.NoError(err)
		s.require.Equal(config.Name, "Solo")
	})
}

type MySimpleWithCustomConfig struct {
	Name        string  `env:"NAME"`
	SocialMedia url.URL `env:"SOCIAL_MEDIA"`
}
