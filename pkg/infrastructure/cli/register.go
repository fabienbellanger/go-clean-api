package cli

import (
	"fmt"
	"go-clean-api/pkg/adapters/db"
	"go-clean-api/pkg/adapters/repositories/gorm_mysql"
	"go-clean-api/pkg/domain/usecases"
	vo "go-clean-api/pkg/domain/value_objects"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var (
	userEmail     string
	userPassword  string
	userLastname  string
	userFirstname string
)

func init() {
	userCmd.Flags().StringVarP(&userLastname, "lastname", "l", "", "user lastname")
	userCmd.Flags().StringVarP(&userFirstname, "firstname", "f", "", "user firstname")
	userCmd.Flags().StringVarP(&userEmail, "email", "e", "", "user email")
	userCmd.Flags().StringVarP(&userPassword, "password", "p", "", "user password")

	userCmd.MarkFlagRequired("lastname")
	userCmd.MarkFlagRequired("firstname")
	userCmd.MarkFlagRequired("email")
	userCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:   "register",
	Short: "User creation",
	Long:  `User creation`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := initConfig()
		if err != nil {
			log.Fatalln(err)
		}

		database, err := initDatabase(config)
		if err != nil {
			log.Fatalln(err)
		}

		email, err := vo.NewEmail(strings.TrimSpace(userEmail))
		if err != nil {
			log.Fatalln(err)
		}
		password, err := vo.NewPassword(strings.TrimSpace(userPassword))
		if err != nil {
			log.Fatalln(err)
		}

		// Call use case
		// sqlxDB, ok := database.(*db.SqlxMySQL)
		// if !ok {
		// 	log.Fatalln("db is not of type *db.SqlxMySQL")
		// }
		// userRepo := gorm_mysql.NewUser(sqlxDB)
		gormDB, ok := database.(*db.GormMySQL)
		if !ok {
			log.Fatalln("db is not of type *db.GormMySQL")
		}
		userRepo := gorm_mysql.NewUser(gormDB)
		userUseCase := usecases.NewUser(userRepo, config.JWT)
		res, errRes := userUseCase.Create(usecases.CreateUserRequest{
			Email:     email,
			Password:  password,
			Lastname:  strings.TrimSpace(userLastname),
			Firstname: strings.TrimSpace(userFirstname),
		})
		if errRes != nil {
			fmt.Printf("\nError: %s\n", errRes)
			return
		}

		// Display result
		fmt.Printf(`
User successfully created:
    - ID:        %s
    - Lastname:  %s
    - Firstname: %s
    - Email:     %s
    - Password:  %s
`,
			res.ID.Value(),
			res.Lastname,
			res.Firstname,
			res.Email.Value(),
			res.Password.Value(),
		)
	},
}
