package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
)

func PrepopulateTags() {
	tags := []models.Tag{
		{Title: "Nightlife and Parties"},
		{Title: "Family-Friendly"},
		{Title: "Charity and Fundraising"},
		{Title: "Tech Conferences"},
		{Title: "Workshops"},
		{Title: "Theater Performances"},
		{Title: "Comedy Shows"},
		{Title: "Sports Events"},
		{Title: "Food Festivals"},
		{Title: "Art Exhibitions"},
		{Title: "Music Concerts"},
		{Title: "Travel and Tourism"},
		{Title: "Fashion Shows"},
		{Title: "Cultural Festivals"},
		{Title: "Gaming and Esports"},
		{Title: "Film Screenings"},
		{Title: "Networking Events"},
		{Title: "Educational Seminars"},
		{Title: "Wellness and Fitness"},
		{Title: "Outdoor Adventures"},
	}

	for _, tag := range tags {
		db.DB.Create(&tag)
	}
}

func GetTags() (*[]models.Tag, error) {
	var tags []models.Tag

	error := db.DB.Model(&models.Tag{}).Find(&tags).Error
	if error != nil {
		return nil, error
	}

	return &tags, error
}
