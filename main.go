package main

import (
	// "bytes"
	// "crypto"
	// "crypto/rsa"
	// "crypto/sha256"
	// "fmt"
	// "github.com/Caknoooo/golang-clean_template/utils"
	"os"

	"github.com/Caknoooo/golang-clean_template/config"
	"github.com/Caknoooo/golang-clean_template/controller"
	"github.com/Caknoooo/golang-clean_template/middleware"

	// "github.com/Caknoooo/golang-clean_template/utils"

	// "github.com/Caknoooo/golang-clean_template/migrations"
	"github.com/Caknoooo/golang-clean_template/repository"
	"github.com/Caknoooo/golang-clean_template/routes"
	"github.com/Caknoooo/golang-clean_template/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	var (
		db         *gorm.DB            = config.SetUpDatabaseConnection()
		jwtService services.JWTService = services.NewJWTService()

		// Repo
		userRepository             repository.UserRepository             = repository.NewUserRepository(db)
		fileRepository             repository.FileRepository             = repository.NewFileRepository(db)
		privateAccessRepository    repository.PrivateAccessRepository    = repository.NewPrivateAccessRepository(db)
		digitalSignatureRepository repository.DigitalSignatureRepository = repository.NewDigitalSignatureRepository(db)

		// Service
		userService             services.UserService             = services.NewUserService(userRepository, fileRepository)
		fileService             services.FileService             = services.NewFileService(fileRepository)
		privateAccessService    services.PrivateAccessService    = services.NewPrivateAccessService(userRepository, privateAccessRepository, fileRepository)
		digitalSignatureService services.DigitalSignatureService = services.NewDigitalSignatureService(digitalSignatureRepository, userRepository)

		// Controller
		userController             controller.UserController             = controller.NewUserController(userService, jwtService)
		fileController             controller.FileController             = controller.NewFileController(fileService, jwtService)
		privateAccessController    controller.PrivateAccessController    = controller.NewPrivateAccessController(privateAccessService)
		digitalSignatureController controller.DigitalSignatureController = controller.NewDigitalSignatureController(digitalSignatureService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())
	routes.User(server, userController, jwtService)
	routes.File(server, fileController, jwtService)
	routes.PrivateAccess(server, privateAccessController, jwtService)
	routes.DigitalSignature(server, digitalSignatureController, jwtService)

	// if err := migrations.Seeder(db); err != nil {
	// 	log.Fatalf("error migration seeder: %v", err)
	// }

	// tes()

	// tes2()
	server.Static("/storage", "./storage")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	server.Run(":" + port)
}

// func tes2() {
// 	filePath := "storage/ed77100d-3a9e-4e15-9d88-ad723e79e2ec/digital_signature/67c6500d-bae6-441a-928e-f3f7937c4d9e.pdf"

// 	file, err := os.ReadFile(filePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	const (
// 		DataCommentKey      = "DataKeyF-02_"
// 		SignatureCommentKey = "SignatureKeyF-02_"
// 		PublicKeyCommentKey = "PublicKeyKeyF-02_"
// 	)

// 	tokens := []string{
// 		"%" + DataCommentKey,
// 		"%" + SignatureCommentKey,
// 		"%" + PublicKeyCommentKey,
// 	}

// 	results := [3]string{}

// 	ctr := 0

// 	idx := 0

// 	stringFileContent := string(file)

// 	var buffer []byte
// 	for idx < len(stringFileContent) {
// 		if ctr <= 2 &&
// 			idx+len(tokens[ctr]) < len(stringFileContent) &&
// 			stringFileContent[idx:idx+len(tokens[ctr])] == tokens[ctr] {

// 			if ctr > 0 && len(buffer) > 0 {
// 				results[ctr-1] = string(buffer[:len(buffer)-1])
// 				buffer = []byte{}
// 			}

// 			idx += len(stringFileContent[idx : idx+len(tokens[ctr])])
// 			ctr++
// 		}

// 		if ctr > 0 {
// 			buffer = append(buffer, stringFileContent[idx])
// 		}

// 		if idx == len(stringFileContent)-1 && ctr != 0 {
// 			results[ctr-1] = string(buffer)
// 		}

// 		idx++
// 	}

// 	fmt.Println(results[1])



