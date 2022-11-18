package services

import (
	"context"
	"crm-worker-go/clients"
	"crm-worker-go/entities"
	"crm-worker-go/repositories"
	"crm-worker-go/types"
	"crm-worker-go/utils"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"strings"
	"sync"
	"time"
)

var ctx = context.Background()

type ExportService struct {
	fileManagerClient *clients.FileManagerClient
	employeeClient    *clients.EmployeeClient
	masterDataClient  *clients.MasterDataClient
	saleRepo          *repositories.SaleOpportunityRepository
	leadRepo          *repositories.LeadRepository
	tagRepo           *repositories.TagRepository
	noteRepo          *repositories.NoteRepository
	uploader          StorageService
}

func NewExportService(client *clients.HttpClient, repository *repositories.Repository, uploader StorageService) *ExportService {
	return &ExportService{
		fileManagerClient: client.FileManagerClient,
		employeeClient:    client.EmployeeClient,
		masterDataClient:  client.MasterDataClient,
		saleRepo:          repository.SaleRepo,
		leadRepo:          repository.LeadRepo,
		tagRepo:           repository.TagRepo,
		noteRepo:          repository.NoteRepo,
		uploader:          uploader,
	}
}

type getDataDB struct {
	leads map[string]*entities.Lead
	notes map[string]*entities.Note
	tags  map[string]*entities.Tag
}

type iMasterData struct {
	assetTypes *map[string]string
	sources    *map[string]string
	statuses   *map[string]string
	types      *map[string]string
	provinces  *map[string]string
	groups     *map[string]string
}

type iJobs struct {
	sale       *entities.SaleOpportunity
	items      getDataDB
	masterData iMasterData
	employees  *map[string]string
}

func (s *ExportService) ExportSaleOpp(payload types.PayloadMessageExport) bool {
	s.fileManagerClient.Client.SetToken(payload.Token)
	s.masterDataClient.Client.SetToken(payload.Token)
	s.employeeClient.Client.SetToken(payload.Token)

	var wg = sync.WaitGroup{}
	wg.Add(3)

	exportRequest, err := s.fileManagerClient.FindExportRequests(ctx, payload.ID)

	if err != nil {
		utils.Logger.Error("Get Export Request Failure ", err, "With Id: ", payload.ID)
		return false
	}

	if utils.Contains([]string{types.Completed, types.Failure}, exportRequest.Status) {
		return true
	}

	requestData := exportRequest.Data
	start, _ := time.Parse(types.DDMMYYYY, requestData.Start)
	end, _ := time.Parse(types.DDMMYYYY, requestData.End)

	filter := bson.D{
		{
			"deletedAt", nil,
		},
		{
			"createdAt", bson.M{
				"$gte": primitive.NewDateTimeFromTime(start),
				"$lt":  primitive.NewDateTimeFromTime(end),
			},
		},
	}

	if requestData.StatusExport != "" {
		statuses := strings.Split(requestData.StatusExport, ",")
		filter = append(filter, bson.E{
			Key: "status", Value: bson.M{
				"$in": statuses,
			},
		})
	}

	if requestData.SourceExport != "" {
		sources := strings.Split(requestData.SourceExport, ",")
		filter = append(filter, bson.E{
			Key: "source", Value: bson.M{
				"$in": sources,
			},
		})
	}

	if requestData.StoreCodes != "" {
		storeCodes := strings.Split(requestData.StoreCodes, ",")
		filter = append(filter, bson.E{
			Key: "storeCode", Value: bson.M{
				"$in": storeCodes,
			},
		})
	}

	results, _ := s.saleRepo.BaseRepo.Find(filter, nil)
	if len(results) == 0 {
		return true
	}

	var leadIds []primitive.ObjectID
	var employeeIds []string
	var saleIds []primitive.ObjectID
	for _, sale := range results {
		leadIds = append(leadIds, sale.LeadId)
		saleIds = append(saleIds, sale.ID)
		employeeIds = append(employeeIds, sale.CreatedBy, sale.UpdatedBy, sale.EmployeeBy)
	}

	var items getDataDB
	go func() {
		items = s.getDataDatabase(ctx, saleIds, leadIds)
		wg.Done()
	}()

	var employees *map[string]string
	go func(ctx context.Context) {
		employees, _ = s.employeeClient.GetEmployees(ctx, employeeIds)
		wg.Done()
	}(ctx)

	var masterData iMasterData
	go func(ctx context.Context) {
		masterData = s.getMasterData(ctx, payload.Token)
		wg.Done()
	}(ctx)
	wg.Wait()

	exportData := [][]interface{}{
		{
			"Updated At",
			"Created At",
			"Code",
			"Lead",
			"Phone",
			"Email",
			"Type",
			"Source",
			"Group",
			"Asset Type",
			"Description",
			"DemandLoan",
			"Account",
			"Account Store",
			"Address",
			"Province",
			"Latest Note",
			"Status",
			"Reason",
			"DisbursedAmount",
			"Contract Code",
			"Created By",
			"Updated By",
		},
	}

	number, numberWorker := len(results), 4
	jobs, result := make(chan iJobs, number), make(chan []interface{}, number)

	wg.Add(4)
	for i := 1; i <= numberWorker; i++ {
		go worker(jobs, result)
		wg.Done()
	}
	wg.Wait()

	for _, sale := range results {
		jobs <- iJobs{
			sale:       sale,
			items:      items,
			masterData: masterData,
			employees:  employees,
		}
	}
	close(jobs)

	for _, _ = range results {
		exportData = append(exportData, <-result)
	}

	_ = s.saveFile(exportData, payload.ID)
	return true
}

