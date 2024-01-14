package handlers

import (
	"bp-task/internal/database"
	"bp-task/internal/utils"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

const (
	errMsgGetToken         = "Error fetching token by CID"
	errMsgGetAllTokens     = "Error fetching all tokens"
	errMsgEncodingMetadata = "Error encoding token"
	errMsgParsingCSV       = "Error parsing CSV"
	errMsgInsertMetadata   = "Error inserting data into database"
)

// IPFSHandler handles requests related to IPFS addresses.
func IPFSHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Handle POST request (processing IPFS addresses and inserting into the database)
			handleIPFSPost(w, r, db)
		case http.MethodGet:
			// Handle GET requests for /tokens and /token/{CID}
			handleIPFSGet(w, r, db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleIPFSGet handles the GET requests for /tokens and /token/{CID}
func handleIPFSGet(w http.ResponseWriter, r *http.Request, db *database.DB) {
	parts := strings.Split(r.URL.Path, "/")
	switch {
	case len(parts) > 0 && parts[len(parts)-1] == "tokens":
		// Handle GET request for /tokens
		handleGetAllTokens(w, r.Context(), db)
	case len(parts) > 0 && parts[len(parts)-2] == "tokens" && parts[len(parts)-1] != "":
		// Handle GET request for /token/{CID}
		handleGetToken(w, r.Context(), db, parts[len(parts)-1])
	default:
		http.NotFound(w, r)
	}
}

// handleGetAllTokens handles the GET request for /tokens
func handleGetAllTokens(w http.ResponseWriter, ctx context.Context, db *database.DB) {
	tokens, err := db.GetTokens()
	if err != nil {
		slog.ErrorContext(ctx, errMsgGetAllTokens, "err", err)
		http.Error(w, errMsgGetAllTokens, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		slog.ErrorContext(ctx, errMsgEncodingMetadata, "err", err)
		http.Error(w, errMsgEncodingMetadata, http.StatusInternalServerError)
		return
	}
}

// handleGetToken handles the GET request for /token/{CID}
func handleGetToken(w http.ResponseWriter, ctx context.Context, db *database.DB, cid string) {
	token, err := db.GetToken(cid)
	if err != nil {
		slog.ErrorContext(ctx, errMsgGetToken, "err", err, "cid", cid)
		http.Error(w, errMsgGetToken, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		slog.ErrorContext(ctx, errMsgEncodingMetadata, "err", err, "cid", cid)
		http.Error(w, errMsgEncodingMetadata, http.StatusInternalServerError)
		return
	}
}

// handleIPFSPost handles requests related to IPFS addresses.
func handleIPFSPost(w http.ResponseWriter, r *http.Request, db *database.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 0 && parts[len(parts)-1] != "ipfs" {
		http.NotFound(w, r)
		return
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	data, err := utils.ParseCSVFileFromRequest(r)
	if err != nil {
		slog.ErrorContext(r.Context(), errMsgParsingCSV, "err", err)
		http.Error(w, errMsgParsingCSV, http.StatusBadRequest)
		return
	}

	// Process the IPFS addresses (parser task) asynchronously in a goroutine with a throttler
	go func(data [][]string) {
		throttler := make(chan struct{}, 5) // Limit to 5 concurrent goroutines, this could be a parameter or env variable
		for _, record := range data {
			if len(record) > 0 {
				throttler <- struct{}{}
				go func(ctx context.Context, cid string) {
					defer func() { <-throttler }()
					slog.InfoContext(ctx, "Received IPFS cid ", "cid", cid)
					// get data from ipfs address (cid)
					data, err := utils.GetMetadataFromIPFS(ctx, cid)
					if err != nil {
						slog.ErrorContext(ctx, "Error getting data from IPFS", "err", err, "cid", cid, "data", data)
						return
					}
					// Insert metadata into the database
					err = db.InsertMetadata(cid, data)
					if err != nil {
						slog.ErrorContext(ctx, errMsgInsertMetadata, "err", err, "cid", cid, "data", data)
					}
				}(r.Context(), record[0])
			}
		}
	}(data)

	// Respond with a JSON message
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(NewResponse("IPFS metadata processed successfully and inserted into the database"))

	if err != nil {
		slog.ErrorContext(r.Context(), "Error encoding JSON response: %v", err)
	}
}
