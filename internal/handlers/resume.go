package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/unidoc/unioffice/common/license"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/moverq1337/VTBHack/internal/models"
	"github.com/moverq1337/VTBHack/internal/pb"
	"github.com/moverq1337/VTBHack/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/unidoc/unioffice/document"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var log = logrus.New()

func GetVacancies(c *gin.Context, db *gorm.DB) {
	var vacancies []models.Vacancy
	if err := db.Find(&vacancies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения вакансий из базы данных"})
		return
	}

	var result []map[string]interface{}
	for _, v := range vacancies {
		var tags []string
		if v.Skills != "" {
			tags = strings.Split(v.Skills, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		vacancyMap := map[string]interface{}{
			"id":          v.ID,
			"title":       v.Title,
			"company":     "IT Company",
			"location":    fmt.Sprintf("%s, %s", v.Region, v.City),
			"description": v.Requirements,
			"tags":        tags,
			"experience":  v.Experience,
			"salary":      fmt.Sprintf("%d - %d руб.", v.SalaryMin, v.SalaryMax),
		}
		result = append(result, vacancyMap)
	}

	c.JSON(http.StatusOK, result)
}

func AnalyzeResumeForVacancy(c *gin.Context, db *gorm.DB) {
	type Request struct {
		Vacancy map[string]interface{} `json:"vacancy"`
		Resume  struct {
			Name    string `json:"name"`
			Content string `json:"content"`
			Type    string `json:"type"`
		} `json:"resume"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	resumeText := strings.ToLower(req.Resume.Content)
	vacancyTags := []string{}

	if tags, ok := req.Vacancy["tags"].([]interface{}); ok {
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				vacancyTags = append(vacancyTags, tagStr)
			}
		}
	}

	matchCount := 0
	for _, tag := range vacancyTags {
		if strings.Contains(resumeText, strings.ToLower(tag)) {
			matchCount++
		}
	}

	matchScore := 0.0
	if len(vacancyTags) > 0 {
		matchScore = float64(matchCount) / float64(len(vacancyTags)) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"match_score":    fmt.Sprintf("%.1f", matchScore),
		"matched_skills": matchCount,
		"total_skills":   len(vacancyTags),
	})
}

func UploadResume(c *gin.Context, db *gorm.DB) {
	log.Info("Начало загрузки резюме DOCX")

	file, err := c.FormFile("resume")
	if err != nil {
		log.WithError(err).Error("Ошибка загрузки файла")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен: " + err.Error()})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".docx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поддерживается только DOCX формат"})
		return
	}

	candidateID := uuid.New()
	filePath := filepath.Join("uploads", candidateID.String()+ext)
	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания директории"})
		return
	}
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.WithError(err).Error("Ошибка сохранения файла")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}
	defer os.Remove(filePath)

	text, err := extractTextFromDOCX(filePath)
	if err != nil {
		log.WithError(err).Error("Ошибка извлечения текста из DOCX")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка извлечения текста из DOCX: " + err.Error()})
		return
	}

	diskURL, err := utils.UploadToYandexDisk(filePath, file.Filename)
	if err != nil {
		log.WithError(err).Error("Ошибка загрузки на Яндекс.Диск")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки на Яндекс.Диск: " + err.Error()})
		return
	}

	resume := models.Resume{
		ID:          uuid.New(),
		CandidateID: candidateID,
		Text:        text,
		ParsedData:  "{}",
		FileURL:     diskURL,
	}

	grpchost := "scoring-service:50051"
	conn, err := grpc.NewClient(grpchost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithError(err).Error("Ошибка gRPC-соединения для парсинга")
	} else {
		defer conn.Close()

		client := pb.NewNLPServiceClient(conn)
		parseResp, err := client.ParseResume(context.Background(), &pb.ParseRequest{
			Text: text,
		})

		if err != nil {
			log.WithError(err).Error("Ошибка парсинга резюме")
		} else {
			resume.ParsedData = parseResp.ParsedData
		}
	}

	if err := db.Create(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения в БД"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"candidate_id": candidateID.String(),
		"file_url":     diskURL,
		"resume_id":    resume.ID.String(),
		"text_preview": truncateText(text, 200),
	})
}

func extractTextFromDOCX(filePath string) (string, error) {
	apiKey := os.Getenv("UNIDOC_LICENSE_API_KEY")
	if apiKey == "" {
		log.Fatal("UNIDOC_LICENSE_API_KEY environment variable not set")
	}

	err := license.SetMeteredKey(apiKey)
	if err != nil {
		log.Fatalf("Ошибка инициализации UniDoc license: %v", err)
	}
	doc, err := document.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия DOCX файла: %v", err)
	}

	var text string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text += run.Text()
		}
		text += "\n"
	}

	for _, tbl := range doc.Tables() {
		for _, row := range tbl.Rows() {
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					for _, run := range para.Runs() {
						text += run.Text() + " "
					}
				}
			}
			text += "\n"
		}
	}

	return text, nil
}

func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

func UploadVacancy(c *gin.Context, db *gorm.DB) {
	type VacancyRequest struct {
		Title            string `json:"title"`
		Requirements     string `json:"requirements"`
		Responsibilities string `json:"responsibilities"`
		Region           string `json:"region"`
		City             string `json:"city"`
		EmploymentType   string `json:"employment_type"`
		WorkSchedule     string `json:"work_schedule"`
		Experience       string `json:"experience"`
		Education        string `json:"education"`
		SalaryMin        int    `json:"salary_min"`
		SalaryMax        int    `json:"salary_max"`
		Languages        string `json:"languages"`
		Skills           string `json:"skills"`
	}

	var req VacancyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	vacancy := models.Vacancy{
		ID:               uuid.New(),
		Title:            req.Title,
		Requirements:     req.Requirements,
		Responsibilities: req.Responsibilities,
		Region:           req.Region,
		City:             req.City,
		EmploymentType:   req.EmploymentType,
		WorkSchedule:     req.WorkSchedule,
		Experience:       req.Experience,
		Education:        req.Education,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		Languages:        req.Languages,
		Skills:           req.Skills,
	}
	if err := db.Create(&vacancy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения в БД"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vacancy_id": vacancy.ID.String(),
		"title":      vacancy.Title,
	})
}

func AnalyzeResume(c *gin.Context, db *gorm.DB) {
	type AnalyzeRequest struct {
		ResumeID  uuid.UUID `json:"resume_id"`
		VacancyID uuid.UUID `json:"vacancy_id"`
	}

	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	var resume models.Resume
	if err := db.First(&resume, "id = ?", req.ResumeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Резюме не найдено"})
		return
	}

	var vacancy models.Vacancy
	if err := db.First(&vacancy, "id = ?", req.VacancyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Вакансия не найдена"})
		return
	}

	grpchost := "scoring-service:50051"

	conn, err := grpc.NewClient(grpchost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithError(err).Error("Ошибка gRPC-соединения")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к сервису анализа"})
		return
	}
	defer conn.Close()

	client := pb.NewNLPServiceClient(conn)

	matchResp, err := client.MatchResumeVacancy(context.Background(), &pb.MatchRequest{
		ResumeText:  resume.Text,
		VacancyText: fmt.Sprintf("%s %s %s %s", vacancy.Title, vacancy.Requirements, vacancy.Responsibilities, vacancy.Skills),
	})
	if err != nil {
		log.WithError(err).Error("Ошибка сопоставления")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сопоставления с вакансией"})
		return
	}

	parseResp, err := client.ParseResume(context.Background(), &pb.ParseRequest{
		Text: resume.Text,
	})
	if err != nil {
		log.WithError(err).Error("Ошибка парсинга резюме")
	}

	var parsedData map[string]interface{}
	if err := json.Unmarshal([]byte(parseResp.ParsedData), &parsedData); err != nil {
		log.WithError(err).Error("Ошибка разбора JSON данных парсинга")
		parsedData = make(map[string]interface{})
	}

	analysisResult := models.AnalysisResult{
		ID:         uuid.New(),
		ResumeID:   resume.ID,
		VacancyID:  vacancy.ID,
		MatchScore: float64(matchResp.Score),
		Details:    parseResp.ParsedData,
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&analysisResult).Error; err != nil {
		log.WithError(err).Error("Ошибка сохранения результатов анализа")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения результатов"})
		return
	}

	if skills, ok := parsedData["skills"].(map[string]interface{}); ok {
		for category, skillList := range skills {
			if skillsArr, ok := skillList.([]interface{}); ok {
				for _, skill := range skillsArr {
					analysisDetail := models.AnalysisDetail{
						ID:               uuid.New(),
						AnalysisResultID: analysisResult.ID,
						Category:         "skills",
						Criteria:         category,
						ResumeValue:      fmt.Sprintf("%v", skill),
						VacancyValue:     "",
						MatchScore:       0.8,
						Weight:           0.3,
						CreatedAt:        time.Now(),
					}
					if err := db.Create(&analysisDetail).Error; err != nil {
						log.WithError(err).Error("Ошибка сохранения деталей анализа")
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"analysis_id":  analysisResult.ID.String(),
		"resume_id":    resume.ID.String(),
		"vacancy_id":   vacancy.ID.String(),
		"match_score":  fmt.Sprintf("%.2f%%", matchResp.Score*100),
		"candidate_id": resume.CandidateID.String(),
		"created_at":   analysisResult.CreatedAt,
		"details":      parsedData,
	})
}
