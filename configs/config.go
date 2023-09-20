package configs

import (
    "github.com/joho/godotenv"
    "os"
)

func LoadConfig(filePath string) error {
    if filePath != "" {
        err := godotenv.Load(filePath)
        if err != nil {
            return err
        }
    } else {
        err := godotenv.Load()
        if err != nil {
            return err
        }
    }

    return nil
}
func GetEnv(key string) string {
    return os.Getenv(key)
}
