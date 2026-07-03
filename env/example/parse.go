package sample

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lcnascimento/go-kit/env"
)

type PlanetParams struct {
	Name string `env:"PLANET_NAME,required"`
	Suns int    `env:"PLANET_SUNS"`
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

type DeathStarConfig struct {
	TargetPlanet  PlanetParams
	CommanderName string     `env:"COMMANDER_NAME"`
	ShootPosition Coordinate `env:"SHOOT_POSITION"`
}

func GetPlanetParamsWithoutRequired() {

	cfg := &PlanetParams{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	slog.Info("DeathStarConfig", slog.Any("config", cfg))

	// panic with: env: required environment variable "PLANET_NAME" is not set
}

func GetPlanetParams() {

	_ = os.Setenv("PLANET_NAME", "Coruscant")
	_ = os.Setenv("PLANET_SUNS", "1")

	cfg := &PlanetParams{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	slog.Info("DeathStarConfig", slog.Any("config", cfg))

	// output INFO DeathStarConfig config="&{Name:Coruscant Suns:1}"
}

func GetDeathStarConfigWithEnvs() {

	_ = os.Setenv("PLANET_NAME", "Tatooine")
	_ = os.Setenv("PLANET_SUNS", "2")
	_ = os.Setenv("COMMANDER_NAME", "Governor Tarkin")
	_ = os.Setenv("SHOOT_POSITION", "123.456,789.012,345.678")

	// custom type parser, just convert a string to the custom type
	stringToCoordinate := func(s string) (Coordinate, error) {
		var c Coordinate
		_, err := fmt.Sscanf(s, "%f,%f,%f", &c.Latitude, &c.Longitude, &c.Altitude)
		return c, err
	}

	cfg := &DeathStarConfig{}
	if err := env.Parse(cfg, env.WithCustomTypeParser(stringToCoordinate)); err != nil {
		panic(err)
	}

	slog.Info("DeathStarConfig", slog.Any("config", cfg))

	// output INFO DeathStarConfig config="&{TargetPlanet:{Name:Tatooine Suns:2} CommanderName:Governor Tarkin ShootPosition:{Latitude:123.456 Longitude:789.012 Altitude:345.678}}"
}

func GetDeathStarConfigFromFile() {
	cfg := &DeathStarConfig{}
	if err := env.Parse(cfg, env.FromFiles("sample/.env")); err != nil {
		panic(err)
	}

	slog.Info("DeathStarConfig", slog.Any("config", cfg))
}
