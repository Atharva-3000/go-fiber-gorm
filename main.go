package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	fibre "github.com/gofiber/fiber/v3"
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

	err := context.BodyParser(&book)
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

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBooksByID)
	api.Get("/get_books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	errHandler(err)
	fmt.Println("Go Fiber and Postgres in Go with gorm")

	db, err := storage.NewConnection(config)
	errHandler(err)
	r := Repository{
		DB: db,
	}

	app := fibre.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
