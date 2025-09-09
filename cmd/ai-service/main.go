package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Vacancy struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Requirements string   `json:"requirements"`
	Tags         []string `json:"tags"`
}

type QuestionRequest struct {
	Vacancy Vacancy `json:"vacancy"`
}

type AnswerAnalysisRequest struct {
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
	Vacancy  Vacancy `json:"vacancy"`
}

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

var openRouterAPIKey = os.Getenv("OPENROUTER_API_KEY")

func main() {
	if openRouterAPIKey == "" {
		log.Fatal("❌ OPENROUTER_API_KEY не установлен! Сервис не сможет работать.")
	}

	r := gin.Default()

	r.POST("/generate-questions", generateQuestionsHandler)
	r.POST("/analyze-answer", analyzeAnswerHandler)
	r.POST("/analyze-resume", analyzeResumeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("🚀 AI Service запущен на порту %s", port)
	r.Run(":" + port)
}

func generateQuestionsHandler(c *gin.Context) {
	var req QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Ошибка парсинга запроса /generate-questions: %v", err)
		c.JSON(400, gin.H{"error": "Неверный формат запроса"})
		return
	}

	log.Printf("📥 Запрос /generate-questions: Vacancy=%s", req.Vacancy.Title)

	prompt := fmt.Sprintf(`
		На основе вакансии "%s" сгенерируй 4 технических вопроса для собеседования.
		Описание: %s
		Требования: %s
		Навыки: %s

		Верни только вопросы в формате JSON массива: ["вопрос1", "вопрос2", "вопрос3", "вопрос4"]
	`,
		req.Vacancy.Title,
		req.Vacancy.Description,
		req.Vacancy.Requirements,
		strings.Join(req.Vacancy.Tags, ", "),
	)

	raw, err := callOpenRouter(prompt)
	if err != nil {
		log.Printf("❌ Ошибка вызова OpenRouter: %v", err)
		c.JSON(500, gin.H{"error": "Ошибка генерации вопросов"})
		return
	}

	log.Printf("📦 Ответ от OpenRouter (raw): %s", raw)

	var questions []string
	if err := json.Unmarshal([]byte(raw), &questions); err != nil {
		log.Printf("❌ Ошибка парсинга JSON массива вопросов: %v", err)
		c.JSON(500, gin.H{"error": "ИИ вернул некорректный формат вопросов", "raw": raw})
		return
	}

	log.Printf("✅ Сгенерировано %d вопросов", len(questions))
	c.JSON(200, gin.H{"questions": questions})
}

func analyzeAnswerHandler(c *gin.Context) {
	var req AnswerAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Ошибка парсинга запроса /analyze-answer: %v", err)
		c.JSON(400, gin.H{"error": "Неверный формат запроса"})
		return
	}

	log.Printf("📥 Запрос /analyze-answer: Вопрос='%s', Ответ len=%d",
		req.Question, len(req.Answer))

	prompt := fmt.Sprintf(`
		Проанализируй ответ кандидата на вопрос интервью и оцени его по шкале от 0 до 1.

		Вакансия: %s
		Требования: %s
		Навыки: %s

		Вопрос: %s
		Ответ кандидата: %s

		Верни ответ в формате JSON: {"score": 0.85, "feedback": "конструктивный фидбек"}
	`,
		req.Vacancy.Title,
		req.Vacancy.Requirements,
		strings.Join(req.Vacancy.Tags, ", "),
		req.Question,
		req.Answer,
	)

	raw, err := callOpenRouter(prompt)
	if err != nil {
		log.Printf("❌ Ошибка вызова OpenRouter: %v", err)
		c.JSON(500, gin.H{"error": "Ошибка анализа ответа"})
		return
	}

	log.Printf("📦 Ответ от OpenRouter (raw): %s", raw)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		log.Printf("❌ Ошибка парсинга анализа JSON: %v", err)
		c.JSON(500, gin.H{"error": "ИИ вернул некорректный формат анализа", "raw": raw})
		return
	}

	log.Printf("✅ Анализ ответа: %+v", result)
	c.JSON(200, gin.H{"result": result})
}

func callOpenRouter(prompt string) (string, error) {
	requestBody := OpenRouterRequest{
		Model: "openai/gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка при маршалинге JSON: %v", err)
	}

	log.Printf("➡️ Отправка запроса в OpenRouter")

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openRouterAPIKey)
	req.Header.Set("HTTP-Referer", "https://hr-avatar.com")
	req.Header.Set("X-Title", "HR Avatar AI Service")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	log.Printf("⬅️ Ответ OpenRouter (status %d): %s", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ошибка API OpenRouter: %s", string(body))
	}

	var response OpenRouterResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге JSON-ответа: %v", err)
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("пустой ответ от ИИ")
}

type ResumeAnalysisRequest struct {
	Vacancy Vacancy `json:"vacancy"`
	Resume  struct {
		Name    string `json:"name"`
		Content string `json:"content"`
		Type    string `json:"type"`
	} `json:"resume"`
}

func analyzeResumeHandler(c *gin.Context) {
	var req ResumeAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Ошибка парсинга запроса /analyze-resume: %v", err)
		c.JSON(400, gin.H{"error": "Неверный формат запроса"})
		return
	}

	log.Printf("📥 Запрос /analyze-resume: Vacancy=%s, Resume=%s", req.Vacancy.Title, req.Resume.Name)

	prompt := fmt.Sprintf(`
        Проанализируй резюме кандидата и сравни его с вакансией "%s".
        Вакансия:
        Описание: %s
        Требования: %s
        Навыки: %s

        Резюме кандидата:
        %s

        Оцени, насколько кандидат подходит под вакансию, и верни JSON:
        {"match_score": 75, "feedback": "Кандидат соответствует большинству требований, но не имеет опыта с AWS"}
    `,
		req.Vacancy.Title,
		req.Vacancy.Description,
		req.Vacancy.Requirements,
		strings.Join(req.Vacancy.Tags, ", "),
		req.Resume.Content,
	)

	raw, err := callOpenRouter(prompt)
	if err != nil {
		log.Printf("❌ Ошибка вызова OpenRouter: %v", err)
		c.JSON(500, gin.H{"error": "Ошибка анализа резюме"})
		return
	}

	log.Printf("📦 Ответ от OpenRouter (raw): %s", raw)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		log.Printf("❌ Ошибка парсинга JSON от OpenRouter: %v", err)
		c.JSON(500, gin.H{"error": "ИИ вернул некорректный формат анализа", "raw": raw})
		return
	}

	c.JSON(200, result)
}
