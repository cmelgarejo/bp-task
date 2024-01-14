package utils

import (
	"encoding/csv"
	"io"
	"log/slog"
	"net/http"
)

// ParseCSVFileFromRequest reads a CSV file from a multipart/form-data request.
func ParseCSVFileFromRequest(r *http.Request) ([][]string, error) {
	err := r.ParseMultipartForm(32 << 10) // 32 kb limit for the entire list, so no abuse is allowed if I host this publicly
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("ipfs-list") // field name "ipfs-list", so I avoid people trying with "file" or "csv"
	slog.Info("file", "file", file)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var data [][]string

	for {
		cid, err := reader.Read()
		slog.Info("data:", "cid", cid)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(cid) > 0 && cid[0] != "" {
			data = append(data, cid)
		}
	}

	return data, nil
}
