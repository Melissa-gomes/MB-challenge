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
			Id:         uuid.MustParse("b7050560-3387-4318-812d-f671ae9caa6e"),
			BalanceBRL: 12533,
			BalanceBT:  4,
			Score:      70,
			CreatedAt:  time.Now(),
		},
		{
			Id:         uuid.MustParse("2268237d-1079-47e8-b7b2-8ab9ae1942f5"),
			BalanceBRL: 994533,
			BalanceBT:  12,
			Score:      99,
			CreatedAt:  time.Now(),
		},
		{
			Id:         uuid.MustParse("d3909b31-045b-4c3e-a6f8-2edb54316b37"),
			BalanceBRL: 18485,
			BalanceBT:  7,
			Score:      95,
			CreatedAt:  time.Now(),
		},
		{
			Id:         uuid.MustParse("e65d206b-aa5c-4d47-8684-672b2bc8a826"),
			BalanceBRL: 985,
			BalanceBT:  2,
			Score:      62,
			CreatedAt:  time.Now(),
		},
		{
			Id:         uuid.MustParse("dc333741-4adc-4e28-89a1-f0e45d38b2db"),
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
