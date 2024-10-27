package main

import (
	"fmt"

	"github.com/Witthaya22/golang-backend-itctc/config"
	"github.com/Witthaya22/golang-backend-itctc/databases"
	"github.com/Witthaya22/golang-backend-itctc/entities"
)

func main() {
	conf := config.ConfigGeting()
	db := databases.NewPostgresDb(conf.Database)

	fmt.Println(db.ConnectionGetting())

	createBaseTables(db)
	addForeignKeys(db)
}

func createBaseTables(db databases.IDatabase) {
	tables := []interface{}{
		&entities.Department{},
		&entities.User{},
		&entities.Admin{},
		&entities.Activity{},
		&entities.ActivityResults{},
		&entities.ActivityDetails{},
		&entities.Oauth{},
	}

	for _, table := range tables {
		if db.ConnectionGetting().Migrator().HasTable(table) {
			if err := db.ConnectionGetting().Migrator().DropTable(table); err != nil {
				fmt.Printf("Error dropping table: %v\n", err)
				continue
			}
		}
		if err := db.ConnectionGetting().Migrator().CreateTable(table); err != nil {
			fmt.Printf("Error creating table: %v\n", err)
		}
	}
}

func addForeignKeys(db databases.IDatabase) {
	foreignKeys := []struct {
		table      string
		constraint string
		sql        string
	}{
		{"users", "fk_departments_users", "ALTER TABLE users ADD CONSTRAINT fk_departments_users FOREIGN KEY (department_id) REFERENCES departments(department_id)"},
		{"admins", "fk_users_admin", "ALTER TABLE admins ADD CONSTRAINT fk_users_admin FOREIGN KEY (user_id) REFERENCES users(user_id)"},
		{"activity_results", "fk_activities_results", "ALTER TABLE activity_results ADD CONSTRAINT fk_activities_results FOREIGN KEY (activity_id) REFERENCES activities(activity_id)"},
		{"activity_results", "fk_users_activities", "ALTER TABLE activity_results ADD CONSTRAINT fk_users_activities FOREIGN KEY (user_id) REFERENCES users(user_id)"},
		{"activity_results", "fk_departments_activities", "ALTER TABLE activity_results ADD CONSTRAINT fk_departments_activities FOREIGN KEY (department_id) REFERENCES departments(department_id)"},
		{"activity_details", "fk_users_activity_details", "ALTER TABLE activity_details ADD CONSTRAINT fk_users_activity_details FOREIGN KEY (user_id) REFERENCES users(user_id)"},
		{"activity_details", "fk_activities_activity_details", "ALTER TABLE activity_details ADD CONSTRAINT fk_activities_activity_details FOREIGN KEY (activity_id) REFERENCES activities(activity_id)"},
		// เพิ่มการเชื่อมจาก users ไป oauths
		{"oauths", "fk_users_oauth", "ALTER TABLE oauths ADD CONSTRAINT fk_users_oauth FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL ON UPDATE CASCADE"},
	}

	for _, fk := range foreignKeys {
		if db.ConnectionGetting().Migrator().HasTable(fk.table) {
			var count int64
			db.ConnectionGetting().Raw("SELECT COUNT(*) FROM pg_constraint WHERE conname = ?", fk.constraint).Scan(&count)

			if count == 0 {
				if err := db.ConnectionGetting().Exec(fk.sql).Error; err != nil {
					fmt.Printf("Error adding foreign key to %s: %v\n", fk.table, err)
				} else {
					fmt.Printf("Successfully added foreign key on %s\n", fk.table)
				}
			} else {
				fmt.Printf("Foreign key %s already exists on %s\n", fk.constraint, fk.table)
			}
		} else {
			fmt.Printf("Table %s does not exist, skipping foreign key\n", fk.table)
		}
	}
}
