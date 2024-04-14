package controller

import (
	"GOHW-1/internal/db"
	"GOHW-1/internal/model"
	"GOHW-1/internal/repository/postgresql"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	queryParamKey  = "key"
	configUsername = "ildus"
	configPassword = "erbaev"
	securePort     = ":9001"
	insecurePort   = ":9000"
)

//go:generate mockgen -package controller -destination=./mocks/mock_repository.go . PickUpPointsRepo
type PickUpPointsRepo interface {
	Create(ctx context.Context, pickUpPoint *model.PickUpPoint) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.PickUpPoint, error)
	List(ctx context.Context) ([]model.PickUpPoint, error)
	Update(ctx context.Context, id int64, updateData model.PickUpPoint) error
	Delete(ctx context.Context, id int64) error
}

type Sender interface {
	sendAsyncMessage(message LoggingMessage) error
}

type PickUpPointController struct {
	Repo   PickUpPointsRepo
	Sender Sender
}

func NewPickUpPointController(database *db.Database, sender *KafkaSender) *PickUpPointController {
	pickUpPointRepo := postgresql.NewPickUpPoints(*database)
	//sender := NewKafkaSender(producer, topic)

	return &PickUpPointController{
		Repo:   pickUpPointRepo,
		Sender: sender,
	}
}

func (controller *PickUpPointController) StartHTTPServer() {
	http.Handle("/", createRouter(*controller))
	go func() {
		if err := http.ListenAndServeTLS(securePort, "./server.crt", "./server.key", nil); err != nil {
			fmt.Println(fmt.Errorf("cannot handle: %w", err))
			return
		}
	}()

	redirectToHTTPS()
}

func redirectToHTTPS() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		secureURL := "https://" + r.Host + securePort + r.RequestURI
		http.Redirect(w, r, secureURL, http.StatusMovedPermanently)
	})

	if err := http.ListenAndServe(insecurePort, nil); err != nil {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func createRouter(controller PickUpPointController) *mux.Router {
	router := mux.NewRouter()
	router.Use(AuthMiddleware)
	router.Use(controller.LoggingMiddleware)
	router.HandleFunc("/pick-up-point", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			controller.Create(w, req)
		case http.MethodGet:
			controller.List(w, req)
		default:
			http.Error(w, "method is not implemented", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(fmt.Sprintf("/pick-up-point/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			controller.GetByID(w, req)
		case http.MethodDelete:
			controller.Delete(w, req)
		case http.MethodPut:
			controller.Update(w, req)
		default:
			http.Error(w, "method is not implemented", http.StatusMethodNotAllowed)
		}
	})
	return router
}

// Create handles the creation of a new pick-up point
func (controller *PickUpPointController) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error occured: %v", err), http.StatusBadRequest)
		return
	}
	var request model.PickUpPoint
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, fmt.Sprintf("error occured: %v", err), http.StatusUnprocessableEntity)
		return
	}

	pickUpPointJson, status, err := controller.CreatePickUpPoint(req.Context(), request)
	if err != nil {
		http.Error(w, fmt.Sprintf("error occured: %v", err), status)
		return
	}
	w.WriteHeader(status)
	w.Write(pickUpPointJson)
}

func (controller *PickUpPointController) CreatePickUpPoint(ctx context.Context, request model.PickUpPoint) ([]byte, int, error) {
	pickUpPointRepo := &model.PickUpPoint{
		Name:    request.Name,
		Address: request.Address,
		Contact: request.Contact,
	}

	id, err := controller.Repo.Create(ctx, pickUpPointRepo)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("can not created pick-up point")
	}

	resp := &model.PickUpPoint{
		ID:      id,
		Name:    pickUpPointRepo.Name,
		Address: pickUpPointRepo.Address,
		Contact: pickUpPointRepo.Contact,
	}
	pickUpPointJson, _ := json.Marshal(*resp)
	return pickUpPointJson, http.StatusOK, nil
}

// GetByID returns json of pick-up point by given id
func (controller *PickUpPointController) GetByID(w http.ResponseWriter, req *http.Request) {
	idStr, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	data, status, err := controller.GetJSONByID(req.Context(), idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("error occured: %v", err), status)
		return
	}
	w.WriteHeader(status)
	w.Write(data)
}

func (controller *PickUpPointController) GetJSONByID(ctx context.Context, idStr string) ([]byte, int, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("id must be a number")
	}

	pickUpPoint, err := controller.Repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.StatusNotFound, fmt.Errorf("pick-up point not found")
		}
		return nil, http.StatusNotFound, fmt.Errorf("error occured: %v", err)
	}
	pickUpPointJson, _ := json.Marshal(pickUpPoint)
	return pickUpPointJson, http.StatusOK, nil
}

// List returns json of all pick-up points
func (controller *PickUpPointController) List(w http.ResponseWriter, req *http.Request) {
	pickUpPoints, err := controller.Repo.List(req.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("error occured: %v", err), http.StatusInternalServerError)
		return
	}

	pickUpPointsJson, _ := json.Marshal(pickUpPoints)
	w.Write(pickUpPointsJson)
}

// Update modifies an existing pick-up point
func (controller *PickUpPointController) Update(w http.ResponseWriter, req *http.Request) {
	idStr, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error occured: %v", err), http.StatusBadRequest)
		return
	}

	var unm model.PickUpPoint
	if err = json.Unmarshal(body, &unm); err != nil {
		http.Error(w, fmt.Sprintf("error occured: %v", err), http.StatusUnprocessableEntity)
		return
	}

	pickUpPointJson, status, err := controller.UpdateByID(req.Context(), idStr, unm)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	w.WriteHeader(status)
	w.Write(pickUpPointJson)
}

func (controller *PickUpPointController) UpdateByID(ctx context.Context, idStr string, unm model.PickUpPoint) ([]byte, int, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("id must be a number")
	}

	pickUpPointRepo := &model.PickUpPoint{
		Name:    unm.Name,
		Address: unm.Address,
		Contact: unm.Contact,
	}

	if err := controller.Repo.Update(ctx, id, *pickUpPointRepo); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("pick-up point not found")
	}

	resp := &model.PickUpPoint{
		ID:      id,
		Name:    pickUpPointRepo.Name,
		Address: pickUpPointRepo.Address,
		Contact: pickUpPointRepo.Contact,
	}
	pickUpPointJson, _ := json.Marshal(resp)
	return pickUpPointJson, http.StatusOK, nil
}

// Delete removes a pick-up point by given ID
func (controller *PickUpPointController) Delete(w http.ResponseWriter, req *http.Request) {
	idStr, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	status, err := controller.DeleteByID(req.Context(), idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("error occured: %v", err), status)
		return
	}
	w.WriteHeader(status)

}

func (controller *PickUpPointController) DeleteByID(ctx context.Context, idStr string) (int, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("id must be a number")
	}

	if err = controller.Repo.Delete(ctx, id); err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
