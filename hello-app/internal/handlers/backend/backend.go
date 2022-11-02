package product

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/controllers/backend"
	EntityBackend "bitbucket.org/codapayments/coda-stack-management-app/internal/entities/backend"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/handlers"
	"bitbucket.org/codapayments/coda-stack-management-app/internal/repository/adapter"
	Rules "bitbucket.org/codapayments/coda-stack-management-app/internal/rules"
	RulesBackend "bitbucket.org/codapayments/coda-stack-management-app/internal/rules/backend"
	HttpStatus "bitbucket.org/codapayments/coda-stack-management-app/utils/http"
	"net/http"
	"time"
)

type Handler struct {
	handlers.Interface

	Controller backend.Interface
	Rules      Rules.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Controller: backend.NewController(repository),
		Rules:      RulesBackend.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "ID") != "" {
		h.getOne(w, r)
	} else {
		h.getAll(w, r)
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	response, err := h.Controller.ListOne(ID)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.Controller.ListAll()
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	providerBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(providerBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, map[string]interface{}{"id": ID.String()})
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	providerBody, err := h.getBodyAndValidate(r, ID)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(ID, providerBody); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	if err := h.Controller.Remove(ID); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityBackend.Backend, error) {
	providerBody := &EntityBackend.Backend{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, providerBody)
	if err != nil {
		return &EntityBackend.Backend{}, errors.New("body is required")
	}

	productParsed, err := EntityBackend.InterfaceToModel(body)
	if err != nil {
		return &EntityBackend.Backend{}, errors.New("error on convert body to model")
	}

	setDefaultValues(productParsed, ID)

	return productParsed, h.Rules.Validate(productParsed)
}

func setDefaultValues(workspace *EntityBackend.Backend, ID uuid.UUID) {
	workspace.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		workspace.ID = uuid.New()
		workspace.CreatedAt = time.Now()
	} else {
		workspace.ID = ID
	}
}