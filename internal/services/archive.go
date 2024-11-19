package services

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"dbc/internal/utils"

	model "dbc/internal/domain"
)

type ArchiveService interface {
	GetArchiveInfo(file multipart.File, filename string) (model.ArchiveInfo, error)
	CreateArchive(files []*multipart.FileHeader) ([]byte, error)
}

type archiveService struct{}

func NewArchiveService() ArchiveService {
	return &archiveService{}
}

func (s *archiveService) GetArchiveInfo(file multipart.File, filename string) (model.ArchiveInfo, error) {
	if filepath.Ext(filename) != ".zip" {
		return model.ArchiveInfo{}, errors.New("only .zip files are supported")
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return model.ArchiveInfo{}, fmt.Errorf("failed to read the file: %w", err)
	}

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return model.ArchiveInfo{}, errors.New("invalid archive format")
	}

	var totalSize float64
	var files []model.FileInfo
	for _, f := range reader.File {
		size := float64(f.UncompressedSize64)
		totalSize += size
		mimeType := utils.GetMimeType(f.Name)

		files = append(files, model.FileInfo{
			FilePath: f.Name,
			Size:     size,
			MimeType: mimeType,
		})
	}

	return model.ArchiveInfo{
		Filename:    filename,
		ArchiveSize: float64(buf.Len()),
		TotalSize:   totalSize,
		TotalFiles:  len(reader.File),
		Files:       files,
	}, nil
}

func (s *archiveService) CreateArchive(files []*multipart.FileHeader) ([]byte, error) {
	if len(files) == 0 {
		return nil, errors.New("at least one file is required")
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	validTypes := []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/xml",
		"image/jpeg",
		"image/png",
	}

	for _, file := range files {
		mimeType := utils.GetMimeType(file.Filename)
		if !utils.IsValidMimeType(mimeType, validTypes) {
			return nil, fmt.Errorf("invalid file type: %s", file.Filename)
		}

		zipFile, err := zipWriter.Create(file.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create zip entry: %w", err)
		}

		uploadedFile, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}

		if _, err := io.Copy(zipFile, uploadedFile); err != nil {
			uploadedFile.Close()
			return nil, fmt.Errorf("failed to write file to archive: %w", err)
		}
		uploadedFile.Close()
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize the archive: %w", err)
	}

	return buf.Bytes(), nil
}
