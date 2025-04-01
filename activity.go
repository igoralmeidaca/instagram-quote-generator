package quote

import (
	"context"
	"database/sql"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"

	_ "github.com/lib/pq"
)

const (
	imgWidth       = 1080 // Instagram standard width
	imgHeight      = 1080 // Instagram standard height
	blurSigma      = 3.0  // Adjust this value for more/less blur
	shadowX        = 3    // X offset for shadow
	shadowY        = 3    // Y offset for shadow
	quoteFontSize  = 30
	authorFontSize = 20
)

var (
	dbURL          = os.Getenv("DATABASE_URL")
	outputDir      = os.Getenv("OUTPUT_DIR")
	bgImagePath    = os.Getenv("BG_IMAGE_PATH")
	quoteFontPath  = os.Getenv("QUOTE_FONT_PATH")
	authorFontPath = os.Getenv("AUTHOR_FONT_PATH")
)

func GenerateText(ctx context.Context, input GenerateTextInput) (GenerateTextOutput, error) {
	if dbURL == "" {
		return GenerateTextOutput{}, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return GenerateTextOutput{}, err
	}
	defer db.Close()

	// Try to get an unused quote
	quote, err := getUnusedQuote(db)
	if err != nil {
		return GenerateTextOutput{}, err
	}

	return GenerateTextOutput{Text: quote.Text, Author: quote.Author, Reference: quote.Reference}, nil
}

func GenerateImage(ctx context.Context, input GenerateTextOutput) (GenerateImageOutput, error) {
	err := validateImageVariables()
	if err != nil {
		return GenerateImageOutput{}, err
	}

	// Load the background image
	bgImage, err := gg.LoadImage(bgImagePath)
	if err != nil {
		return GenerateImageOutput{}, fmt.Errorf("failed to load background image: %w", err)
	}

	// Apply Gaussian blur to the image
	blurredImage := imaging.Blur(bgImage, blurSigma)

	// Get background image size
	bgWidth := blurredImage.Bounds().Dx()
	bgHeight := blurredImage.Bounds().Dy()

	// Create an image context
	dc := gg.NewContext(bgWidth, bgHeight)

	// Draw the blurred background image onto the context
	dc.DrawImage(blurredImage, 0, 0)

	// Load font for text
	if err := dc.LoadFontFace(quoteFontPath, quoteFontSize); err != nil {
		return GenerateImageOutput{}, fmt.Errorf("failed to load font: %w", err)
	}

	// Set positions
	textX := float64(bgWidth) / 2
	textY := float64(bgHeight) / 2.5

	// Draw shadow (darker text slightly offset)
	dc.SetColor(color.Black) // Shadow color
	dc.DrawStringWrapped(input.Text, textX+shadowX, textY+shadowY, 0.5, 0.5, float64(bgWidth)-100, 1.5, gg.AlignCenter)

	// Draw main text
	dc.SetColor(color.White) // Main text color
	dc.DrawStringWrapped(input.Text, textX, textY, 0.5, 0.5, float64(bgWidth)-100, 1.5, gg.AlignCenter)

	// Load font for author & reference
	if err := dc.LoadFontFace(authorFontPath, authorFontSize); err != nil {
		return GenerateImageOutput{}, fmt.Errorf("failed to load font: %w", err)
	}

	// Draw shadow for author text
	authorX := float64(bgWidth) / 2
	authorY := float64(bgHeight) * 0.9
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(fmt.Sprintf("- %s (%s)", input.Author, input.Reference), authorX+shadowX, authorY+shadowY, 0.5, 0.5)

	// Draw main author text
	dc.SetColor(color.White)
	dc.DrawStringAnchored(fmt.Sprintf("- %s (%s)", input.Author, input.Reference), authorX, authorY, 0.5, 0.5)

	// Save the image
	fileName := fmt.Sprintf("quote_%d.png", time.Now().Unix())
	filePath := filepath.Join(outputDir, fileName)
	if err := dc.SavePNG(filePath); err != nil {
		return GenerateImageOutput{}, fmt.Errorf("failed to save image: %w", err)
	}

	return GenerateImageOutput{fileName}, nil
}

func PublishInstagramPost(input GenerateImageOutput) (string, error) {
	// TODO implement
	return "", nil
}

func validateImageVariables() error {
	if outputDir == "" {
		return fmt.Errorf("OUTPUT_DIR is not set")
	}

	if bgImagePath == "" {
		return fmt.Errorf("BG_IMAGE_PATH is not set")
	}

	if quoteFontPath == "" {
		return fmt.Errorf("QUOTE_FONT_PATH is not set")
	}

	if authorFontPath == "" {
		return fmt.Errorf("AUTHOR_FONT_PATH is not set")
	}
	return nil
}

// getUnusedQuote fetches an unused quote and marks it as used
func getUnusedQuote(db *sql.DB) (*Quote, error) {
	var quote Quote
	err := db.QueryRow(`
		SELECT id, text, author, reference FROM quotes 
		WHERE used = FALSE 
		ORDER BY id LIMIT 1
	`).Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Reference)

	if err != nil {
		return nil, err
	}

	// Mark the selected quote as used
	_, err = db.Exec("UPDATE quotes SET used = TRUE WHERE id = $1", quote.ID)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}