func worker(jobs <-chan iJobs, results chan<- []interface{}) {
	for n := range jobs {
		results <- handleData(n)
	}
}

func handleData(
	payload iJobs,
) []interface{} {
	sale := payload.sale
	items := payload.items
	masterData := payload.masterData
	employees := payload.employees

	assets := sale.Assets
	lead := items.leads[sale.LeadId.Hex()]
	note := items.notes[sale.ID.Hex()]

	LatestNote := ""
	if note != nil {
		LatestNote = note.Content
	}
	Reason := ""
	if len(sale.Tags) > 0 {
		for i := 0; i < len(sale.Tags); i++ {
			if tag := items.tags[sale.Tags[i]]; tag != nil {
				Reason += tag.Name + ", "
			}
		}
	}

	CreatedAt := (sale.CreatedAt).String()
	UpdatedAt := (sale.UpdatedAt).String()
	Code := sale.Code
	FullName := lead.FullName
	Phone := lead.Phone
	Email := lead.Email
	Type := (*masterData.types)[sale.Type]
	Source := (*masterData.sources)[sale.Source]
	Group := (*masterData.groups)[sale.Group]
	AssetType := (*masterData.assetTypes)[assets.AssetType]
	Description := assets.Description
	Employee := (*employees)[sale.EmployeeBy]
	StoreCode := sale.StoreCode
	Address := lead.Address
	Province := (*masterData.provinces)[lead.Province]
	Status := (*masterData.statuses)[sale.Status]
	DisbursedAmount := sale.DisbursedAmount
	ContractCode := sale.ContractCode
	CreatedName := (*employees)[sale.CreatedBy]
	UpdatedName := (*employees)[sale.UpdatedBy]
	DemandLoan := assets.DemandLoan
	if _, ok := DemandLoan.(int32); !ok {
		DemandLoan = 0
	}

	return []interface{}{
		CreatedAt, UpdatedAt, Code, FullName, Phone, Email, Type, Source, Group,
		AssetType, Description, DemandLoan, Employee, StoreCode, Address, Province,
		LatestNote, Status, Reason, DisbursedAmount, ContractCode, CreatedName, UpdatedName,
	}
}

