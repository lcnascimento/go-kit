package env_test

import (
	"errors"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/lcnascimento/go-kit/env"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetterTestSuite struct {
	suite.Suite
	require *require.Assertions
}

func TestGetterTestSuite(t *testing.T) {
	suite.Run(t, new(GetterTestSuite))
}

func (s *GetterTestSuite) SetupTest() {
	s.require = require.New(s.T())
}

func (s *GetterTestSuite) TestGetter() {

	s.Run("Simple Getter", func() {
		s.require.NoError(os.Setenv("NAME", "Han-Solo"))
		s.require.NoError(os.Setenv("AGE", "35"))
		s.require.NoError(os.Setenv("IS_SMUGGLER", "TRUE"))
		s.require.NoError(os.Setenv("SPOUSE_NAME", "Leia"))

		name := env.Get[string]("NAME")
		age := env.Get[int]("AGE")
		isSmuggler := env.Get[bool]("IS_SMUGGLER")
		masterName := env.Get[string]("MASTER_NAME")
		isJedi := env.Get[bool]("IS_JEDI")
		spouseName := env.Get[*string]("SPOUSE_NAME")
		spouseAge := env.Get[*int]("SPOUSE_AGE")
		friends := env.Get[[]string]("FRIENDS")

		s.require.Equal("Han-Solo", name)
		s.require.Equal(35, age)
		s.require.Equal(true, isSmuggler)
		s.require.Equal("", masterName) // string zero value
		s.require.Equal(false, isJedi)  // bool zero value
		s.require.Nil(spouseName)       // pointer not supported
		s.require.Nil(spouseAge)        // pointer not supported
		s.require.Nil(friends)          // slice not supported
	})

	s.Run("Getter With Default", func() {
		s.require.NoError(os.Setenv("PLANET_NAME", "Coruscant"))
		s.require.NoError(os.Setenv("NUMBER_OF_MOONS", "4"))

		name := env.Get[string]("PLANET_NAME")
		moons := env.Get("NUMBER_OF_MOONS", env.WithDefaultValue(1))
		suns := env.Get("NUMBER_OF_SUNS", env.WithDefaultValue(1))
		hasJediTemple := env.Get[bool]("HAS_JEDI_TEMPLE")

		s.require.Equal("Coruscant", name)    // defined value (no default)
		s.require.Equal(4, moons)             // defined value (with default)
		s.require.Equal(1, suns)              // default value
		s.require.Equal(false, hasJediTemple) // zero value
	})

	s.Run("Getter With Enum", func() {
		s.require.NoError(os.Setenv("SOCCER_TEAM", "Cruzeiro"))
		s.require.NoError(os.Setenv("SEASON", "2019/2020"))

		team := env.Get("SOCCER_TEAM", env.WithEnum([]string{"Cruzeiro", "Palmeiras"}))
		season := env.Get("SEASON", env.WithEnum([]string{"2020/2021", "2021/2022"}), env.WithDefaultValue("2020/2021"))

		s.require.Equal("Cruzeiro", team)    // defined value (no default)
		s.require.Equal("2020/2021", season) // defined value (with default)
	})

	s.Run("Getter With Custom Parser", func() {
		defaultProfile, err := url.Parse("https://starwars.fandom.com/pt/wiki/P%C3%A1gina_principal")
		s.require.NoError(err)

		planetParser := func(value string) (Planet, error) {
			return Planet{Name: value}, nil
		}
		pointerPlanetParser := func(value string) (*Planet, error) {
			return &Planet{Name: value}, nil
		}
		planetErrorParser := func(value string) (Planet, error) {
			return Planet{}, errors.New("planet not found")
		}

		s.require.NoError(os.Setenv("NAME", "Jabba"))
		s.require.NoError(os.Setenv("HOME_PLANET", "Nal Hutta"))

		name := env.Get[string]("NAME")
		homePlanet := env.Get("HOME_PLANET", env.WithCustomParser(planetParser))
		deathPlanet := env.Get("DEATH_PLANET", env.WithCustomParser(planetParser))
		workPlanet := env.Get("WORK_PLANET", env.WithCustomParser(pointerPlanetParser))
		datePlanet := env.Get("DATE_PLANET", env.WithCustomParser(planetErrorParser))
		profile := env.Get("PROFILE_URI", env.WithCustomParser(url.Parse), env.WithDefaultValue(defaultProfile))

		s.require.Equal("Jabba", name)                             // defined value
		s.require.Equal(Planet{Name: "Nal Hutta"}, homePlanet)     // defined value (with custom parser)
		s.require.Equal(defaultProfile.String(), profile.String()) // default value (with custom parser)
		s.require.Equal(Planet{Name: ""}, deathPlanet)             // zero value (with custom parser)
		s.require.Nil(workPlanet)                                  // pointer not supported
		s.require.Equal(Planet{}, datePlanet)                      // zero value (with custom parser with error)
	})

	s.Run("Validate Built In Parsers", func() {
		s.require.NoError(os.Setenv("VAR_BOOL", "1"))
		s.require.NoError(os.Setenv("VAR_STRING", "abc"))
		s.require.NoError(os.Setenv("VAR_INT", "66"))
		s.require.NoError(os.Setenv("VAR_INT16", "-55"))
		s.require.NoError(os.Setenv("VAR_INT32", "44"))
		s.require.NoError(os.Setenv("VAR_INT64", "33"))
		s.require.NoError(os.Setenv("VAR_INT8", "22"))
		s.require.NoError(os.Setenv("VAR_UINT", "-11"))
		s.require.NoError(os.Setenv("VAR_UINT16", "80"))
		s.require.NoError(os.Setenv("VAR_UINT32", "81"))
		s.require.NoError(os.Setenv("VAR_UINT64", "82"))
		s.require.NoError(os.Setenv("VAR_UINT8", "71"))
		s.require.NoError(os.Setenv("VAR_FLOAT64", "5.734"))
		s.require.NoError(os.Setenv("VAR_FLOAT32", "8.1243"))
		s.require.NoError(os.Setenv("VAR_DURATION", "1s"))

		varBool := env.Get[bool]("VAR_BOOL")
		varString := env.Get[string]("VAR_STRING")
		varInt := env.Get[int]("VAR_INT")
		varInt16 := env.Get[int16]("VAR_INT16")
		varInt32 := env.Get[int32]("VAR_INT32")
		varInt64 := env.Get[int64]("VAR_INT64")
		varInt8 := env.Get[int8]("VAR_INT8")
		varUint := env.Get[uint]("VAR_UINT")
		varUint16 := env.Get[uint16]("VAR_UINT16")
		varUint32 := env.Get[uint32]("VAR_UINT32")
		varUint64 := env.Get[uint64]("VAR_UINT64")
		varUint8 := env.Get[uint8]("VAR_UINT8")
		varFloat64 := env.Get[float64]("VAR_FLOAT64")
		varFloat32 := env.Get[float32]("VAR_FLOAT32")
		varDuration := env.Get[time.Duration]("VAR_DURATION")

		s.require.Equal(true, varBool)
		s.require.Equal("abc", varString)
		s.require.Equal(66, varInt)
		s.require.Equal(int16(-55), varInt16)
		s.require.Equal(int32(44), varInt32)
		s.require.Equal(int64(33), varInt64)
		s.require.Equal(int8(22), varInt8)
		s.require.Equal(uint(0), varUint)
		s.require.Equal(uint16(80), varUint16)
		s.require.Equal(uint32(81), varUint32)
		s.require.Equal(uint64(82), varUint64)
		s.require.Equal(uint8(71), varUint8)
		s.require.Equal(5.734, varFloat64)
		s.require.Equal(float32(8.1243), varFloat32)
		s.require.Equal(time.Second, varDuration)
	})
}

