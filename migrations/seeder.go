package migrations

import (
	"errors"

	"github.com/Caknoooo/golang-clean_template/constants"
	"github.com/Caknoooo/golang-clean_template/entities"
	"github.com/Caknoooo/golang-clean_template/utils"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := ListUserSeeder(db); err != nil {
		return err
	}

	return nil
}

func ListUserSeeder(db *gorm.DB) error {
	var listUser = []entities.User{
		{
			Name:       "Admin",
			TelpNumber: "081234567890",
			Email:      "admin@gmail.com",
			Password:   "admin123",
			Role:       constants.ENUM_ROLE_ADMIN,
			IsVerified: true,
		},
		{
			Name:       "User",
			TelpNumber: "081234567891",
			Email:      "user@gmail.com",
			Password:   "user123",
			Role:       constants.ENUM_ROLE_USER,
			IsVerified: true,
		},
	}

	hasTable := db.Migrator().HasTable(&entities.User{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entities.User{}); err != nil {
			return err
		}
	}

	for _, data := range listUser {
		var user entities.User
		var err error

		err = db.Where(&entities.User{Email: data.Email}).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if user != (entities.User{}) {
			break
		}

		data.TelpNumber, _, err = utils.AESEncrypt(data.TelpNumber, utils.KEY)
		if err != nil {
			return err
		}

		data.Name, _, err = utils.AESEncrypt(data.Name, utils.KEY)
		if err != nil {
			return err
		}

		symKey, err := utils.GenerateKey(32)
		if err != nil {
			return err
		}

		pubSymKey, err := utils.GenerateKey(32)
		if err != nil {
			return err
		}

		privKey, pubKey, err := utils.GenerateRSAKey()
		if err != nil {
			return err
		}

		symKey, _, err = utils.AESEncrypt(symKey, utils.KEY)
		if err != nil {
			return err
		}

		pubSymKey, _, err = utils.AESEncrypt(pubSymKey, utils.KEY)
		if err != nil {
			return err
		}

		pubKey, _, err = utils.AESEncrypt(pubKey, utils.KEY)
		if err != nil {
			return err
		}

		privKey, _, err = utils.AESEncrypt(privKey, utils.KEY)
		if err != nil {
			return err
		}

		keys := [4]string{pubKey, privKey, symKey, pubSymKey}
		var encryptedKeys []string
		for i := 0; i < 4; i++ {
			encryptedKey, _, err := utils.AESEncrypt(keys[i], utils.KEY)
			if err != nil {
				return err
			}
			encryptedKeys = append(encryptedKeys, encryptedKey)
		}

		data.PublicKey = encryptedKeys[0]
		data.PrivateKey = encryptedKeys[1]
		data.SymmetricKey = encryptedKeys[2]
		data.PublicSymmetricKey = encryptedKeys[3]

		isData := db.Find(&user, "email = ?", data.Email).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
