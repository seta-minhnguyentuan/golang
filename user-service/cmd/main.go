package main

import (
	"user-service/internal/team"
	"user-service/internal/user"
	"user-service/router"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=user_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&user.User{}, &team.Team{}, &team.TeamMember{})
	if err != nil {
		panic(err)
	}

	_, err = sqlDB.Exec(`
		DO $$ 
		BEGIN
			-- Add team_id foreign key constraint if it doesn't exist
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_team_members_team_id'
			) THEN
				ALTER TABLE team_members 
				ADD CONSTRAINT fk_team_members_team_id 
				FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE;
			END IF;

			-- Add user_id foreign key constraint if it doesn't exist
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_team_members_user_id'
			) THEN
				ALTER TABLE team_members 
				ADD CONSTRAINT fk_team_members_user_id 
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`)
	if err != nil {
		println("Warning: Could not create foreign key constraints:", err.Error())
	}

	userRepo := &user.GormRepository{DB: db}
	teamRepo := &team.GormRepository{DB: db}

	userService := &user.Service{Repo: userRepo}
	teamService := &team.Service{TeamRepo: teamRepo, UserRepo: userRepo}

	appRouter := router.NewRouter(userService, teamService)

	r := gin.Default()

	appRouter.SetupRoutes(r)

	r.Run(":8080")
}
