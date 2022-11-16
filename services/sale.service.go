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
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var currentTime = time.Now()

type SaleService struct {
	topicService *TopicService
	saleRepo     *repositories.SaleOpportunityRepository
	leadRepo     *repositories.LeadRepository
	logRepo      *repositories.LogRepository
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

func NewSaleService(topicService *TopicService, repository *repositories.Repository) *SaleService {
	return &SaleService{
		topicService: topicService,
		leadRepo:     repository.LeadRepo,
		saleRepo:     repository.SaleRepo,
		logRepo:      repository.LogRepo,
	}
}

func (s *SaleService) ExecuteMessage(messages types.RequestMessageOrder, source string) bool {
	order := messages.Order
	metadata := messages.Metadata
	images := messages.Images
	customerName := order.CustomerName
	phone := order.Phone
	email := order.Email
	assetType := order.AssetType
	customerId := order.CustomerId

	lead := s.findOrCreateLead(PayloadFindOrCreateLead{
		FullName:   customerName,
		Email:      email,
		Phone:      phone,
		Source:     source,
		Metadata:   metadata,
		CustomerId: customerId,
	})

	if lead != nil {
		group := s.getSaleGroup(phone, lead)
		days := order.Days
		code := order.Code
		code = s.saleRepo.GenerateCode(code)
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

		if isExistsHash(hash, s.saleRepo) {
			return true
		}
		entity.Code = s.saleRepo.GenerateCode(code)
		entity.DisbursedAmount = 0
		entity.DisbursedAt = nil
		entity.CreatedBy = createdBy
		entity.UpdatedBy = createdBy
		entity.Hash = hash
		entities.CreatingEntity(&entity.BaseEntity)

		saleOpp, err := s.saleRepo.BaseRepo.Create(entity)
		if err != nil {
			return false
		}
		s.afterSaleOppCreated(saleOpp)

		// Notification To Customer
		s.notification(lead.CustomerId, saleOpp)
	}

	return false
}

func (s *SaleService) disbursed(payload types.PayloadMessageDisbursed) {
	s.borrowDisbursed(types.PayloadBorrowDisbursed{
		LoanAmount:     0,
		ModifiedAmount: 0,
		ContractCode:   "",
	})

}

func (s *SaleService) borrowDisbursed(payload types.PayloadBorrowDisbursed) bool {
	contractCode, loanAmount, modifiedAmount := payload.ContractCode, payload.LoanAmount, payload.ModifiedAmount
	findOneOptions := options.FindOne()
	findOneOptions.SetSort(bson.D{{"createdAt", -1}})
	sale, _ := s.saleRepo.BaseRepo.FindOne(bson.M{"contractCode": contractCode}, &options.FindOneOptions{
		Sort: bson.M{"disbursedAt": -1},
	})
	if sale != nil {
		currentTime := time.Now()
		currentMonth := currentTime.Format(types.YYMM)
		disbursedMonth := sale.DisbursedAt.Format(types.YYMM)
		if currentMonth != disbursedMonth {
		_:
			s.saleRepo.BaseRepo.UpdateByID(sale.ID, bson.M{"disbursedAmount": loanAmount})
			afterSaleOppUpdated(sale)
			return true
		}
		newSale := sale
		newSale.DisbursedAmount = modifiedAmount
		newSale.Code = s.saleRepo.GenerateCode("")
		newSale.DisbursedAt = &currentTime
		entities.CreatingEntity(&newSale.BaseEntity)
		saleOpp, _ := s.saleRepo.BaseRepo.Create(newSale)
		s.afterSaleOppCreated(saleOpp)
	}

	return false
}

func (s *SaleService) createSaleOppDisbursed(ctx context.Context, payload types.PayloadMessageDisbursed) {
	saleOppCode, lead, saleOpp := payload.SaleOppCode, payload.Lead, payload.SaleOpp
	description := saleOpp.Description
	contractCode := saleOpp.ContractCode
	assetType := saleOpp.AssetType
	loanTerm := saleOpp.LoanTerm
	disbursedAmount := saleOpp.DisbursedAmount
	accountStore := saleOpp.AccountStore
	createdId := saleOpp.CreatedId
	demandLoan := saleOpp.DemandLoan
	disbursedAt := saleOpp.DisbursedAt

	fullName, phone, nationalId, account, customerId, source :=
		lead.FullName, lead.Phone, lead.NationalId, lead.Account, lead.CustomerId, lead.Source

	disbursedLead := s.findOrCreateLead(PayloadFindOrCreateLead{
		FullName:   fullName,
		Phone:      phone,
		Source:     source,
		CreatedBy:  account,
		CustomerId: customerId,
		NationalId: nationalId,
	})

	if disbursedLead != nil {
		entity := &entities.SaleOpportunity{
			SourceRefs: entities.SourceRefs{
				Source:     source,
				RefId:      contractCode,
				CustomerId: customerId,
			},
			Status: types.SaleOppStatusSuccess,
			Type:   types.SaleOppTypeBorrower,
			Source: source,
			Group:  s.getSaleGroup(phone, disbursedLead),
			Assets: entities.Asset{
				Description: description,
				AssetType:   assetType,
				DemandLoan:  demandLoan,
				LoanTerm:    loanTerm,
			},
			LeadId:          disbursedLead.ID,
			ContractCode:    contractCode,
			DisbursedAmount: disbursedAmount,
			DisbursedAt:     &disbursedAt,
			EmployeeBy:      createdId,
			BaseEntity: entities.BaseEntity{
				CreatedBy: createdId,
				UpdatedBy: createdId,
				CreatedAt: disbursedAt,
				UpdatedAt: disbursedAt,
			},
			StoreCode: accountStore,
			Code:      s.saleRepo.GenerateCode(saleOppCode),
		}

		saleOppDisbursed, _ := s.saleRepo.BaseRepo.Create(entity)
		s.afterSaleOppCreated(saleOppDisbursed)

		s.pushEventInternal(saleOppDisbursed)
	}
}

func (s *SaleService) findOrCreateLead(payload PayloadFindOrCreateLead) *entities.Lead {
	fullName := payload.FullName
	phone := payload.Phone
	email := payload.Email
	source := payload.Source
	createdBy := payload.CreatedBy
	nationalId := payload.NationalId
	customerId := payload.CustomerId

	item, err := s.leadRepo.BaseRepo.FindOne(bson.D{{"phone", phone}}, nil)
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
		item, _ = s.leadRepo.BaseRepo.Create(entity)
	}
	return item
}

