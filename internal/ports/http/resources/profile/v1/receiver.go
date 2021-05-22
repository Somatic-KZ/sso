package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/JetBrainer/sso/internal/ports/http/resources/profile"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReceiversRequest struct {
	Receivers []*models.ReceiverResponse `json:"receivers"`
}

// Receivers @Summary Получатели
// @Description Позволяет получить информацию(кроме адреса) обо всех получателях привязанных к пользователю
// @Produce json
// @Tags profile
// @Security JWT
// @Success 200 {object} ReceiversRequest
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /profile/receivers [get]
func (p *ProfileResource) Receivers(w http.ResponseWriter, r *http.Request) {
	var tdid interface{}
	if tdid = r.Context().Value("tdid"); tdid == "" {
		_ = render.Render(w, r, resources.BadRequest(profile.ErrUnknownTDID))
		return
	}

	id, err := primitive.ObjectIDFromHex(tdid.(string))
	if err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	users := p.authManager.Users()

	var user *models.User

	user, err = users.ByTDID(id)
	if err != nil {
		_ = render.Render(w, r, resources.ResourceNotFound(err))
		return
	}

	m := make(map[string]*models.ReceiverResponse, len(user.Receivers))

	for _, receiver := range user.Receivers {
		_, ok := m[receiver.PrimaryPhone]
		if ok {
			continue
		}
		m[receiver.PrimaryPhone] = &models.ReceiverResponse{
			ID:              receiver.ID,
			FirstName:       receiver.FirstName,
			LastName:        receiver.LastName,
			Email:           receiver.Email,
			PrimaryPhone:    receiver.PrimaryPhone,
			AdditionalPhone: receiver.AdditionalPhone,
			IsDefault:       receiver.IsDefault,
			IsOrganization:  receiver.IsOrganization,
			Organization:    receiver.Organization,
		}

	}

	receivers := make([]*models.ReceiverResponse, 0, len(user.Receivers))
	for _, value := range m {
		receivers = append(receivers, value)
	}
	if len(receivers) == 0 {
		_ = render.Render(w, r, resources.ResourceNotFound(errors.New("not found")))
		return
	}

	render.JSON(w, r, &ReceiversRequest{Receivers: receivers})
}

type AddressesResponse struct {
	Addresses []*models.ReceiverAddressResponse `json:"addresses"`
}

// ReceiversAddress @Summary Получатели
// @Description Позволяет получить информацию об адресах всех получателей привязанных к пользователю
// @Produce json
// @Tags profile
// @Security JWT
// @Success 200 {object} AddressesResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /profile/addresses [get]
func (p *ProfileResource) ReceiversAddress(w http.ResponseWriter, r *http.Request) {
	var tdid interface{}
	if tdid = r.Context().Value("tdid"); tdid == "" {
		_ = render.Render(w, r, resources.BadRequest(profile.ErrUnknownTDID))
		return
	}

	id, err := primitive.ObjectIDFromHex(tdid.(string))
	if err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	users := p.authManager.Users()

	var user *models.User

	user, err = users.ByTDID(id)
	if err != nil {
		_ = render.Render(w, r, resources.ResourceNotFound(err))
		return
	}

	var sb strings.Builder

	m := make(map[string]*models.ReceiverAddressResponse)

	for _, receiver := range user.Receivers {
		if receiver.Address == nil {
			continue
		}
		sb.WriteString(receiver.Address.City)
		sb.WriteString(receiver.Address.Street)
		sb.WriteString(receiver.Address.House)
		_, ok := m[sb.String()]
		if ok {
			sb.Reset()
			continue
		}
		m[sb.String()] = &models.ReceiverAddressResponse{
			ID:        receiver.ID,
			Region:    receiver.Address.Region,
			City:      receiver.Address.City,
			Street:    receiver.Address.Street,
			House:     receiver.Address.House,
			Floor:     receiver.Address.Floor,
			Apartment: receiver.Address.Apartment,
			Zipcode:   receiver.Address.Zipcode,
			Geo:       receiver.Address.Geo,
		}
		sb.Reset()
	}

	receiversAddress := make([]*models.ReceiverAddressResponse, 0, len(user.Receivers))
	for _, value := range m {
		receiversAddress = append(receiversAddress, value)
	}
	if len(receiversAddress) == 0 {
		_ = render.Render(w, r, resources.ResourceNotFound(errors.New("not found")))
		return
	}

	render.JSON(w, r, &AddressesResponse{Addresses: receiversAddress})

}
