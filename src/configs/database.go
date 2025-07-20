package configs

import (
	"MB-test/src/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(env Env) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		env.POSTGRES_HOST,
		env.POSTGRES_USER,
		env.POSTGRES_PASSWORD,
		env.POSTGRES_DB,
		env.POSTGRES_PORT,
		env.SSLMODE,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	return db
}

func MigrateDb(db *gorm.DB) {
	err := db.AutoMigrate(&models.Client{}, &models.Orders{})
	if err != nil {
		panic("Erro na migração")
	}
}

func Seeders(db *gorm.DB) {
	log.Println("running seeds...")

	clients := []models.Client{
		{
			Id:         uuid.New(),
			BalanceBRL: 12.533,
			BalanceBT:  4,
			Score:      70,
			CreatedAt:  time.Now(),
		},
		{
			Id:         uuid.New(),
			BalanceBRL: 994.533,
			BalanceBT:  12,
			Score:      99,
			CreatedAt:  time.Now(),
		},
		{
			Id:         uuid.New(),
			BalanceBRL: 18.485,
			BalanceBT:  7,
			Score:      95,
			CreatedAt:  time.Now(),
		},
		{
			Id:         uuid.New(),
			BalanceBRL: 985,
			BalanceBT:  2,
			Score:      62,
			CreatedAt:  time.Now(),
		},
		{
			Id:         uuid.New(),
			BalanceBRL: 62875,
			BalanceBT:  35,
			Score:      100,
			CreatedAt:  time.Now(),
		},
	}

	for _, client := range clients {
		if err := db.FirstOrCreate(&client, models.Client{Id: client.Id}).Error; err != nil {
			log.Printf("Error seeding client %v: %v\n", client.Id, err)
		}
	}

	log.Println("Seeds completed")
}
