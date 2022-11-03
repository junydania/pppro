package product

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/junydania/pppro/hello-app/internal/controllers/hello"
	EntityHello "github.com/junydania/pppro/hello-app/internal/entities/hello"
	"github.com/junydania/pppro/hello-app/internal/handlers"
	"github.com/junydania/pppro/hello-app/internal/repository/adapter"
	Rules "github.com/junydania/pppro/hello-app/internal/rules"
	RulesHello "github.com/junydania/pppro/hello-app/internal/rules/hello"
	HttpStatus "github.com/junydania/pppro/hello-app/utils/http"
	"net/http"
	"time"
)

type Handler struct {
	handlers.Interface
	Controller hello.Interface
	Rules      Rules.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Controller: hello.NewController(repository),
		Rules:      RulesHello.NewRules(),
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
	helloBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(helloBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, map[string]interface{}{"id": ID.String()})
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityHello.Hello, error) {
	helloBody := &EntityHello.Hello{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, helloBody)
	if err != nil {
		return &EntityHello.Hello{}, errors.New("body is required")
	}

	helloParsed, err := EntityHello.InterfaceToModel(body)
	if err != nil {
		return &EntityHello.Hello{}, errors.New("error on convert body to model")
	}

	setDefaultValues(helloParsed, ID)

	return helloParsed, err
}

func setDefaultValues(hello *EntityHello.Hello, ID uuid.UUID) {
	hello.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		hello.ID = uuid.New()
		hello.CreatedAt = time.Now()
	} else {
		hello.ID = ID
	}
}