package main

import (
	"fmt"
	"log"
	"time"

	"agrocontrol-api/configs"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/utils"
)

func main() {
	cfg := configs.LoadConfig()
	db := configs.ConnectDatabase(cfg)

	if err := db.AutoMigrate(
		&entities.User{}, &entities.Farm{}, &entities.Field{},
		&entities.Crop{}, &entities.Season{}, &entities.Planting{},
		&entities.Input{}, &entities.Application{}, &entities.Monitoring{},
		&entities.Harvest{}, &entities.Alert{},
	); err != nil {
		log.Fatal("migrate: ", err)
	}

	// ── Usuários
	adminHash, _ := utils.HashPassword("Admin@123")
	managerHash, _ := utils.HashPassword("Manager@123")
	operatorHash, _ := utils.HashPassword("Operator@123")

	users := []entities.User{
		{Name: "Administrador", Email: "admin@agrocontrol.com", PasswordHash: adminHash, Role: "admin", Active: true},
		{Name: "Gerente Silva", Email: "gerente@agrocontrol.com", PasswordHash: managerHash, Role: "manager", Active: true},
		{Name: "Operador João", Email: "operador@agrocontrol.com", PasswordHash: operatorHash, Role: "operator", Active: true},
	}
	for i := range users {
		db.Where("email = ?", users[i].Email).FirstOrCreate(&users[i])
	}
	fmt.Println("✔ Usuários criados")

	// ── Fazenda
	farm := entities.Farm{
		Name: "Fazenda Bom Futuro", OwnerName: "José da Silva",
		Location: "Zona Rural", TotalArea: 1500, City: "Campo Mourão", State: "PR", CreatedBy: 1,
	}
	db.Where("name = ?", farm.Name).FirstOrCreate(&farm)
	fmt.Println("✔ Fazenda criada")

	// ── Talhões
	fields := []entities.Field{
		{Name: "Talhão 01", Area: 120, SoilType: "Argiloso", Status: entities.FieldStatusActive, FarmID: farm.ID, CreatedBy: 1},
		{Name: "Talhão 02", Area: 95, SoilType: "Latossolo", Status: entities.FieldStatusActive, FarmID: farm.ID, CreatedBy: 1},
		{Name: "Talhão 03", Area: 80, SoilType: "Arenoso", Status: entities.FieldStatusFallow, FarmID: farm.ID, CreatedBy: 1},
	}
	for i := range fields {
		db.Where("name = ? AND farm_id = ?", fields[i].Name, fields[i].FarmID).FirstOrCreate(&fields[i])
	}
	fmt.Println("✔ Talhões criados")

	// ── Culturas
	crops := []entities.Crop{
		{Name: "Soja", Variety: "TMG 7062 IPRO", Type: "leguminosa", CycleDays: 115, SpacingCm: 45, PlantPopulation: 280000, CreatedBy: 1},
		{Name: "Milho", Variety: "DKB 290 PRO3", Type: "cereal", CycleDays: 130, SpacingCm: 70, PlantPopulation: 65000, CreatedBy: 1},
		{Name: "Trigo", Variety: "BRS Gaivota", Type: "cereal", CycleDays: 100, SpacingCm: 17, PlantPopulation: 350000, CreatedBy: 1},
	}
	for i := range crops {
		db.Where("name = ? AND variety = ?", crops[i].Name, crops[i].Variety).FirstOrCreate(&crops[i])
	}
	fmt.Println("✔ Culturas criadas")

	// ── Safra
	season := entities.Season{
		Name:      "Safra 2025/2026",
		StartDate: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
		Status:    entities.SeasonStatusActive,
		CreatedBy: 1,
	}
	db.Where("name = ?", season.Name).FirstOrCreate(&season)
	fmt.Println("✔ Safra criada")

	// ── Insumos
	inputs := []entities.Input{
		{Name: "Glifosato 480", Category: entities.InputCategoryHerbicide, Manufacturer: "Bayer", Unit: "L", StockQty: 500, MinStockQty: 50, CostPerUnit: 18.50, Active: true, CreatedBy: 1},
		{Name: "Ureia 45%", Category: entities.InputCategoryFertilizer, Manufacturer: "Yara", Unit: "sc", StockQty: 200, MinStockQty: 20, CostPerUnit: 145.00, Active: true, CreatedBy: 1},
		{Name: "Opera", Category: entities.InputCategoryFungicide, Manufacturer: "BASF", Unit: "L", StockQty: 80, MinStockQty: 10, CostPerUnit: 85.00, Active: true, CreatedBy: 1},
	}
	for i := range inputs {
		db.Where("name = ?", inputs[i].Name).FirstOrCreate(&inputs[i])
	}
	fmt.Println("✔ Insumos criados")

	fmt.Println()
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  Seeds concluídos com sucesso!")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  admin@agrocontrol.com    / Admin@123")
	fmt.Println("  gerente@agrocontrol.com  / Manager@123")
	fmt.Println("  operador@agrocontrol.com / Operator@123")
}
