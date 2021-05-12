package database

import (
	"atlas-cos/character"
	"atlas-cos/equipment"
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/location"
	"atlas-cos/retry"
	"atlas-cos/skill"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDatabase(l logrus.FieldLogger) *gorm.DB {
	var db *gorm.DB

	tryToConnect := func(attempt int) (bool, error) {
		var err error
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       "root:the@tcp(atlas-db:3306)/atlas-cos?charset=utf8&parseTime=True&loc=Local",
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}), &gorm.Config{})
		if err != nil {
			return true, err
		}
		return false, err
	}

	err := retry.Try(tryToConnect, 10)
	if err != nil {
		l.WithError(err).Fatalf("Failed to connect to database.")
	}

	// Migrate the schema
	character.Migration(db)
	equipment.Migration(db)
	item.Migration(db)
	location.Migration(db)
	skill.Migration(db)
	inventory.Migration(db)
	return db
}
