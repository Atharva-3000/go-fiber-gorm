package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atharva-3000/go-fiber-gorm-postgres/models"
	"github.com/atharva-3000/go-fiber-gorm-postgres/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

// Create Book
func (r *Repository) CreateBook(context fiber.Ctx) error {
	book := Book{}

	err := context.Bind().Body(&book)
	if err != nil {
		context.Status(fiber.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Faialed !"})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Failed to create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book created successfully"})
	return nil
}

//  Get Books

func (r *Repository) GetBooks(context fiber.Ctx) error {
	bookModels := &[]models.Books{}
	err := r.DB.Find(&bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not fetch the books"},
		)
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Books fetched Successfully", "data": bookModels},
	)
	return nil
}

// delete book function

func (r *Repository) DeleteBook(context fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Id cannot be Empty!"})
	}
	err := r.DB.Delete(&bookModel, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Failed ot delete the book with the passed id"})
		return err
	}
	context.Status(http.StatusOK).JSON(fiber.Map{"message": "Book deleted successfully"})

	return nil
}

// getbook by id
func (r *Repository) GetBookByID(context fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Id cannot be Empty!"})
		return nil
	}
	fmt.Println(id)

	err := r.DB.Where("id = ?", id).First(&bookModel).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Failed to fetch the book with the passed"})
		return err
	}

	context.Status(http.StatusOK).JSON(fiber.Map{"message": "Book fetched successfully", "data": bookModel})
	return nil

}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/get_books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	errHandler(err)
	fmt.Println("Go Fiber and Postgres in Go with gorm")

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	errHandler(err)
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate db", err)
	}

	errHandler(err)
	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
