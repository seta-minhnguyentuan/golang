package database

import (
	"user-service/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.User{}, &models.Team{}, &models.TeamMember{})
	if err != nil {
		return err
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
		// Log warning but don't fail the migration
		return nil
	}

	return nil
}
