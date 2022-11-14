package services

import (
	"context"
	"crm-worker-go/config"
	"crm-worker-go/entities"
	"crm-worker-go/repositories"
	"crm-worker-go/types"
	"crm-worker-go/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var currentTime = time.Now()

type SaleService struct {
}

type PayloadFindOrCreateLead struct {
	FullName   string
	Email      string
	Phone      string
	Source     string
	NationalId string
	Metadata   interface{}
	CustomerId string
	CreatedBy  string
}

func NewSaleService() *SaleService {
	return &SaleService{}
}

func (s *SaleService) ExecuteMessage(messages types.RequestMessageOrder, source string) bool {

	ctx := context.Background()

	order := messages.Order
	metadata := messages.Metadata
	images := messages.Images
	customerName := order.CustomerName
	phone := order.Phone
	email := order.Email
	assetType := order.AssetType
	customerId := order.CustomerId

	lead := findOrCreateLead(ctx, PayloadFindOrCreateLead{
		FullName:   customerName,
		Email:      email,
		Phone:      phone,
		Source:     source,
		Metadata:   metadata,
		CustomerId: customerId,
	})
	utils.Logger.Info(lead)
	if lead != nil {
		saleRepo := repositories.NewSaleOpportunityRepository(ctx)

		group := getSaleGroup(ctx, phone, lead)
		days := order.Days
		code := order.Code
		code = saleRepo.GenerateCode(code)
		detail := order.Detail
		bill := order.Bill
		createdBy := config.DefaultDataConfig["createdBy"]
		id := order.Id

		assets := entities.Asset{
			Description: detail,
			AssetType:   assetType,
			DemandLoan:  bill,
			LoanTerm:    days,
			Media:       getMedia(images),
		}
		entity := &entities.SaleOpportunity{
			SourceRefs: entities.SourceRefs{
				Source:     source,
				RefId:      id,
				CustomerId: order.CreatedBy,
			},
			Status: types.SaleOppStatusNew,
			Type:   types.SaleOppTypeBorrower,
			Source: source,
			Group:  group,
			Assets: assets,
			LeadId: lead.ID,
		}

		hash := utils.Hash(assets)

		if isExistsHash(hash, saleRepo) {
			return true
		}
		entity.Code = saleRepo.GenerateCode(code)
		entity.DisbursedAmount = 0
		entity.CreatedBy = createdBy
		entity.UpdatedBy = createdBy
		entity.Hash = hash
		entities.CreatingEntity(&entity.BaseEntity)

		saleOpp, err := saleRepo.BaseRepo.Create(entity)
		if err != nil {
			return false
		}

		// Notification To Customer
		notification(lead.CustomerId, saleOpp)

		utils.Logger.Info(saleOpp)
	}

	return false
}

func findOrCreateLead(ctx context.Context, payload PayloadFindOrCreateLead) *entities.Lead {
	fullName := payload.FullName
	phone := payload.Phone
	email := payload.Email
	source := payload.Source
	createdBy := payload.CreatedBy
	nationalId := payload.NationalId
	customerId := payload.CustomerId

	leadRepo := repositories.NewLeadRepository(ctx)

	item, err := leadRepo.BaseRepo.FindOne(bson.D{{"phone", phone}})
	if err != nil || item == nil {
		entity := &entities.Lead{
			FullName:   fullName,
			Phone:      phone,
			Email:      email,
			NationalId: nationalId,
			Source:     source,
			CustomerId: customerId,
			BaseEntity: entities.BaseEntity{
				CreatedBy: createdBy,
			},
		}
		entities.CreatingEntity(&entity.BaseEntity)
		item, _ = leadRepo.BaseRepo.Create(entity)
	}
	return item
}

func getSaleGroup(ctx context.Context, phone string, lead *entities.Lead) string {
	group := types.GroupNew

	if lead == nil {
		leadRepo := repositories.NewLeadRepository(ctx)
		lead, _ = leadRepo.BaseRepo.FindOne(bson.D{{"phone", phone}})
	}
	if lead != nil {
		saleRepo := repositories.NewSaleOpportunityRepository(ctx)
		filter := bson.M{
			"leadId": lead.ID,
			"disbursedAt": bson.M{
				"$gte": currentTime.AddDate(0, 0, -90),
			},
			"contractCode": bson.M{
				"$ne": nil,
			},
			"disbursedAmount": bson.M{
				"$ne": 0,
			},
		}
		sale, err := saleRepo.BaseRepo.FindOne(filter)
		if sale != nil && err == nil {
			group = types.GroupOld
		}
	}

	return group
}

func getMedia(images []interface{}) []entities.AssetMedia {
	medias := make([]entities.AssetMedia, 0)

	for i := 0; i < len(images); i++ {
		medias = append(medias, entities.AssetMedia{
			Url:      "",
			MimeType: "",
		})
	}
	return medias
}

func isExistsHash(hash string, saleRepo *repositories.SaleOpportunityRepository) bool {
	total, _ := saleRepo.BaseRepo.Count(bson.M{"hash": hash})
	return total != 0
}

func notification(customerId string, sale *entities.SaleOpportunity) {
	utils.Logger.Info(customerId)
	if customerId != "" {
		code := sale.Code
		status := sale.Status

		dataNotify := map[string]interface{}{
			"code":  code,
			"order": sale,
		}

		if status != types.SaleOppStatusNew {
			dataNotify["content"] = fmt.Sprintf("đã được cập nhật thành %v", getOrderStatusDisplay(status))
		}
		var subscriptionType string
		if status == types.SaleOppStatusNew {
			subscriptionType = types.TopicSubscriptionTypeOrderCreated
		} else {
			subscriptionType = types.TopicSubscriptionTypeOrderUpdated
		}

		topicService := NewTopicService()
		topicService.Send(config.TopicConfig["customerOrderUpdated"], map[string]interface{}{
			"data":      dataNotify,
			"receivers": []string{customerId},
		}, map[string]string{
			"subscriptionType": subscriptionType,
		})
	}
}

func getOrderStatusDisplay(status string) string {
	switch status {
	case types.SaleOppStatusNew:
		return "Moi"
	case types.SaleOppStatusPending:
		return "Đã xác nhận"
	case types.SaleOppStatusConsulting:
	case types.SaleOppStatusDealt:
	case types.SaleOppStatusUnContactable:
		return "Đang xử lý"
	case types.SaleOppStatusSuccess:
		return "Thành công"
	case types.SaleOppStatusCancel:
	case types.SaleOppStatusDenied:
		return "Đã huỷ"
	default:
		return ""
	}
	return ""
}
