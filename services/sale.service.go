package services

import (
	"crm-worker-go/config"
	"crm-worker-go/entities"
	"crm-worker-go/repositories"
	"crm-worker-go/types"
	"crm-worker-go/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var currentTime = time.Now()
var createdBy = config.GetConfig().DefaultDataConfig.CreatedBy

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

func (s *SaleService) Disbursed(payload types.PayloadMessageDisbursed) bool {
	if payload.ContractCode != "" {
		return s.borrowDisbursed(types.PayloadBorrowDisbursed{
			ContractCode:   payload.ContractCode,
			LoanAmount:     payload.LoanAmount,
			ModifiedAmount: payload.ModifiedAmount,
		})
	}

	saleOppCode := payload.SaleOppCode

	if saleOppCode != "" {
		s.updateSaleOppDisbursed(payload)
		return true
	}

	s.createSaleOppDisbursed(payload)
	return true
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

		if currentMonth == disbursedMonth {
		_:
			s.saleRepo.BaseRepo.UpdateByID(sale.ID, bson.M{"disbursedAmount": loanAmount})
			s.afterSaleOppUpdated(sale)
			return true
		}
		newSale := sale
		newSale.ID = primitive.NewObjectID()
		newSale.DisbursedAmount = modifiedAmount
		newSale.Code = s.saleRepo.GenerateCode("")
		newSale.DisbursedAt = &currentTime
		entities.CreatingEntity(&newSale.BaseEntity)
		saleOpp, _ := s.saleRepo.BaseRepo.Create(newSale)
		s.afterSaleOppCreated(saleOpp)
		return true
	}

	return false
}

func (s *SaleService) createSaleOppDisbursed(payload types.PayloadMessageDisbursed) {
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

func (s *SaleService) updateSaleOppDisbursed(payload types.PayloadMessageDisbursed) bool {
	loanPackageCode, saleOppCode, saleOpp, lead :=
		payload.LoanPackageCode, payload.SaleOppCode, payload.SaleOpp, payload.Lead
	description,
		contractCode,
		assetType,
		disbursedAmount,
		accountStore,
		createdId := saleOpp.Description,
		saleOpp.ContractCode,
		saleOpp.AssetType,
		saleOpp.DisbursedAmount,
		saleOpp.AccountStore,
		saleOpp.CreatedId

	saleItem, _ := s.saleRepo.BaseRepo.FindOne(bson.M{"code": saleOppCode}, nil)

	if saleItem != nil {
		leadUpdated, err := s.updateLead(lead)
		if err != nil {
			utils.Logger.Debug(err)
			return false
		}

		entity := bson.M{
			"loanPackageCode": loanPackageCode,
			"status":          types.SaleOppStatusSuccess,
			"assets": entities.Asset{
				Description: description,
				Media:       []entities.AssetMedia{},
				AssetType:   assetType,
				DemandLoan:  saleItem.Assets.DemandLoan,
				LoanTerm:    saleItem.Assets.LoanTerm,
			},
			"employeeBy":      createdId,
			"storeCode":       accountStore,
			"disbursedAt":     currentTime,
			"contractCode":    contractCode,
			"disbursedAmount": disbursedAmount,
			"updatedBy":       createdId,
			"updatedAt":       currentTime,
		}
		if lead.Phone != leadUpdated.Phone {
			metadata := saleItem.Metadata
			metadata["phone"] = lead.Phone
			entity["metadata"] = metadata
			entity["group"] = s.getSaleGroup(leadUpdated.Phone, leadUpdated)
		}
		_, err = s.saleRepo.BaseRepo.UpdateByID(saleItem.ID, entity)
		utils.Logger.Debug(err)
		return false
	}

	return false
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

		s.topicService.Send(config.GetConfig().TopicConfig.CustomerOrderUpdated, map[string]interface{}{
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

func (s *SaleService) afterSaleOppUpdated(before *entities.SaleOpportunity) {
	after, _ := s.saleRepo.BaseRepo.FindById(before.ID)
	beforeData := utils.Omit(before, []string{
		"leadId",
		"source_refs",
		"code",
		"createdAt",
		"updatedAt",
		"id",
		"createdBy",
		"metadata",
		"deletedAt",
		"assets",
		"source_refs",
	})

	afterData := utils.Omit(after, []string{
		"leadId",
		"source_refs",
		"code",
		"createdAt",
		"updatedAt",
		"id",
		"createdBy",
		"metadata",
		"deletedAt",
		"assets",
		"source_refs",
	})

	id, createdBy := after.ID, after.UpdatedBy

	var keyChange []string
	for key, val := range afterData {
		if val != beforeData[key] {
			keyChange = append(keyChange, key)
		}
	}

	s.logRepo.BaseRepo.Create(&entities.Log{
		ID:                  primitive.ObjectID{},
		BeforeAttributes:    utils.Pick(beforeData, keyChange),
		AfterAttributes:     utils.Pick(afterData, keyChange),
		SaleOpportunitiesId: id,
		CreatedBy:           createdBy,
		CreatedAt:           time.Now(),
	})
}

func (s *SaleService) updateLead(payload types.MessageDisbursedLead) (*entities.Lead, error) {
	phone, fullName, nationalId, account, accountStore, customerId :=
		payload.Phone, payload.FullName, payload.NationalId, payload.Account, payload.AccountStore, payload.CustomerId

	lead, err := s.leadRepo.BaseRepo.FindOne(bson.M{"phone": phone}, nil)
	if lead != nil {
		item, err := s.leadRepo.BaseRepo.UpdateByID(lead.ID, bson.M{
			"fullName":   fullName,
			"nationalId": nationalId,
			"employeeBy": account,
			"customerId": customerId,
			"storeCode":  accountStore,
		})
		if err != nil {
			return nil, err
		}
		return item, nil
	}
	return nil, err
}

func (s *SaleService) afterSaleOppCreated(sale *entities.SaleOpportunity) {
_:
	s.logRepo.BaseRepo.Create(&entities.Log{
		BeforeAttributes:    utils.Omit(sale, []string{"leadId", "source_refs", "code", "createdAt", "updatedAt", "ID", "createdBy", "hash"}),
		AfterAttributes:     nil,
		SaleOpportunitiesId: sale.ID,
		CreatedBy:           sale.CreatedBy,
		CreatedAt:           time.Now(),
	})
}

func (s *SaleService) pushEventInternal(saleOpp *entities.SaleOpportunity) {
	s.topicService.Send(config.GetConfig().TopicConfig.CustomerOrderUpdated, map[string]interface{}{
		"data":      "",
		"receivers": []string{"customerId"},
	}, map[string]string{
		"subscriptionType": "subscriptionType",
	})
}
