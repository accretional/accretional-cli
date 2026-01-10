package file

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// extractTextFromFile extracts text content from file data based on file type.
// Returns the extracted text and detected MIME type.
// Currently supports plain text files. Can be extended for PDF, DOCX, etc.
func extractTextFromFile(fileData []byte, fileName string) (textContent string, mimeType string, err error) {
	// Detect file type from extension
	ext := strings.ToLower(filepath.Ext(fileName))

	// Detect MIME type based on extension
	switch ext {
	case ".txt", ".text", ".md", ".markdown", ".log":
		mimeType = "text/plain"
		// For plain text files, decode from bytes directly
		textContent = string(fileData)
		return textContent, mimeType, nil

	case ".json":
		mimeType = "application/json"
		// JSON files are text-based
		textContent = string(fileData)
		return textContent, mimeType, nil

	case ".xml", ".html", ".htm":
		mimeType = "text/" + ext[1:] // text/xml, text/html
		textContent = string(fileData)
		return textContent, mimeType, nil

	case ".csv":
		mimeType = "text/csv"
		textContent = string(fileData)
		return textContent, mimeType, nil

	case ".yaml", ".yml":
		mimeType = "text/yaml"
		textContent = string(fileData)
		return textContent, mimeType, nil

	case ".go", ".js", ".ts", ".py", ".java", ".cpp", ".c", ".h", ".rs", ".rb", ".php", ".sh", ".bash", ".proto":
		mimeType = "text/x-" + ext[1:] // text/x-go, text/x-js, text/x-proto, etc.
		textContent = string(fileData)
		return textContent, mimeType, nil

	case ".pdf":
		mimeType = "application/pdf"
		// TODO: Add PDF text extraction library (e.g., gofpdf, unidoc)
		// For now, return empty - PDF extraction requires external library
		return "", mimeType, fmt.Errorf("PDF text extraction not yet implemented")

	case ".docx":
		mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		// TODO: Add DOCX text extraction library
		// For now, return empty
		return "", mimeType, fmt.Errorf("DOCX text extraction not yet implemented")

	case ".doc":
		mimeType = "application/msword"
		// TODO: Add DOC text extraction library
		return "", mimeType, fmt.Errorf("DOC text extraction not yet implemented")

	case ".xlsx":
		mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		// TODO: Add XLSX text extraction library
		return "", mimeType, fmt.Errorf("XLSX text extraction not yet implemented")

	default:
		// Try to detect if it's text by checking if it's valid UTF-8
		// and doesn't contain too many null bytes
		if isValidText(fileData) {
			mimeType = "text/plain"
			textContent = string(fileData)
			return textContent, mimeType, nil
		}

		// Unknown binary file type
		mimeType = "application/octet-stream"
		return "", mimeType, fmt.Errorf("cannot extract text from binary file type: %s", ext)
	}
}

// isValidText checks if data appears to be text (valid UTF-8, not too many null bytes)
func isValidText(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// Check if it's valid UTF-8
	if !isValidUTF8(data) {
		return false
	}

	// Check for too many null bytes (more than 5% suggests binary)
	nullCount := 0
	for _, b := range data {
		if b == 0 {
			nullCount++
		}
	}
	if nullCount > len(data)/20 {
		return false
	}

	return true
}

// isValidUTF8 checks if data is valid UTF-8
func isValidUTF8(data []byte) bool {
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			return false
		}
		data = data[size:]
	}
	return true
}

// detectMimeType detects MIME type from file extension without extracting text content.
// This is used when text extraction is disabled.
func detectMimeType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".txt", ".text", ".md", ".markdown", ".log":
		return "text/plain"
	case ".json":
		return "application/json"
	case ".xml":
		return "text/xml"
	case ".html", ".htm":
		return "text/html"
	case ".csv":
		return "text/csv"
	case ".yaml", ".yml":
		return "text/yaml"
	case ".go", ".js", ".ts", ".py", ".java", ".cpp", ".c", ".h", ".rs", ".rb", ".php", ".sh", ".bash", ".proto":
		return "text/x-" + ext[1:]
	case ".pdf":
		return "application/pdf"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".doc":
		return "application/msword"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	default:
		return "application/octet-stream"
	}
}