// 	// Decrypt content file pdf
// 	decryptPublicKey, err := utils.AESDecrypt("40fcc710cbd3ac6416e79bae3111672318554abf0031bb56cfacd1d567be14dfb314c0ae35a69438ff4c247f4e7351ffc23667831be09e79786df1e329f56ad94c64961581eb0749a5b5faef58c2c43d6e1982501d8ed4bc021ac0c44d3c520d82febff779d76cd2dd53313f1ad9dd8358a9bde76ce6fe0be3677c9b2405e56d4688984cb7546d2cf83ffbdc1ba755d6ce394d03d896a8fa128eca295c6f68d95e11bca82e869eb196b4dd6f960afb70d72abad884a6a7f6f36913d6b38b33a8d9450dc06db3c63eee0c3308371a038b4d83d060a60f44be125f73182ea7152a47d21dc1903d238bbc971920e2ffe98f64147bd1eda15e0e0a5632b1197d39d9221589e7c822740a4e1d1918e54f343ae48d5adebe45c1c2353328c16c3530a81ef897a8f985b824eca081cea8c2da020b489759d699f2ca4c55c29d4049d420560f2bf1061e4b37c0c3c076b995631cef9a79fb27eeeeab4dc7dadc37f74ee031cafcfb513e9d50676875b54e0c279428198e641def17d697e521f2a930bf492363995cb060162d5a458e3251c3c3e82d6e60c79efd67c0d671f12ebcf2266967dbd0b374526dd7340534769c6e4fac794951f06e924e01f5cabcbb7b941c66107f4de181d988a039a4b8ffbaa377fcdbaa1e7689f87c7c6055d3e6f685ae0876302e2c17630c", utils.KEY)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// fmt.Println(decryptPublicKey)

// 	parts := bytes.Split(file, []byte("\n"+"%"+DataCommentKey))
// 	if len(parts) < 2 {
// 		panic("error")
// 	}

// 	// fmt.Println(string(parts[1]))

// 	data := parts[0]
// 	signature := results[0]

// 	//

// 	hash := sha256.Sum256(data)

// 	pubKey, err := utils.ParsePublicKeyFromPEM(decryptPublicKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// fmt.Println("pubKey: ", pubKey)

// 	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], []byte(signature))
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Verified!")
// }

// func tes() {
// 	filePath := "storage/CV_M Naufal Badruttamam.pdf"
// 	newPath := "storage/tes.pdf"
// 	messagesAdded := []byte{}

// Open the file for reading and writing
// file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
// if err != nil {
// 	panic(err)
// }
// defer file.Close()

// // Write the bytes to the file
// if _, err := file.Write(messagesAdded); err != nil {
// 	panic(err)
// }

// Write "This is a test" to the file
// if _, err := file.Write([]byte("\nThis is a test")); err != nil {
// 	panic(err)
// }

// 	file, err := os.ReadFile(filePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	data, _, err := utils.AESEncrypt(string(file), utils.FILE_KEY_AES)
// 	if err != nil {
// 		panic(err)
// 	}

// 	const (
// 		DataCommentKey      = "DataKeyF-02_"
// 		SignatureCommentKey = "SignatureF-02_"
// 		PublicKeyCommentKey = "PublicKeyF-02_"
// 	)

// 	messagesAdded = append(messagesAdded, file...)

// 	byteData := []byte(data)
// 	messagesAdded = append(messagesAdded, []byte("\n%"+DataCommentKey)...)
// 	messagesAdded = append(messagesAdded, byteData...)

// 	messagesAdded = append(messagesAdded, []byte("\n%"+SignatureCommentKey)...)
// 	messagesAdded = append(messagesAdded, []byte("tes")...)

// 	messagesAdded = append(messagesAdded, []byte("\n%"+PublicKeyCommentKey)...)
// 	messagesAdded = append(messagesAdded, []byte("tes2")...)

// 	// write to new pdf
// 	err = os.WriteFile(newPath, messagesAdded, 0644)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fileContent, err := os.ReadFile(newPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Decrypt content file pdf
// 	tokens := []string{
// 		"%" + DataCommentKey,
// 		"%" + SignatureCommentKey,
// 		"%" + PublicKeyCommentKey,
// 	}

// 	results := [3]string{}

// 	ctr := 0

// 	idx := 0

// 	stringFileContent := string(fileContent)

// 	var buffer []byte
// 	for idx < len(stringFileContent) {
// 		if ctr <= 2 &&
// 			idx+len(tokens[ctr]) < len(stringFileContent) &&
// 			stringFileContent[idx:idx+len(tokens[ctr])] == tokens[ctr] {

// 			if ctr > 0 && len(buffer) > 0 {
// 				results[ctr-1] = string(buffer[:len(buffer)-1])
// 				buffer = []byte{}
// 			}

// 			idx += len(stringFileContent[idx : idx+len(tokens[ctr])])
// 			ctr++
// 		}

// 		if ctr > 0 {
// 			buffer = append(buffer, stringFileContent[idx])
// 		}

// 		if idx == len(stringFileContent)-1 && ctr != 0 {
// 			results[ctr-1] = string(buffer)
// 		}

// 		idx++
// 	}

// 	decrypt, err := utils.AESDecrypt(results[0], utils.FILE_KEY_AES)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// fmt.Println(decrypt)

// 	// Compare filepath and newPath
// 	old, err := os.ReadFile(filePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// news, err := os.ReadFile(decrypt)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	if compareByteArray(old, []byte(decrypt)) {
// 		fmt.Println("The PDF files have the same content.")
// 	} else {
// 		fmt.Println("The PDF files have different content.")
// 	}
// }

// func compareByteArray(arr1, arr2 []byte) bool {
//     if len(arr1) != len(arr2) {
//         return false
//     }

//     for i := range arr1 {
//         if arr1[i] != arr2[i] {
//             return false
//         }
//     }

//     return true
// }
