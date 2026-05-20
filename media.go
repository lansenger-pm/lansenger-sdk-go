package lansenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (c *LansengerClient) UploadMedia(ctx context.Context, filePath string, mediaType int) (*UploadMediaResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "medias", "create", token,
		WithMediaType(mediaType),
	)

	result, err := uploadMediaInternal(ctx, c.httpClient, url, filePath, mediaType)
	if err != nil {
		return &UploadMediaResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &UploadMediaResult{Success: true}
	if data != nil {
		res.MediaID = strFromMap(data, "mediaId")
	}
	return res, nil
}

func (c *LansengerClient) DownloadMedia(ctx context.Context, mediaID string) (*DownloadMediaResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "medias", "fetch", token,
		WithPathVar("media_id", mediaID),
	)

	data, err := c.doGetRaw(ctx, url)
	if err != nil {
		return &DownloadMediaResult{Success: false, Error: err.Error()}, nil
	}

	return &DownloadMediaResult{
		Success: true,
		Data:    data,
	}, nil
}

func (c *LansengerClient) DownloadMediaToFile(ctx context.Context, mediaID string, targetPath string) (string, error) {
	result, err := c.DownloadMedia(ctx, mediaID)
	if err != nil {
		return "", err
	}
	if !result.Success {
		return "", fmt.Errorf("download failed: %s", result.Error)
	}

	if targetPath == "" {
		targetPath = mediaID
	}

	if err := os.WriteFile(targetPath, result.Data, 0644); err != nil {
		return "", fmt.Errorf("writing file: %w", err)
	}

	return targetPath, nil
}

func uploadMediaInternal(ctx context.Context, httpClient *http.Client, url string, filePath string, mediaType int) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, NewFileError("cannot open file: " + err.Error())
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("creating multipart form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copying file to multipart: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("closing multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("creating upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, NewNetworkError("upload request failed: " + err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading upload response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decoding upload response: %w", err)
	}

	errCode, _ := result["errCode"].(float64)
	if errCode != 0 {
		errMsg, _ := result["errMsg"].(string)
		return result, NewAPIError(errMsg, int(errCode))
	}

	return result, nil
}