func (s *GetterTestSuite) TestListGetter() {

	s.Run("Simple Getter", func() {
		s.require.NoError(os.Setenv("LIST_NAME", "Han-Solo1,Han-Solo2,Han-Solo3"))
		s.require.NoError(os.Setenv("LIST_AGE", "35,36,37"))
		s.require.NoError(os.Setenv("LIST_IS_SMUGGLER", "TRUE,FALSE,TRUE"))
		s.require.NoError(os.Setenv("LIST_SPOUSE_NAME", "Leia1,Leia2,Leia3"))

		names := env.GetList[string]("LIST_NAME")
		ages := env.GetList[int]("LIST_AGE")
		isSmugglers := env.GetList[bool]("LIST_IS_SMUGGLER")
		spouseNames := env.GetList[*string]("LIST_SPOUSE_NAME")

		s.require.Equal([]string{"Han-Solo1", "Han-Solo2", "Han-Solo3"}, names)
		s.require.Equal([]int{35, 36, 37}, ages)
		s.require.Equal([]bool{true, false, true}, isSmugglers)
		s.require.Equal([]*string{}, spouseNames)
	})

	s.Run("Getter With Enum", func() {
		s.require.NoError(os.Setenv("LIST_SOCCER_TEAM", "Cruzeiro,Palmeiras,Internacional"))
		s.require.NoError(os.Setenv("LIST_SEASON", "2019/2020,2020/2021,2021/2022"))

		teams := env.GetList("LIST_SOCCER_TEAM", env.WithEnum([]string{"Cruzeiro", "Palmeiras", "Internacional"}))
		seasons := env.GetList("LIST_SEASON", env.WithEnum([]string{"2020/2021", "2021/2022"}))
		tournaments := env.GetList("LIST_TOURNAMENT",
			env.WithEnum([]string{"Copa do Brasil", "Brasileirão"}),
			env.WithDefaultListValue([]string{"Copa do Brasil", "Brasileirão"}),
		)

		s.require.Equal([]string{"Cruzeiro", "Palmeiras", "Internacional"}, teams) // defined value (no default)
		s.require.Equal([]string{"2020/2021", "2021/2022"}, seasons)               // defined value (with default)
		s.require.Equal([]string{"Copa do Brasil", "Brasileirão"}, tournaments)    // defined value (with default)
	})
}

type Planet struct {
	Name string
}
