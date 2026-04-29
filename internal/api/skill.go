package api

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type SkillService struct {
	client *Client
}

func NewSkillService(client *Client) *SkillService {
	return &SkillService{client: client}
}

func (s *SkillService) DownloadAndExtract(projectID, skillID, authHeader, outputDir string) error {
	projectID = strings.TrimSpace(projectID)
	skillID = strings.TrimSpace(skillID)
	if projectID == "" {
		return fmt.Errorf("project ID is required")
	}
	if skillID == "" {
		return fmt.Errorf("skill ID is required")
	}

	path := fmt.Sprintf("/api/v1/%s/skills/%s/.well-known/skills/archive.zip", projectID, skillID)
	return s.downloadAndExtract(path, authHeader, outputDir)
}

func (s *SkillService) DownloadMarketplaceAndExtract(skillID, authHeader, outputDir string) error {
	skillID = strings.TrimSpace(skillID)
	if skillID == "" {
		return fmt.Errorf("skill ID is required")
	}

	path := fmt.Sprintf("/api/v1/marketplace/skills/%s/.well-known/skills/archive.zip", skillID)
	return s.downloadAndExtract(path, authHeader, outputDir)
}

func (s *SkillService) downloadAndExtract(path, authHeader, outputDir string) error {
	body, statusCode, err := s.client.Get(path, authHeader)
	if err != nil {
		return fmt.Errorf("download archive: %w", err)
	}
	if statusCode != 200 {
		return fmt.Errorf("archive request failed with status %d: %s", statusCode, strings.TrimSpace(string(body)))
	}

	if err = os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	return unzipArchive(body, outputDir)
}

func unzipArchive(archiveBody []byte, outputDir string) error {
	reader, err := zip.NewReader(bytes.NewReader(archiveBody), int64(len(archiveBody)))
	if err != nil {
		return fmt.Errorf("open downloaded archive: %w", err)
	}

	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("resolve extraction path: %w", err)
	}

	for _, file := range reader.File {
		cleanName := filepath.Clean(file.Name)
		if cleanName == "." || cleanName == "/" {
			continue
		}

		targetPath := filepath.Join(absOutputDir, cleanName)
		relPath, err := filepath.Rel(absOutputDir, targetPath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			return fmt.Errorf("archive contains unsafe path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(targetPath, 0o755); err != nil {
				return fmt.Errorf("create directory %s: %w", targetPath, err)
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return fmt.Errorf("create parent directories for %s: %w", targetPath, err)
		}

		source, err := file.Open()
		if err != nil {
			return fmt.Errorf("open archive file %s: %w", file.Name, err)
		}

		destination, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			_ = source.Close()
			return fmt.Errorf("create extracted file %s: %w", targetPath, err)
		}

		_, copyErr := io.Copy(destination, source)
		closeDestErr := destination.Close()
		closeSourceErr := source.Close()

		if copyErr != nil {
			return fmt.Errorf("extract file %s: %w", targetPath, copyErr)
		}
		if closeDestErr != nil {
			return fmt.Errorf("close extracted file %s: %w", targetPath, closeDestErr)
		}
		if closeSourceErr != nil {
			return fmt.Errorf("close archive file %s: %w", file.Name, closeSourceErr)
		}
	}

	return nil
}