func (s *SaleService) getSaleGroup(phone string, lead *entities.Lead) string {
	group := types.GroupNew

	if lead == nil {
		lead, _ = s.leadRepo.BaseRepo.FindOne(bson.D{{"phone", phone}}, nil)
	}
	if lead != nil {
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
		sale, err := s.saleRepo.BaseRepo.FindOne(filter, nil)
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

func (s *SaleService) notification(customerId string, sale *entities.SaleOpportunity) {
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

		s.topicService.Send(config.TopicConfig["customerOrderUpdated"], map[string]interface{}{
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

func afterSaleOppUpdated(sale *entities.SaleOpportunity) {

}

func (s *SaleService) afterSaleOppCreated(sale *entities.SaleOpportunity) {
_:
	s.logRepo.BaseRepo.Create(&entities.Log{
		BeforeAttributes:    utils.Omit(sale, []string{"leadId", "sourceRefs", "code", "createdAt", "updatedAt", "ID", "createdBy", "hash"}),
		AfterAttributes:     nil,
		SaleOpportunitiesId: sale.ID,
		CreatedBy:           sale.CreatedBy,
		CreatedAt:           time.Time{},
	})
}

func (s *SaleService) pushEventInternal(saleOpp *entities.SaleOpportunity) {
	s.topicService.Send(config.TopicConfig["customerOrderUpdated"], map[string]interface{}{
		"data":      "",
		"receivers": []string{"customerId"},
	}, map[string]string{
		"subscriptionType": "subscriptionType",
	})
}
