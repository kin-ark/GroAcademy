package seeds

import (
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
	"github.com/kin-ark/GroAcademy/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedAll() error {
	log.Println("Starting database seeding...")

	if err := s.clearData(); err != nil {
		return fmt.Errorf("failed to clear data: %w", err)
	}

	if err := s.seedUsers(50); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := s.seedCourses(20); err != nil {
		return fmt.Errorf("failed to seed courses: %w", err)
	}

	if err := s.seedModules(); err != nil {
		return fmt.Errorf("failed to seed modules: %w", err)
	}

	if err := s.seedPurchases(100); err != nil {
		return fmt.Errorf("failed to seed purchases: %w", err)
	}

	if err := s.seedModuleProgress(); err != nil {
		return fmt.Errorf("failed to seed module progress: %w", err)
	}

	if err := s.seedCertificates(30); err != nil {
		return fmt.Errorf("failed to seed certificates: %w", err)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func (s *Seeder) clearData() error {
	log.Println("Clearing existing data...")

	tables := []string{
		"certificates",
		"module_progresses",
		"purchases",
		"modules",
		"courses",
		"users",
	}

	for _, table := range tables {
		if err := s.db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) seedUsers(count int) error {
	log.Printf("Seeding %d users...", count)

	users := make([]models.User, count)
	roles := []string{"student", "admin"}

	for i := 0; i < count; i++ {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		users[i] = models.User{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Username:  faker.Name(),
			Email:     faker.Email(),
			Password:  string(hashedPassword),
			Role:      roles[rand.Intn(len(roles))],
			Balance:   float64(rand.Intn(1000)),
		}
	}

	return s.db.Create(&users).Error
}

func (s *Seeder) seedCourses(count int) error {
	log.Printf("Seeding %d courses...", count)

	courses := make([]models.Course, count)
	topics := [][]string{
		{"Programming", "Web Development", "Backend"},
		{"Data Science", "Machine Learning", "AI"},
		{"Mobile Development", "React Native", "Flutter"},
		{"DevOps", "Docker"},
		{"Design", "UI/UX", "Figma"},
		{"Marketing", "SEO", "Social Media"},
		{"Business", "Management", "Leadership"},
		{"Photography", "Editing", "Lightroom"},
		{"Music", "Production", "Theory"},
		{"Language", "English", "Spanish"},
	}

	for i := 0; i < count; i++ {
		selectedTopics := topics[rand.Intn(len(topics))]

		courses[i] = models.Course{
			Title:          faker.Sentence(),
			Description:    faker.Paragraph(),
			Instructor:     faker.Name(),
			Topics:         selectedTopics,
			Price:          float64(rand.Intn(500) + 50),
			ThumbnailImage: "https://i.imgflip.com/9grj9y.png?a487656",
		}
	}

	return s.db.Create(&courses).Error
}

func (s *Seeder) seedModules() error {
	log.Println("Seeding modules...")

	var courses []models.Course
	if err := s.db.Find(&courses).Error; err != nil {
		return err
	}

	var modules []models.Module

	for _, course := range courses {
		moduleCount := rand.Intn(8) + 3

		for j := 0; j < moduleCount; j++ {
			module := models.Module{
				CourseID:     course.ID,
				Title:        fmt.Sprintf("Module %d: %s", j+1, faker.Sentence()),
				Description:  faker.Paragraph(),
				Order:        j + 1,
				VideoContent: "https://cdn.mtdv.me/video/rick.mp4",
			}
			modules = append(modules, module)
		}
	}

	return s.db.Create(&modules).Error
}

func (s *Seeder) seedPurchases(count int) error {
	log.Printf("Seeding %d purchases...", count)

	var users []models.User
	var courses []models.Course

	if err := s.db.Find(&users).Error; err != nil {
		return err
	}
	if err := s.db.Find(&courses).Error; err != nil {
		return err
	}

	purchases := make([]models.Purchase, count)
	usedCombinations := make(map[string]bool)

	for i := 0; i < count; i++ {
		var userID, courseID uint
		var amount float64
		var combination string

		for {
			user := users[rand.Intn(len(users))]
			course := courses[rand.Intn(len(courses))]
			combination = fmt.Sprintf("%d-%d", user.ID, course.ID)

			if !usedCombinations[combination] {
				userID = user.ID
				courseID = course.ID
				amount = course.Price
				usedCombinations[combination] = true
				break
			}
		}

		purchases[i] = models.Purchase{
			UserID:   userID,
			CourseID: courseID,
			Amount:   amount,
		}
	}

	return s.db.Create(&purchases).Error
}

func (s *Seeder) seedModuleProgress() error {
	log.Println("Seeding module progress...")

	var purchases []models.Purchase
	if err := s.db.Preload("Course").Find(&purchases).Error; err != nil {
		return err
	}

	var progresses []models.ModuleProgress

	for _, purchase := range purchases {
		var modules []models.Module
		if err := s.db.Where("course_id = ?", purchase.CourseID).Find(&modules).Error; err != nil {
			continue
		}

		completedCount := rand.Intn(len(modules) + 1)

		for i, module := range modules {
			isCompleted := i < completedCount

			progress := models.ModuleProgress{
				UserID:      purchase.UserID,
				ModuleID:    module.ID,
				IsCompleted: isCompleted,
			}
			progresses = append(progresses, progress)
		}
	}

	if len(progresses) > 0 {
		return s.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "module_id"}}, // unique key
			DoUpdates: clause.AssignmentColumns([]string{"is_completed", "updated_at"}),
		}).Create(&progresses).Error
	}

	return nil
}

func (s *Seeder) seedCertificates(count int) error {
	log.Printf("Seeding %d certificates...", count)

	var completedCourses []struct {
		UserID   uint
		CourseID uint
	}

	query := `
		SELECT DISTINCT p.user_id, p.course_id 
		FROM purchases p
		WHERE NOT EXISTS (
			SELECT 1 FROM modules m 
			LEFT JOIN module_progresses mp ON m.id = mp.module_id AND mp.user_id = p.user_id
			WHERE m.course_id = p.course_id AND (mp.is_completed IS NULL OR mp.is_completed = false)
		)
		LIMIT ?
	`

	if err := s.db.Raw(query, count).Scan(&completedCourses).Error; err != nil {
		return err
	}

	for _, cc := range completedCourses {
		err := s.generateCertificate(cc.CourseID, cc.UserID)
		if err != nil {
			log.Println(err.Error())
		}
	}

	return nil
}

func SeedDatabase() error {
	if database.DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	seeder := NewSeeder(database.DB)
	return seeder.SeedAll()
}

func (s *Seeder) generateCertificate(courseId, userId uint) error {
	courseRepo := repositories.NewCourseRepository()
	userRepo := repositories.NewUserRepository()

	course, err := courseRepo.FindById(courseId)
	if err != nil {
		return err
	}

	user, err := userRepo.FindById(userId)
	if err != nil {
		return err
	}

	img, err := utils.GenerateCertificate(
		user.Username,
		course.Title,
		course.Instructor,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("cert_user%d_course%d.png", user.ID, courseId)
	filePath := filepath.Join("uploads/certificates", fileName)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	publicURL := os.Getenv("BASE_URL") + "uploads/certificates/" + fileName

	certificate := models.Certificate{
		UserID:   user.ID,
		CourseID: courseId,
		FileURL:  publicURL,
	}

	if err := courseRepo.CreateCourseCertificate(&certificate); err != nil {
		return err
	}

	return nil
}
