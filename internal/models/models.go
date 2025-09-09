package models

import (
	"github.com/google/uuid"
	"time"
)

type Vacancy struct {
	ID               uuid.UUID `gorm:"primaryKey;autoIncrement" json:"-"`
	Title            string    `gorm:"type:varchar(255)"`
	Requirements     string    `gorm:"type:text"`
	Responsibilities string    `gorm:"type:text"`
	Region           string    `gorm:"type:varchar(100)"`
	City             string    `gorm:"type:varchar(100)"`
	EmploymentType   string    `gorm:"type:varchar(50)"`  
	WorkSchedule     string    `gorm:"type:varchar(50)"`  
	Experience       string    `gorm:"type:varchar(50)"`  
	Education        string    `gorm:"type:varchar(100)"` 
	SalaryMin        int       `gorm:"type:integer"`
	SalaryMax        int       `gorm:"type:integer"`
	Languages        string    `gorm:"type:text"` 
	Skills           string    `gorm:"type:text"` 
	CreatedAt        time.Time
}

type Resume struct {
	ID           uuid.UUID `gorm:"primaryKey;autoIncrement" json:"-"`
	CandidateID  uuid.UUID `gorm:"type:uuid"`
	Text         string    `gorm:"type:text"`
	ParsedData   string    `gorm:"type:jsonb"`
	FileURL      string    `gorm:"type:text"`
	Experience   int       `gorm:"type:integer"` 
	Education    string    `gorm:"type:varchar(100)"`
	Skills       string    `gorm:"type:text"`
	Languages    string    `gorm:"type:text"`
	SalaryExpect int       `gorm:"type:integer"`
	CreatedAt    time.Time
}

type AnalysisResult struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid" json:"-"`
	ResumeID   uuid.UUID `gorm:"type:uuid"`
	VacancyID  uuid.UUID `gorm:"type:uuid"`
	MatchScore float64   `gorm:"type:decimal(5,2)"`
	Details    string    `gorm:"type:jsonb;default:'{}'"`
	CreatedAt  time.Time
}

type AnalysisDetail struct {
	ID               uuid.UUID `gorm:"primaryKey;type:uuid" json:"-"`
	AnalysisResultID uuid.UUID `gorm:"type:uuid"`
	Category         string    `gorm:"type:varchar(100)"` 
	Criteria         string    `gorm:"type:text"`         
	ResumeValue      string    `gorm:"type:text"`         
	VacancyValue     string    `gorm:"type:text"`         
	MatchScore       float64   `gorm:"type:decimal(3,2)"` 
	Weight           float64   `gorm:"type:decimal(3,2)"` 
	CreatedAt        time.Time
}