func (s *ExportService) getMasterData(ctx context.Context, token string) iMasterData {
	var result []map[string]string

	var wg = sync.WaitGroup{}
	wg.Add(6)

	var assetTypes, sources, statuses, dataTypes, provinces, groups *map[string]string

	go func(ctx context.Context) {
		assetTypes = s.masterDataClient.GetAssetType(ctx)
		result = append(result, *assetTypes)
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		sources = s.masterDataClient.GetSource(ctx)
		result = append(result, *sources)
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		statuses = s.masterDataClient.GetStatuses(ctx)
		result = append(result, *statuses)
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		dataTypes = s.masterDataClient.GetTypes(ctx)
		result = append(result, *dataTypes)
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		groups = s.masterDataClient.GetGroups(ctx)
		result = append(result, *groups)
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		provinces = s.masterDataClient.GetProvinces(ctx)
		result = append(result, *provinces)
		wg.Done()
	}(ctx)

	wg.Wait()

	return iMasterData{
		assetTypes: assetTypes,
		sources:    sources,
		statuses:   statuses,
		types:      dataTypes,
		provinces:  provinces,
		groups:     groups,
	}
}

func (s *ExportService) getDataDatabase(
	ctx context.Context,
	saleIds []primitive.ObjectID,
	leadIds []primitive.ObjectID,
) getDataDB {
	var wg = sync.WaitGroup{}
	wg.Add(3)

	var (
		leads = make(map[string]*entities.Lead)
		tags  = make(map[string]*entities.Tag)
		notes = make(map[string]*entities.Note)
	)

	go func() {
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"createdAt", -1}})
		findOptions.SetProjection(bson.D{
			{"fullName", 1},
			{"phone", 1},
			{"email", 1},
			{"address", 1},
			{"province", 1},
			{"_id", 1},
		})
		items, _ := s.leadRepo.BaseRepo.Find(bson.M{
			"_id": bson.M{
				"$in": leadIds,
			},
		}, findOptions)

		for _, lead := range items {
			leads[lead.ID.Hex()] = lead
		}
		wg.Done()
	}()

	go func() {
		findOptions := options.Find()
		findOptions.SetProjection(bson.D{
			{"name", 1},
			{"code", 1},
		})
		items, _ := s.tagRepo.BaseRepo.Find(bson.M{}, findOptions)
		for _, item := range items {
			tags[item.Code] = item
		}
		wg.Done()
	}()

	go func() {
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"createdAt", -1}})
		findOptions.SetProjection(bson.D{{"content", 1}, {"_id", 1}})
		items, _ := s.noteRepo.BaseRepo.Find(bson.M{
			"saleOpportunitiesId": bson.M{
				"$in": saleIds,
			},
		}, findOptions)
		for _, item := range items {
			notes[item.SaleOpportunitiesId] = item
		}
		wg.Done()
	}()

	wg.Wait()

	return getDataDB{
		leads: leads,
		notes: notes,
		tags:  tags,
	}
}

func (s *ExportService) saveFile(values [][]interface{}, exportRequestId string) error {
	f := excelize.NewFile()
	for i, row := range values {
		startCell, err := excelize.JoinCellName("A", i+1)
		if err != nil {
			utils.Logger.Error(err)
			return nil
		}

		if err := f.SetSheetRow("Sheet1", startCell, &row); err != nil {
			utils.Logger.Error(err)
			return nil
		}
	}

	file, _ := f.WriteToBuffer()
	reader := io.Reader(file)
	result, err := s.uploader.UploadFile(reader, exportRequestId+".xlsx")
	if err != nil {
		return err
	}

	err = s.fileManagerClient.CreateFile(ctx, exportRequestId, result.Name, map[string]interface{}{
		"name": result.Name,
		"type": result.ContentType,
		"size": result.Size,
	})

	if err != nil {
		_ = s.fileManagerClient.UpdateExportRequestFailure(ctx, exportRequestId)
		return err
	}
	return nil
}
