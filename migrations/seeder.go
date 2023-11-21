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

		symKey, err := utils.GenerateKey(32)
		if err != nil {
			return err
		}

		privKey, pubKey, err := utils.GenerateRSAKey()
		if err != nil {
			return err
		}

		datas := []string{data.Name, data.TelpNumber, pubKey, privKey, symKey}
		var encryptedDatas []string
		for i := 0; i < len(datas); i++ {
			encryptedData, _, err := utils.AESEncrypt(datas[i], utils.KEY)
			if err != nil {
				return err
			}
			encryptedDatas = append(encryptedDatas, encryptedData)
		}

		data.Name = encryptedDatas[0]
		data.TelpNumber = encryptedDatas[1]
		data.PublicKey = encryptedDatas[2]
		data.PrivateKey = encryptedDatas[3]
		data.SymmetricKey = encryptedDatas[4]

		isData := db.Find(&user, "email = ?", data.Email).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
