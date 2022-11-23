package services

import (
	"crm-worker-go/clients"
	"crm-worker-go/entities"
	"crm-worker-go/repositories"
	"crm-worker-go/types"
	"crm-worker-go/utils"
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"strings"
	"sync"
	"time"
)

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

type GetDataDB struct {
	leads map[string]*entities.Lead
	notes map[string]*entities.Note
	tags  map[string]*entities.Tag
}

type IMasterData struct {
	assetTypes *map[string]string
	sources    *map[string]string
	statuses   *map[string]string
	types      *map[string]string
	provinces  *map[string]string
	groups     *map[string]string
}

type IJobs struct {
	sale       *entities.SaleOpportunity
	items      GetDataDB
	masterData IMasterData
	employees  *map[string]string
}

func (s *ExportService) ExportData(payload types.PayloadMessageExport) bool {
	token := payload.Token
	_, err := utils.ExtractClaims(token)
	if err != nil {
		utils.Logger.Error(err)
		return false
	}

	s.fileManagerClient.Client.SetToken(token)
	s.masterDataClient.Client.SetToken(token)
	s.employeeClient.Client.SetToken(token)

	exportRequest, err := s.fileManagerClient.FindExportRequests(payload.ID)

	if err != nil {
		utils.Logger.Error("Get Export Request Failure ", err, "With Id: ", payload.ID)
		return false
	}

	if utils.Contains([]string{types.Completed, types.Failure}, exportRequest.Status) {
		return true
	}

	switch exportRequest.Type {
	case types.ExportRequestTypeSaleOpp:
		return s.exportSaleOpp(exportRequest)
	case types.ExportRequestTypeLead:
		return s.exportLead(exportRequest)
	}
	return false
}

func (s *ExportService) exportSaleOpp(payload *types.IExportRequest) bool {
	requestData := payload.Data

	var wg = sync.WaitGroup{}
	wg.Add(3)

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

	var items GetDataDB
	go func() {
		items = s.getDataDatabase(saleIds, leadIds)
		wg.Done()
	}()

	var employees *map[string]string
	go func() {
		employees, _ = s.employeeClient.GetEmployees(employeeIds)
		wg.Done()
	}()

	var masterData IMasterData
	go func() {
		masterData = s.getMasterData()
		wg.Done()
	}()
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

	numberWorker := 4
	jobs, result := make(chan IJobs, 100), make(chan []interface{}, 100)

	for i := 1; i <= numberWorker; i++ {
		go worker(jobs, result, fmt.Sprintf("%d", i))
	}

	go func(results []*entities.SaleOpportunity) {
		for i, sale := range results {
			jobs <- IJobs{
				sale:       sale,
				items:      items,
				masterData: masterData,
				employees:  employees,
			}
			fmt.Printf("[JOB] ===>>> %d has been enqueued with sale code %v \n", i, sale.Code)
		}
		close(jobs)
	}(results)

	for i := 0; i < len(results); i++ {
		fmt.Printf("[STT] ===>>> %d has been export with sale code %v \n", i, results[i].Code)
		exportData = append(exportData, <-result)
	}

	err := s.saveFile(exportData, payload.ID)
	if err != nil {
		utils.Logger.Error(err)
		return false
	}
	return true
}

func (s *ExportService) exportLead(payload *types.IExportRequest) bool {
	requestData := payload.Data

	var wg = sync.WaitGroup{}
	wg.Add(3)

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

	results, _ := s.leadRepo.BaseRepo.Find(filter, nil)
	if len(results) == 0 {
		return true
	}

	exportData := [][]interface{}{
		{
			"Full Name",
			"Phone",
			"Email",
			"Source",
			"Account Store",
			"Address",
		},
	}

	for _, lead := range results {
		exportData = append(exportData, []interface{}{
			lead.FullName, lead.Phone, lead.Email, lead.Source, lead.EmployeeBy, lead.Address,
		})
	}
	err := s.saveFile(exportData, payload.ID)
	if err != nil {
		utils.Logger.Error(err)
		return false
	}
	return true
}

func worker(jobs <-chan IJobs, results chan<- []interface{}, name string) {
	for n := range jobs {
		fmt.Printf("Worker %s is handle sale code %v\n", name, n.sale.Code)
		results <- handleData(n)
	}
}

func handleData(
	payload IJobs,
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
	Phone := utils.MaskStr(lead.Phone, 3, 4, "*")
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

func (s *ExportService) getMasterData() IMasterData {
	var result []map[string]string

	var wg = sync.WaitGroup{}
	wg.Add(6)

	var assetTypes, sources, statuses, dataTypes, provinces, groups *map[string]string

	go func() {
		assetTypes = s.masterDataClient.GetAssetType()
		result = append(result, *assetTypes)
		wg.Done()
	}()

	go func() {
		sources = s.masterDataClient.GetSource()
		result = append(result, *sources)
		wg.Done()
	}()

	go func() {
		statuses = s.masterDataClient.GetStatuses()
		result = append(result, *statuses)
		wg.Done()
	}()

	go func() {
		dataTypes = s.masterDataClient.GetTypes()
		result = append(result, *dataTypes)
		wg.Done()
	}()

	go func() {
		groups = s.masterDataClient.GetGroups()
		result = append(result, *groups)
		wg.Done()
	}()

	go func() {
		provinces = s.masterDataClient.GetProvinces()
		result = append(result, *provinces)
		wg.Done()
	}()

	wg.Wait()

	return IMasterData{
		assetTypes: assetTypes,
		sources:    sources,
		statuses:   statuses,
		types:      dataTypes,
		provinces:  provinces,
		groups:     groups,
	}
}

func (s *ExportService) getDataDatabase(
	saleIds []primitive.ObjectID,
	leadIds []primitive.ObjectID,
) GetDataDB {
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

	return GetDataDB{
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

		if err = f.SetSheetRow("Sheet1", startCell, &row); err != nil {
			utils.Logger.Error(err)
			return nil
		}
	}
	//if err := f.SaveAs("Book1.xlsx"); err != nil {
	//	fmt.Println(err)
	//}
	file, _ := f.WriteToBuffer()
	reader := io.Reader(file)
	result, err := s.uploader.UploadFile(reader, exportRequestId+".xlsx")
	if err != nil {
		_ = s.fileManagerClient.UpdateExportRequestFailure(exportRequestId)
		utils.Logger.Error(err)
		return err
	}

	_ = s.fileManagerClient.CreateFile(exportRequestId, result.Name, map[string]interface{}{
		"name": result.Name,
		"type": result.ContentType,
		"size": result.Size,
	})
	return nil
}
