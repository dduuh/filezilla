package psql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewPostgreSQL() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
			viper.GetString("postgres_db.user"),
			viper.GetString("postgres_db.password"),
			viper.GetString("postgres_db.host"),
			viper.GetString("postgres_db.port"),
			viper.GetString("postgres_db.dbname"),
			viper.GetString("postgres_db.sslmode")))
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logrus.Fatal(err)
		return nil, err
	}

	return db, nil
}
