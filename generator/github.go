package generator

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/daniel1302/vega-asistant/utils"
)

type ArtifactType string

const (
	ArtifactVega  ArtifactType = "vega"
	ArtifactVisor ArtifactType = "visor"
)

func DownloadArtifact(
	repository, version, outputDir string,
	artifactType ArtifactType,
) (string, error) {
	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH

	artifactName := fmt.Sprintf("%s-%s-%s.zip", artifactType, operatingSystem, architecture)

	artifactURL := fmt.Sprintf(
		"https://github.com/%s/releases/download/%s/%s",
		repository,
		version,
		artifactName,
	)

	filePath := filepath.Join(outputDir, artifactName)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local artifact file: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(artifactURL)
	if err != nil {
		return "", fmt.Errorf("failed to get file from '%s': %w", artifactURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad http status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf(
			"failed to copy bytes from downloaded file(%s) to the local destination(%s): %w",
			artifactURL,
			filePath,
			err,
		)
	}

	if err := utils.Unzip(filePath, outputDir); err != nil {
		return "", fmt.Errorf("failed to unzip downloaded artifact(%s): %w", filePath, err)
	}

	binaryPath := filepath.Join(outputDir, string(artifactType))
	if err := os.Chmod(binaryPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to change permissions mod for binary %s: %w", binaryPath, err)
	}

	return binaryPath, nil
}
