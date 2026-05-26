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

	result, err := uploadMediaInternal(ctx, c.httpClient, url, filePath)
	if err != nil {
		return &UploadMediaResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &UploadMediaResult{Success: true}
	if data != nil {
		res.MediaID = strFromMap(data, "mediaId")
		res.CreatedTime = strFromMap(data, "createdTime")
	}
	return res, nil
}

func (c *LansengerClient) UploadAppMedia(ctx context.Context, filePath string, mediaType string, width, height, duration int) (*UploadAppMediaResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "app_medias", "create", token,
		WithMediaTypeString(mediaType),
		WithIntParam("width", width),
		WithIntParam("height", height),
		WithIntParam("duration", duration),
	)

	result, err := uploadMediaInternal(ctx, c.httpClient, url, filePath)
	if err != nil {
		return &UploadAppMediaResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &UploadAppMediaResult{Success: true}
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

func (c *LansengerClient) FetchMediaPath(ctx context.Context, mediaID string, userToken string) (*MediaPathResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "medias", "path_fetch", token,
		WithUserToken(userToken),
		WithPathVar("media_id", mediaID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &MediaPathResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &MediaPathResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &MediaPathResult{
		Success:     true,
		MediaPath:   strFromMap(data, "mediaPath"),
		Name:        strFromMap(data, "name"),
		Type:        strFromMap(data, "type"),
		Size:        strFromMap(data, "size"),
		RawResponse: result,
	}, nil
}

func uploadMediaInternal(ctx context.Context, httpClient *http.Client, url string, filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, NewFileError("cannot open file: " + err.Error())
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("media", filepath.Base(filePath))
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
