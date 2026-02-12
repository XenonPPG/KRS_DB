package initializers

import (
	"DB/internal/models"

	"github.com/ilyakaznacheev/cleanenv"
)

func LoadConfig(path string) (config models.Config, err error) {
	err = cleanenv.ReadConfig(path, &config)
	return
}
