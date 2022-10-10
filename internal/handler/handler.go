package handler

import (
	"github.com/itsabgr/omp/internal/db"
	"github.com/itsabgr/omp/internal/model"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

type Handler struct {
	router *httprouter.Router
	db     *db.DB
}

func New(db *db.DB) *Handler {
	h := &Handler{
		router: httprouter.New(),
		db:     db,
	}
	h.init()
	return h
}

func (h *Handler) init() {
	h.router.POST("/image", h.createImage)
	h.router.POST("/image/:sha256/chunks", h.uploadChunk)
	h.router.GET("/image/:sha256", h.downloadImage)
}

func (h *Handler) createImage(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	if request.Header.Get("Content-Type") != "application/json" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusUnsupportedMediaType)
		response.Write(model.MarshalError("unsupported media type", http.StatusUnsupportedMediaType))
		return
	}
	body, err := io.ReadAll(io.LimitReader(request.Body, 10000))
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(model.MarshalError(err.Error(), http.StatusBadRequest))
		return
	}
	image := new(model.Image)
	if err := image.Unmarshal(body); err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(model.MarshalError(err.Error(), http.StatusBadRequest))
		return
	}
	if err := h.db.CreateImage(image); err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusConflict)
		response.Write(model.MarshalError(err.Error(), http.StatusConflict))
		return
	}
	response.WriteHeader(http.StatusCreated)
}

func (h *Handler) uploadChunk(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if request.Header.Get("Content-Type") != "application/json" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusUnsupportedMediaType)
		response.Write(model.MarshalError("unsupported media type", http.StatusUnsupportedMediaType))
		return
	}
	body, err := io.ReadAll(io.LimitReader(request.Body, 10000))
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(model.MarshalError(err.Error(), http.StatusBadRequest))
		return
	}
	chunk := new(model.Chunk)
	chunk.Image = params.ByName("sha256")
	if err := chunk.Unmarshal(body); err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(model.MarshalError(err.Error(), http.StatusBadRequest))
		return
	}
	switch err := h.db.PutChunk(chunk); err {
	case nil:
		break
	case db.ErrNotFound:
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusNotFound)
		response.Write(model.MarshalError(err.Error(), http.StatusNotFound))
		return
	case db.ErrExists:
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusConflict)
		response.Write(model.MarshalError(err.Error(), http.StatusConflict))
		return
	default:
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(model.MarshalError(err.Error(), http.StatusBadRequest))
		return
	}
	response.WriteHeader(http.StatusCreated)
}

func (h *Handler) downloadImage(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	imageSha256 := params.ByName("sha256")
	response.Header().Set("Content-Type", "text/plain")
	switch err := h.db.WriteImageTo(imageSha256, response); err {
	case nil:
		break
	case db.ErrNotFound:
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusNotFound)
		response.Write(model.MarshalError(err.Error(), http.StatusNotFound))
		return
	default:
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(model.MarshalError(err.Error(), http.StatusBadRequest))
		return
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
