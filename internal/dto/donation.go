package dto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type CreateDonationRequest struct {
	PayeeUserID uuid.UUID  `json:"payee_user_id" form:"payee_user_id"`
	PayerUserID *uuid.UUID `json:"payer_user_id,omitempty" form:"payer_user_id"`

	PayerName string `json:"payer_name" form:"payer_name"`
	Email     string `json:"email" form:"email"`
	Message   string `json:"message" form:"message"`

	MediaType string `json:"media_type" form:"media_type"`

	// Media input dari user
	MediaURL          *string `json:"media_url,omitempty" form:"media_url"`
	MediaStartSeconds *int32  `json:"media_start_seconds,omitempty" form:"media_start_seconds"`

	// Transaction
	Amount         int64  `json:"amount" form:"amount"`
	PaymentChannel string `json:"payment_channel" form:"payment_channel"`
}

// Validate validates all fields including business rules
func (r *CreateDonationRequest) Validate() error {
	// Basic field validations
	if r.PayeeUserID == uuid.Nil {
		return errors.New("PayeeUserID is required")
	}

	if r.PayerUserID != nil && r.PayeeUserID == *r.PayerUserID {
		return errors.New("PayeeUserID cannot be the same as PayerUserID")
	}

	if r.PayerName == "" {
		return errors.New("PayerName is required")
	}
	if len(r.PayerName) < 1 || len(r.PayerName) > 20 {
		return errors.New("PayerName must be between 1 and 20 characters")
	}

	if r.Email == "" {
		return errors.New("Email is required")
	}
	if !isValidEmail(r.Email) {
		return errors.New("Email is not valid")
	}

	if r.Message == "" {
		return errors.New("Message is required")
	}
	if len(r.Message) < 1 || len(r.Message) > 300 {
		return errors.New("Message must be between 1 and 300 characters")
	}

	if r.MediaType == "" {
		return errors.New("MediaType is required")
	}
	if !isValidMediaType(r.MediaType) {
		return errors.New("MediaType must be one of: TEXT, YOUTUBE, GIF")
	}

	if r.Amount < 1000 {
		return errors.New("Amount must be at least 1000")
	}

	if r.Amount > 100_000_000 {
		return errors.New("Amount cannot exceed 100,000,000")
	}

	if r.PaymentChannel != "qris" {
		return errors.New("PaymentChannel must be qris")
	}

	// Media type specific validations
	return r.validateMediaFields()
}

// validateMediaFields validates media fields based on MediaType
func (r *CreateDonationRequest) validateMediaFields() error {
	switch r.MediaType {
	case "TEXT":
		return r.validateTextMedia()
	case "YOUTUBE":
		return r.validateYouTubeMedia()
	case "GIF":
		return r.validateGIFMedia()
	default:
		return fmt.Errorf("unsupported media type: %s", r.MediaType)
	}
}

// validateTextMedia ensures TEXT type doesn't have media fields
func (r *CreateDonationRequest) validateTextMedia() error {
	if r.MediaURL != nil && *r.MediaURL != "" {
		return errors.New("MediaURL is not allowed when MediaType is TEXT")
	}
	if r.MediaStartSeconds != nil {
		return errors.New("MediaStartSeconds is not allowed when MediaType is TEXT")
	}
	return nil
}

// validateYouTubeMedia ensures YOUTUBE type has required media fields
func (r *CreateDonationRequest) validateYouTubeMedia() error {
	// MediaURL is required for YOUTUBE
	if r.MediaURL == nil || *r.MediaURL == "" {
		return errors.New("MediaURL is required when MediaType is YOUTUBE")
	}

	// Validate URL format
	if !isValidURL(*r.MediaURL) {
		return errors.New("MediaURL must be a valid URL")
	}

	// Validate YouTube URL format
	if !isValidYouTubeURL(*r.MediaURL) {
		return errors.New("MediaURL must be a valid YouTube URL")
	}

	// MediaStartSeconds is required for YOUTUBE
	if r.MediaStartSeconds == nil {
		return errors.New("MediaStartSeconds is required when MediaType is YOUTUBE")
	}

	// Additional validation: reasonable time limit
	if *r.MediaStartSeconds < 0 {
		return errors.New("MediaStartSeconds must be greater than or equal to 0")
	}

	// Optional: Max duration check (24 hours)
	maxSeconds := int32(86400)
	if *r.MediaStartSeconds > maxSeconds {
		return fmt.Errorf("MediaStartSeconds cannot exceed %d seconds (24 hours)", maxSeconds)
	}

	return nil
}

// validateGIFMedia ensures GIF type has MediaURL but not MediaStartSeconds
func (r *CreateDonationRequest) validateGIFMedia() error {
	// MediaURL is required for GIF
	if r.MediaURL == nil || *r.MediaURL == "" {
		return errors.New("MediaURL is required when MediaType is GIF")
	}

	// Validate URL format
	if !isValidURL(*r.MediaURL) {
		return errors.New("MediaURL must be a valid URL")
	}

	// MediaStartSeconds is not allowed for GIF
	if r.MediaStartSeconds != nil {
		return errors.New("MediaStartSeconds is not allowed when MediaType is GIF")
	}

	// Optional: Validate GIF URL format
	if !isValidGIFURL(*r.MediaURL) {
		return errors.New("MediaURL must be a valid GIF URL (e.g., from Giphy or Tenor)")
	}

	return nil
}

// Helper validation functions

func isValidEmail(email string) bool {
	// Simple email validation
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

func isValidMediaType(mediaType string) bool {
	validTypes := []string{"TEXT", "YOUTUBE", "GIF"}
	for _, t := range validTypes {
		if mediaType == t {
			return true
		}
	}
	return false
}

func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func isValidYouTubeURL(url string) bool {
	url = strings.ToLower(url)
	validDomains := []string{
		"youtube.com",
		"www.youtube.com",
		"m.youtube.com",
		"youtu.be",
	}

	for _, domain := range validDomains {
		if strings.Contains(url, domain) {
			return true
		}
	}
	return false
}

func isValidGIFURL(url string) bool {
	url = strings.ToLower(url)

	// Check if it's from common GIF providers
	validDomains := []string{
		"giphy.com",
		"tenor.com",
		"media.giphy.com",
		"media.tenor.com",
	}

	for _, domain := range validDomains {
		if strings.Contains(url, domain) {
			return true
		}
	}

	// Or check if URL ends with .gif
	return strings.HasSuffix(url, ".gif")
}
