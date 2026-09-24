package cdek

import (
	"encoding/json"
	"fmt"
)

// dtoMapper - компонент преобразования между Service API и автогенерированными CDEK типами (camelCase по регламенту 4.6)
type dtoMapper struct{}

// newDtoMapper создает новый маппер DTO
func newDtoMapper() *dtoMapper {
	return &dtoMapper{}
}

// ========================
// Calculator (Cost Estimation)
// ========================

// toCDEKCalculatorRequest преобразует CostRequest → CalculatorTariffListRequestDto
func (m *dtoMapper) toCDEKCalculatorRequest(req *CostRequest) CalculatorTariffListRequestDto {
	// Тип услуги: 1 = интернет-магазин
	serviceType := int32(1)
	// Валюта: 1 = рубли
	currency := int32(1)

	// Преобразование упаковок
	packages := make([]CalcPackageRequestDto, len(req.Packages))
	for i, pkg := range req.Packages {
		packages[i] = CalcPackageRequestDto{
			Weight: pkg.Weight,
			Length: &pkg.Length,
			Width:  &pkg.Width,
			Height: &pkg.Height,
		}
	}

	return CalculatorTariffListRequestDto{
		Type:     &serviceType,
		Currency: &currency,
		FromLocation: CalculatorLocationDto{
			Code: &req.FromCityCode,
		},
		ToLocation: CalculatorLocationDto{
			Code: &req.ToCityCode,
		},
		Packages: packages,
	}
}

// fromCDEKCalculatorResponse преобразует CalculatorTariffListResponseDto → CostResponse
func (m *dtoMapper) fromCDEKCalculatorResponse(data []byte) (*CostResponse, error) {
	// Парсим как map т.к. автогенерированные типы используют interface{}
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal calculator response: %w", err)
	}

	// Проверяем наличие тарифов
	tariffCodes, ok := rawResp["tariff_codes"].([]interface{})
	if !ok || len(tariffCodes) == 0 {
		return nil, fmt.Errorf("no tariffs available for this route")
	}

	// Преобразование тарифов
	tariffs := make([]TariffOption, 0, len(tariffCodes))
	for _, t := range tariffCodes {
		tariffMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}

		tariff := TariffOption{}

		if code, ok := tariffMap["tariff_code"].(float64); ok {
			tariff.TariffCode = int(code)
		}
		if name, ok := tariffMap["tariff_name"].(string); ok {
			tariff.TariffName = name
		}
		if mode, ok := tariffMap["delivery_mode"].(float64); ok {
			tariff.DeliveryMode = int(mode)
		}
		if sum, ok := tariffMap["delivery_sum"].(float64); ok {
			tariff.DeliverySum = sum
		}
		if periodMin, ok := tariffMap["period_min"].(float64); ok {
			tariff.PeriodMin = int(periodMin)
		}
		if periodMax, ok := tariffMap["period_max"].(float64); ok {
			tariff.PeriodMax = int(periodMax)
		}

		tariffs = append(tariffs, tariff)
	}

	return &CostResponse{
		Tariffs: tariffs,
	}, nil
}

// fromCDEKAvailableTariffs преобразует CalculatorAvailableTariffsResponseDto → []AvailableTariff
func (m *dtoMapper) fromCDEKAvailableTariffs(data []byte) ([]AvailableTariff, error) {
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal available tariffs: %w", err)
	}

	tariffCodes, ok := rawResp["tariff_codes"].([]interface{})
	if !ok {
		return nil, nil
	}

	tariffs := make([]AvailableTariff, 0, len(tariffCodes))
	for _, t := range tariffCodes {
		tariffMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}

		tariff := AvailableTariff{}

		if name, ok := tariffMap["tariff_name"].(string); ok {
			tariff.TariffName = name
		}
		if v, ok := tariffMap["weight_min"].(float64); ok {
			tariff.WeightMin = v
		}
		if v, ok := tariffMap["weight_max"].(float64); ok {
			tariff.WeightMax = v
		}
		if v, ok := tariffMap["weight_calc_max"].(float64); ok {
			tariff.WeightCalcMax = v
		}
		if v, ok := tariffMap["length_min"].(float64); ok {
			tariff.LengthMin = v
		}
		if v, ok := tariffMap["length_max"].(float64); ok {
			tariff.LengthMax = v
		}
		if v, ok := tariffMap["width_min"].(float64); ok {
			tariff.WidthMin = v
		}
		if v, ok := tariffMap["width_max"].(float64); ok {
			tariff.WidthMax = v
		}
		if v, ok := tariffMap["height_min"].(float64); ok {
			tariff.HeightMin = v
		}
		if v, ok := tariffMap["height_max"].(float64); ok {
			tariff.HeightMax = v
		}

		tariff.OrderTypes = toIntSlice(tariffMap["order_types"])
		tariff.PayerContragentType = toStringSlice(tariffMap["payer_contragent_type"])
		tariff.SenderContragentType = toStringSlice(tariffMap["sender_contragent_type"])
		tariff.RecipientContragentType = toStringSlice(tariffMap["recipient_contragent_type"])

		if modes, ok := tariffMap["delivery_modes"].([]interface{}); ok {
			tariff.DeliveryModes = make([]AvailableDeliveryMode, 0, len(modes))
			for _, dm := range modes {
				modeMap, ok := dm.(map[string]interface{})
				if !ok {
					continue
				}
				mode := AvailableDeliveryMode{}
				if v, ok := modeMap["delivery_mode"].(float64); ok {
					mode.DeliveryMode = int(v)
				}
				if v, ok := modeMap["delivery_mode_name"].(string); ok {
					mode.DeliveryModeName = v
				}
				if v, ok := modeMap["tariff_code"].(float64); ok {
					mode.TariffCode = int(v)
				}
				tariff.DeliveryModes = append(tariff.DeliveryModes, mode)
			}
		}

		if param, ok := tariffMap["additional_order_types_param"].(map[string]interface{}); ok {
			additional := &AvailableTariffAdditionalOrderTypes{
				AdditionalOrderTypes: toIntSlice(param["additional_order_types"]),
			}
			if v, ok := param["without_additional_order_type"].(bool); ok {
				additional.WithoutAdditionalOrderType = v
			}
			tariff.AdditionalOrderTypes = additional
		}

		tariffs = append(tariffs, tariff)
	}

	return tariffs, nil
}

// toIntSlice преобразует []interface{} с числами в []int
func toIntSlice(v interface{}) []int {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]int, 0, len(raw))
	for _, item := range raw {
		if n, ok := item.(float64); ok {
			result = append(result, int(n))
		}
	}
	return result
}

// toStringSlice преобразует []interface{} со строками в []string
func toStringSlice(v interface{}) []string {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// ========================
// Orders
// ========================

// toCDEKOrderRequest преобразует OrderRequest → map для OrderCreateRequestDto
// Используем map[string]interface{} т.к. автогенерированные типы используют interface{}
func (m *dtoMapper) toCDEKOrderRequest(req *OrderRequest) (map[string]interface{}, error) {
	order := make(map[string]interface{})

	// Обязательные поля — CDEK v2 expects integer type: 1=online_store, 2=delivery
	switch req.Type {
	case "delivery", "2":
		order["type"] = 2
	case "online_store", "1":
		order["type"] = 1
	default:
		order["type"] = 1
	}
	order["tariff_code"] = req.TariffCode

	// Опциональный комментарий
	if req.Comment != nil {
		order["comment"] = *req.Comment
	}

	// Отправитель (теперь с поддержкой ИНН и паспортных данных)
	if req.Sender.Name != "" {
		sender := map[string]interface{}{
			"name": req.Sender.Name,
		}
		if req.Sender.Company != nil {
			sender["company"] = *req.Sender.Company
		}
		if req.Sender.Email != nil {
			sender["email"] = *req.Sender.Email
		}
		if len(req.Sender.Phones) > 0 {
			phones := make([]map[string]interface{}, len(req.Sender.Phones))
			for i, phone := range req.Sender.Phones {
				p := map[string]interface{}{"number": phone.Number}
				if phone.Additional != nil {
					p["additional"] = *phone.Additional
				}
				phones[i] = p
			}
			sender["phones"] = phones
		}
		// ИНН для отправителя (компания или ИП)
		if req.Sender.TIN != nil {
			sender["tin"] = *req.Sender.TIN
		}
		// Паспортные данные для отправителя (физлицо)
		if req.Sender.PassportSeries != nil {
			sender["passport_series"] = *req.Sender.PassportSeries
		}
		if req.Sender.PassportNumber != nil {
			sender["passport_number"] = *req.Sender.PassportNumber
		}
		if req.Sender.PassportDateOfIssue != nil {
			sender["passport_date_of_issue"] = *req.Sender.PassportDateOfIssue
		}
		if req.Sender.PassportOrganization != nil {
			sender["passport_organization"] = *req.Sender.PassportOrganization
		}
		if req.Sender.PassportDateOfBirth != nil {
			sender["passport_date_of_birth"] = *req.Sender.PassportDateOfBirth
		}
		order["sender"] = sender
	}

	// Получатель (обязательный)
	recipient := map[string]interface{}{
		"name": req.Recipient.Name,
	}
	if req.Recipient.Company != nil {
		recipient["company"] = *req.Recipient.Company
	}
	if req.Recipient.Email != nil {
		recipient["email"] = *req.Recipient.Email
	}
	if len(req.Recipient.Phones) > 0 {
		phones := make([]map[string]interface{}, len(req.Recipient.Phones))
		for i, phone := range req.Recipient.Phones {
			p := map[string]interface{}{"number": phone.Number}
			if phone.Additional != nil {
				p["additional"] = *phone.Additional
			}
			phones[i] = p
		}
		recipient["phones"] = phones
	}
	// ИНН (для юридических лиц и ИП)
	if req.Recipient.TIN != nil {
		recipient["tin"] = *req.Recipient.TIN
	}
	// Паспортные данные (для физических лиц)
	if req.Recipient.PassportSeries != nil {
		recipient["passport_series"] = *req.Recipient.PassportSeries
	}
	if req.Recipient.PassportNumber != nil {
		recipient["passport_number"] = *req.Recipient.PassportNumber
	}
	if req.Recipient.PassportDateOfIssue != nil {
		recipient["passport_date_of_issue"] = *req.Recipient.PassportDateOfIssue
	}
	if req.Recipient.PassportOrganization != nil {
		recipient["passport_organization"] = *req.Recipient.PassportOrganization
	}
	if req.Recipient.PassportDateOfBirth != nil {
		recipient["passport_date_of_birth"] = *req.Recipient.PassportDateOfBirth
	}
	order["recipient"] = recipient

	// Истинный продавец (третье лицо) - только для интернет-магазинов
	if req.Seller != nil {
		seller := make(map[string]interface{})
		if req.Seller.Name != nil {
			seller["name"] = *req.Seller.Name
		}
		if req.Seller.INN != nil {
			seller["inn"] = *req.Seller.INN
		}
		if req.Seller.Phone != nil {
			seller["phone"] = *req.Seller.Phone
		}
		if req.Seller.OwnershipForm != nil {
			seller["ownership_form"] = *req.Seller.OwnershipForm
		}
		if req.Seller.Address != nil {
			seller["address"] = *req.Seller.Address
		}
		order["seller"] = seller
	}

	// Адрес отправителя
	if req.FromLocation.Code != nil || req.FromLocation.Address != nil {
		fromLoc := make(map[string]interface{})
		if req.FromLocation.Code != nil {
			fromLoc["code"] = *req.FromLocation.Code
		}
		if req.FromLocation.Address != nil {
			fromLoc["address"] = *req.FromLocation.Address
		}
		if req.FromLocation.City != nil {
			fromLoc["city"] = *req.FromLocation.City
		}
		order["from_location"] = fromLoc
	}

	// Адрес получателя
	if req.ToLocation.Code != nil || req.ToLocation.Address != nil {
		toLoc := make(map[string]interface{})
		if req.ToLocation.Code != nil {
			toLoc["code"] = *req.ToLocation.Code
		}
		if req.ToLocation.Address != nil {
			toLoc["address"] = *req.ToLocation.Address
		}
		if req.ToLocation.City != nil {
			toLoc["city"] = *req.ToLocation.City
		}
		order["to_location"] = toLoc
	}

	// Упаковки (обязательно)
	packages := make([]map[string]interface{}, len(req.Packages))
	for i, pkg := range req.Packages {
		p := map[string]interface{}{
			"number": pkg.Number,
			"weight": pkg.Weight,
		}
		if pkg.Length != nil {
			p["length"] = *pkg.Length
		}
		if pkg.Width != nil {
			p["width"] = *pkg.Width
		}
		if pkg.Height != nil {
			p["height"] = *pkg.Height
		}

		// Товары в упаковке
		items := make([]map[string]interface{}, len(pkg.Items))
		for j, item := range pkg.Items {
			items[j] = map[string]interface{}{
				"name":     item.Name,
				"ware_key": item.WareKey,
				"payment":  map[string]interface{}{"value": item.Payment},
				"cost":     item.Cost,
				"weight":   item.Weight,
				"amount":   item.Amount,
			}
		}
		p["items"] = items
		packages[i] = p
	}
	order["packages"] = packages

	return order, nil
}

// fromCDEKOrderResponse преобразует ответ API → OrderResponse
func (m *dtoMapper) fromCDEKOrderResponse(data []byte) (*OrderResponse, error) {
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal order response: %w", err)
	}

	orderResp := &OrderResponse{}

	// CreateOrder возвращает либо entity (успех), либо requests (создание запроса)
	// Пробуем entity сначала (для успешного создания)
	if entity, ok := rawResp["entity"].(map[string]interface{}); ok {
		// UUID (обязательный)
		if uuid, ok := entity["uuid"].(string); ok {
			orderResp.UUID = uuid
		}

		// Номер заказа CDEK (может быть null)
		if cdekNum, ok := entity["cdek_number"].(string); ok {
			orderResp.Number = &cdekNum
		}

		// Код тарифа
		if tariff, ok := entity["tariff_code"].(float64); ok {
			orderResp.TariffCode = int(tariff)
		}

		// Статусы
		if statuses, ok := entity["statuses"].([]interface{}); ok {
			orderResp.Statuses = make([]StatusEvent, 0, len(statuses))
			for _, s := range statuses {
				if status, ok := s.(map[string]interface{}); ok {
					event := StatusEvent{}
					if code, ok := status["code"].(string); ok {
						event.Code = code
					}
					if name, ok := status["name"].(string); ok {
						event.Name = name
					}
					if dt, ok := status["date_time"].(string); ok {
						event.DateTime = dt
					}
					if city, ok := status["city"].(string); ok {
						event.City = &city
					}
					orderResp.Statuses = append(orderResp.Statuses, event)
				}
			}
		}

		// Дата создания из requests внутри entity
		if requests, ok := entity["requests"].([]interface{}); ok && len(requests) > 0 {
			if firstReq, ok := requests[0].(map[string]interface{}); ok {
				if dt, ok := firstReq["date_time"].(string); ok {
					orderResp.CreatedAt = dt
				}
			}
		}
	}

	// Если нет entity, используем requests на верхнем уровне
	// Это происходит при асинхронном создании заказа
	if orderResp.UUID == "" {
		// Ищем UUID в related_entities
		if relatedEntities, ok := rawResp["related_entities"].([]interface{}); ok && len(relatedEntities) > 0 {
			for _, rel := range relatedEntities {
				if relMap, ok := rel.(map[string]interface{}); ok {
					if uuid, ok := relMap["uuid"].(string); ok {
						orderResp.UUID = uuid
						break
					}
				}
			}
		}

		// Дата создания из requests на верхнем уровне
		if requests, ok := rawResp["requests"].([]interface{}); ok && len(requests) > 0 {
			if firstReq, ok := requests[0].(map[string]interface{}); ok {
				if dt, ok := firstReq["date_time"].(string); ok {
					orderResp.CreatedAt = dt
				}
			}
		}
	}

	if orderResp.UUID == "" {
		return nil, fmt.Errorf("uuid not found in response")
	}

	return orderResp, nil
}

// fromCDEKOrderToTracking преобразует GetOrder ответ → TrackingInfo
func (m *dtoMapper) fromCDEKOrderToTracking(data []byte) (*TrackingInfo, error) {
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal tracking response: %w", err)
	}

	// Проверяем наличие entity
	entity, ok := rawResp["entity"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("entity not found in response")
	}

	tracking := &TrackingInfo{}

	// UUID (обязательный)
	if uuid, ok := entity["uuid"].(string); ok {
		tracking.UUID = uuid
	}

	// Номер заказа
	if cdekNum, ok := entity["cdek_number"].(string); ok {
		tracking.Number = &cdekNum
	}

	// История статусов
	if statuses, ok := entity["statuses"].([]interface{}); ok && len(statuses) > 0 {
		tracking.StatusHistory = make([]StatusEvent, 0, len(statuses))
		for _, s := range statuses {
			if status, ok := s.(map[string]interface{}); ok {
				event := StatusEvent{}
				if code, ok := status["code"].(string); ok {
					event.Code = code
				}
				if name, ok := status["name"].(string); ok {
					event.Name = name
				}
				if dt, ok := status["date_time"].(string); ok {
					event.DateTime = dt
				}
				if city, ok := status["city"].(string); ok {
					event.City = &city
				}
				tracking.StatusHistory = append(tracking.StatusHistory, event)
			}
		}

		// Текущий статус - последний в списке
		if len(tracking.StatusHistory) > 0 {
			tracking.CurrentStatus = tracking.StatusHistory[len(tracking.StatusHistory)-1]
		}
	}

	// Плановая дата доставки
	if date, ok := entity["planned_delivery_date"].(string); ok {
		tracking.EstimatedDelivery = &date
	}

	// Фактическая дата доставки (из информации о вручении)
	if delivery, ok := entity["delivery_detail"].(map[string]interface{}); ok {
		if date, ok := delivery["date"].(string); ok {
			tracking.ActualDelivery = &date
		}
	}

	return tracking, nil
}

// ========================
// Delivery Points
// ========================

// fromCDEKDeliveryPoints преобразует GetDeliverypoints ответ → []DeliveryPoint
func (m *dtoMapper) fromCDEKDeliveryPoints(data []byte) ([]DeliveryPoint, error) {
	var rawResp []map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal delivery points: %w", err)
	}

	points := make([]DeliveryPoint, 0, len(rawResp))
	for _, p := range rawResp {
		point := DeliveryPoint{}

		if code, ok := p["code"].(string); ok {
			point.Code = code
		}
		if uuid, ok := p["uuid"].(string); ok {
			point.UUID = &uuid
		}
		if name, ok := p["name"].(string); ok {
			point.Name = name
		}
		if pvzType, ok := p["type"].(string); ok {
			point.Type = pvzType
		}
		if nearestStation, ok := p["nearest_station"].(string); ok {
			point.NearestStation = &nearestStation
		}
		if workTime, ok := p["work_time"].(string); ok {
			point.WorkTime = workTime
		}
		if email, ok := p["email"].(string); ok {
			point.Email = &email
		}
		if note, ok := p["note"].(string); ok {
			point.Note = &note
		}
		if ownerCode, ok := p["owner_code"].(string); ok {
			point.OwnerCode = &ownerCode
		}
		if takeOnly, ok := p["take_only"].(bool); ok {
			point.TakeOnly = takeOnly
		}
		if isHandout, ok := p["is_handout"].(bool); ok {
			point.IsHandout = isHandout
		}
		if isReception, ok := p["is_reception"].(bool); ok {
			point.IsReception = isReception
		}
		if isDressingRoom, ok := p["is_dressing_room"].(bool); ok {
			point.IsDressingRoom = isDressingRoom
		}
		if isLtl, ok := p["is_ltl"].(bool); ok {
			point.IsLtl = isLtl
		}
		if haveCashless, ok := p["have_cashless"].(bool); ok {
			point.HaveCashless = haveCashless
		}
		if haveCash, ok := p["have_cash"].(bool); ok {
			point.HaveCash = haveCash
		}
		if haveFastPaymentSystem, ok := p["have_fast_payment_system"].(bool); ok {
			point.HaveFastPaymentSystem = haveFastPaymentSystem
		}
		if allowedCod, ok := p["allowed_cod"].(bool); ok {
			point.AllowedCod = allowedCod
		}
		if site, ok := p["site"].(string); ok {
			point.Site = &site
		}
		if weightMin, ok := p["weight_min"].(float64); ok {
			point.WeightMin = &weightMin
		}
		if weightMax, ok := p["weight_max"].(float64); ok {
			point.WeightMax = &weightMax
		}
		if status, ok := p["status"].(string); ok {
			point.Status = status
		}

		// График работы по дням недели
		if workTimeList, ok := p["work_time_list"].([]interface{}); ok {
			point.WorkTimeList = make([]WorkTimeEntry, 0, len(workTimeList))
			for _, wt := range workTimeList {
				if entry, ok := wt.(map[string]interface{}); ok {
					workTimeEntry := WorkTimeEntry{}
					if day, ok := entry["day"].(float64); ok {
						workTimeEntry.Day = int(day)
					}
					if t, ok := entry["time"].(string); ok {
						workTimeEntry.Time = t
					}
					point.WorkTimeList = append(point.WorkTimeList, workTimeEntry)
				}
			}
		}

		// Исключения в графике работы
		if exceptionList, ok := p["work_time_exception_list"].([]interface{}); ok {
			point.WorkTimeExceptionList = make([]WorkTimeException, 0, len(exceptionList))
			for _, exc := range exceptionList {
				if entry, ok := exc.(map[string]interface{}); ok {
					exception := WorkTimeException{}
					if dateStart, ok := entry["date_start"].(string); ok {
						exception.DateStart = dateStart
					}
					if dateEnd, ok := entry["date_end"].(string); ok {
						exception.DateEnd = dateEnd
					}
					if timeStart, ok := entry["time_start"].(string); ok {
						exception.TimeStart = timeStart
					}
					if timeEnd, ok := entry["time_end"].(string); ok {
						exception.TimeEnd = timeEnd
					}
					if isWorking, ok := entry["is_working"].(bool); ok {
						exception.IsWorking = isWorking
					}
					point.WorkTimeExceptionList = append(point.WorkTimeExceptionList, exception)
				}
			}
		}

		// Телефоны
		if phones, ok := p["phones"].([]interface{}); ok {
			point.Phones = make([]Phone, 0, len(phones))
			for _, ph := range phones {
				if phone, ok := ph.(map[string]interface{}); ok {
					phoneObj := Phone{}
					if number, ok := phone["number"].(string); ok {
						phoneObj.Number = number
					}
					if add, ok := phone["additional"].(string); ok {
						phoneObj.Additional = &add
					}
					point.Phones = append(point.Phones, phoneObj)
				}
			}
		}

		// Местоположение
		if loc, ok := p["location"].(map[string]interface{}); ok {
			if countryCode, ok := loc["country_code"].(string); ok {
				point.Location.CountryCode = &countryCode
			}
			if region, ok := loc["region"].(string); ok {
				point.Location.Region = region
			}
			if regionCode, ok := loc["region_code"].(float64); ok {
				regionCodeInt := int32(regionCode)
				point.Location.RegionCode = &regionCodeInt
			}
			if city, ok := loc["city"].(string); ok {
				point.Location.City = city
			}
			if cityCode, ok := loc["city_code"].(float64); ok {
				cityCodeInt := int32(cityCode)
				point.Location.CityCode = &cityCodeInt
			}
			if cityUUID, ok := loc["city_uuid"].(string); ok {
				point.Location.CityUUID = &cityUUID
			}
			if fiasGUID, ok := loc["fias_guid"].(string); ok {
				point.Location.FiasGUID = &fiasGUID
			}
			if address, ok := loc["address"].(string); ok {
				point.Location.Address = address
			}
			if addressFull, ok := loc["address_full"].(string); ok {
				point.Location.AddressFull = &addressFull
			}
			if postal, ok := loc["postal_code"].(string); ok {
				point.Location.PostalCode = postal
			}
			if lat, ok := loc["latitude"].(float64); ok {
				point.Location.Latitude = lat
			}
			if lon, ok := loc["longitude"].(float64); ok {
				point.Location.Longitude = lon
			}
		}

		// Изображения офиса
		if images, ok := p["office_image_list"].([]interface{}); ok && len(images) > 0 {
			point.OfficeImageList = make([]string, 0, len(images))
			for _, img := range images {
				if imgMap, ok := img.(map[string]interface{}); ok {
					if url, ok := imgMap["url"].(string); ok {
						point.OfficeImageList = append(point.OfficeImageList, url)
					}
				}
			}
			if firstImg, ok := images[0].(map[string]interface{}); ok {
				if url, ok := firstImg["url"].(string); ok {
					point.OfficeImage = &url
				}
			}
		}

		points = append(points, point)
	}

	return points, nil
}

// ========================
// Order Info (GetOrder)
// ========================

// fromCDEKOrderToInfo преобразует GetOrder ответ → OrderInfo (полная информация)
func (m *dtoMapper) fromCDEKOrderToInfo(data []byte) (*OrderInfo, error) {
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal order info: %w", err)
	}

	// Проверяем наличие entity
	entity, ok := rawResp["entity"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("entity not found in response")
	}

	info := &OrderInfo{}

	// UUID (обязательный)
	if uuid, ok := entity["uuid"].(string); ok {
		info.UUID = uuid
	}

	// Номер заказа
	if cdekNum, ok := entity["cdek_number"].(string); ok {
		info.Number = &cdekNum
	}

	// Тип заказа
	if orderType, ok := entity["type"].(string); ok {
		info.Type = orderType
	}

	// Код тарифа
	if tariff, ok := entity["tariff_code"].(float64); ok {
		info.TariffCode = int(tariff)
	}

	// Отправитель (теперь Recipient с поддержкой ИНН и паспорта)
	if sender, ok := entity["sender"].(map[string]interface{}); ok {
		if name, ok := sender["name"].(string); ok {
			info.Sender.Name = name
		}
		if company, ok := sender["company"].(string); ok {
			info.Sender.Company = &company
		}
		if email, ok := sender["email"].(string); ok {
			info.Sender.Email = &email
		}
		if phones, ok := sender["phones"].([]interface{}); ok {
			info.Sender.Phones = make([]Phone, 0, len(phones))
			for _, ph := range phones {
				if phone, ok := ph.(map[string]interface{}); ok {
					phoneObj := Phone{}
					if number, ok := phone["number"].(string); ok {
						phoneObj.Number = number
					}
					if add, ok := phone["additional"].(string); ok {
						phoneObj.Additional = &add
					}
					info.Sender.Phones = append(info.Sender.Phones, phoneObj)
				}
			}
		}
		// ИНН (для юридических лиц)
		if tin, ok := sender["tin"].(string); ok {
			info.Sender.TIN = &tin
		}
		// Паспортные данные (для физических лиц)
		if passportSeries, ok := sender["passport_series"].(string); ok {
			info.Sender.PassportSeries = &passportSeries
		}
		if passportNumber, ok := sender["passport_number"].(string); ok {
			info.Sender.PassportNumber = &passportNumber
		}
		if passportDateOfIssue, ok := sender["passport_date_of_issue"].(string); ok {
			info.Sender.PassportDateOfIssue = &passportDateOfIssue
		}
		if passportOrg, ok := sender["passport_organization"].(string); ok {
			info.Sender.PassportOrganization = &passportOrg
		}
		if passportDOB, ok := sender["passport_date_of_birth"].(string); ok {
			info.Sender.PassportDateOfBirth = &passportDOB
		}
	}

	// Получатель
	if recipient, ok := entity["recipient"].(map[string]interface{}); ok {
		if name, ok := recipient["name"].(string); ok {
			info.Recipient.Name = name
		}
		if company, ok := recipient["company"].(string); ok {
			info.Recipient.Company = &company
		}
		if email, ok := recipient["email"].(string); ok {
			info.Recipient.Email = &email
		}
		if phones, ok := recipient["phones"].([]interface{}); ok {
			info.Recipient.Phones = make([]Phone, 0, len(phones))
			for _, ph := range phones {
				if phone, ok := ph.(map[string]interface{}); ok {
					phoneObj := Phone{}
					if number, ok := phone["number"].(string); ok {
						phoneObj.Number = number
					}
					if add, ok := phone["additional"].(string); ok {
						phoneObj.Additional = &add
					}
					info.Recipient.Phones = append(info.Recipient.Phones, phoneObj)
				}
			}
		}
		// ИНН (для юридических лиц)
		if tin, ok := recipient["tin"].(string); ok {
			info.Recipient.TIN = &tin
		}
		// Паспортные данные (для физических лиц)
		if passportSeries, ok := recipient["passport_series"].(string); ok {
			info.Recipient.PassportSeries = &passportSeries
		}
		if passportNumber, ok := recipient["passport_number"].(string); ok {
			info.Recipient.PassportNumber = &passportNumber
		}
		if passportDateOfIssue, ok := recipient["passport_date_of_issue"].(string); ok {
			info.Recipient.PassportDateOfIssue = &passportDateOfIssue
		}
		if passportOrg, ok := recipient["passport_organization"].(string); ok {
			info.Recipient.PassportOrganization = &passportOrg
		}
		if passportDOB, ok := recipient["passport_date_of_birth"].(string); ok {
			info.Recipient.PassportDateOfBirth = &passportDOB
		}
	}

	// Продавец (третье лицо, для интернет-магазинов)
	if seller, ok := entity["seller"].(map[string]interface{}); ok {
		info.Seller = &Seller{}
		if name, ok := seller["name"].(string); ok {
			info.Seller.Name = &name
		}
		if inn, ok := seller["inn"].(string); ok {
			info.Seller.INN = &inn
		}
		if phone, ok := seller["phone"].(string); ok {
			info.Seller.Phone = &phone
		}
		if ownership, ok := seller["ownership_form"].(float64); ok {
			ownershipInt := int(ownership)
			info.Seller.OwnershipForm = &ownershipInt
		}
		if address, ok := seller["address"].(string); ok {
			info.Seller.Address = &address
		}
	}

	// Статусы
	if statuses, ok := entity["statuses"].([]interface{}); ok && len(statuses) > 0 {
		info.Statuses = make([]StatusEvent, 0, len(statuses))
		for _, s := range statuses {
			if status, ok := s.(map[string]interface{}); ok {
				event := StatusEvent{}
				if code, ok := status["code"].(string); ok {
					event.Code = code
				}
				if name, ok := status["name"].(string); ok {
					event.Name = name
				}
				if dt, ok := status["date_time"].(string); ok {
					event.DateTime = dt
				}
				if city, ok := status["city"].(string); ok {
					event.City = &city
				}
				info.Statuses = append(info.Statuses, event)
			}
		}
	}

	// Адрес отправления
	if fromLoc, ok := entity["from_location"].(map[string]interface{}); ok {
		info.FromLocation = parseCDEKLocation(fromLoc)
	}

	// Адрес доставки
	if toLoc, ok := entity["to_location"].(map[string]interface{}); ok {
		info.ToLocation = parseCDEKLocation(toLoc)
	}

	// Список мест (упаковок)
	if packages, ok := entity["packages"].([]interface{}); ok {
		info.Packages = make([]OrderPackage, 0, len(packages))
		for _, p := range packages {
			if pkgMap, ok := p.(map[string]interface{}); ok {
				info.Packages = append(info.Packages, parseCDEKOrderPackage(pkgMap))
			}
		}
	}

	// Дата создания из requests
	if requests, ok := entity["requests"].([]interface{}); ok && len(requests) > 0 {
		if firstReq, ok := requests[0].(map[string]interface{}); ok {
			if dt, ok := firstReq["date_time"].(string); ok {
				info.CreatedAt = dt
			}
		}
	}

	// Плановая дата доставки
	if date, ok := entity["planned_delivery_date"].(string); ok {
		info.EstimatedDelivery = &date
	}

	// Стоимость доставки и фактическая дата доставки (из информации о вручении)
	if delivery, ok := entity["delivery_detail"].(map[string]interface{}); ok {
		if sum, ok := delivery["delivery_sum"].(float64); ok {
			info.DeliveryCost = &sum
		}
		if date, ok := delivery["date"].(string); ok {
			info.ActualDelivery = &date
		}
	}

	return info, nil
}

// parseCDEKLocation преобразует map с данными адреса → Location
func parseCDEKLocation(loc map[string]interface{}) Location {
	location := Location{}
	if code, ok := loc["code"].(float64); ok {
		codeInt := int32(code)
		location.Code = &codeInt
	}
	if fiasGUID, ok := loc["fias_guid"].(string); ok {
		location.FiasGUID = &fiasGUID
	}
	if postalCode, ok := loc["postal_code"].(string); ok {
		location.PostalCode = &postalCode
	}
	if countryCode, ok := loc["country_code"].(string); ok {
		location.CountryCode = &countryCode
	}
	if region, ok := loc["region"].(string); ok {
		location.Region = &region
	}
	if city, ok := loc["city"].(string); ok {
		location.City = &city
	}
	if address, ok := loc["address"].(string); ok {
		location.Address = &address
	}
	return location
}

// parseCDEKOrderPackage преобразует map с данными упаковки → OrderPackage
func parseCDEKOrderPackage(p map[string]interface{}) OrderPackage {
	pkg := OrderPackage{}
	if number, ok := p["number"].(string); ok {
		pkg.Number = number
	}
	if weight, ok := p["weight"].(float64); ok {
		pkg.Weight = int32(weight)
	}
	if length, ok := p["length"].(float64); ok {
		lengthInt := int32(length)
		pkg.Length = &lengthInt
	}
	if width, ok := p["width"].(float64); ok {
		widthInt := int32(width)
		pkg.Width = &widthInt
	}
	if height, ok := p["height"].(float64); ok {
		heightInt := int32(height)
		pkg.Height = &heightInt
	}
	if comment, ok := p["comment"].(string); ok {
		pkg.Comment = &comment
	}
	if items, ok := p["items"].([]interface{}); ok {
		pkg.Items = make([]Item, 0, len(items))
		for _, it := range items {
			if itemMap, ok := it.(map[string]interface{}); ok {
				item := Item{}
				if name, ok := itemMap["name"].(string); ok {
					item.Name = name
				}
				if wareKey, ok := itemMap["ware_key"].(string); ok {
					item.WareKey = wareKey
				}
				if payment, ok := itemMap["payment"].(map[string]interface{}); ok {
					if value, ok := payment["value"].(float64); ok {
						item.Payment = value
					}
				}
				if cost, ok := itemMap["cost"].(float64); ok {
					item.Cost = cost
				}
				if weight, ok := itemMap["weight"].(float64); ok {
					item.Weight = int32(weight)
				}
				if amount, ok := itemMap["amount"].(float64); ok {
					item.Amount = int32(amount)
				}
				pkg.Items = append(pkg.Items, item)
			}
		}
	}
	return pkg
}

// ========================
// Update Order
// ========================

// toCDEKUpdateOrderRequest преобразует UpdateOrderRequest → map для Update API
func (m *dtoMapper) toCDEKUpdateOrderRequest(req *UpdateOrderRequest) (map[string]interface{}, error) {
	update := make(map[string]interface{})

	// UUID обязательный для обновления
	update["uuid"] = req.OrderUUID

	// Получатель (если указан)
	if req.Recipient != nil {
		recipient := map[string]interface{}{
			"name": req.Recipient.Name,
		}
		if req.Recipient.Company != nil {
			recipient["company"] = *req.Recipient.Company
		}
		if req.Recipient.Email != nil {
			recipient["email"] = *req.Recipient.Email
		}
		if len(req.Recipient.Phones) > 0 {
			phones := make([]map[string]interface{}, len(req.Recipient.Phones))
			for i, phone := range req.Recipient.Phones {
				p := map[string]interface{}{"number": phone.Number}
				if phone.Additional != nil {
					p["additional"] = *phone.Additional
				}
				phones[i] = p
			}
			recipient["phones"] = phones
		}
		// ИНН для получателя
		if req.Recipient.TIN != nil {
			recipient["tin"] = *req.Recipient.TIN
		}
		// Паспортные данные для получателя
		if req.Recipient.PassportSeries != nil {
			recipient["passport_series"] = *req.Recipient.PassportSeries
		}
		if req.Recipient.PassportNumber != nil {
			recipient["passport_number"] = *req.Recipient.PassportNumber
		}
		if req.Recipient.PassportDateOfIssue != nil {
			recipient["passport_date_of_issue"] = *req.Recipient.PassportDateOfIssue
		}
		if req.Recipient.PassportOrganization != nil {
			recipient["passport_organization"] = *req.Recipient.PassportOrganization
		}
		if req.Recipient.PassportDateOfBirth != nil {
			recipient["passport_date_of_birth"] = *req.Recipient.PassportDateOfBirth
		}
		update["recipient"] = recipient
	}

	// Отправитель (если указан) - теперь с ИНН и паспортными данными
	if req.Sender != nil {
		sender := map[string]interface{}{
			"name": req.Sender.Name,
		}
		if req.Sender.Company != nil {
			sender["company"] = *req.Sender.Company
		}
		if req.Sender.Email != nil {
			sender["email"] = *req.Sender.Email
		}
		if len(req.Sender.Phones) > 0 {
			phones := make([]map[string]interface{}, len(req.Sender.Phones))
			for i, phone := range req.Sender.Phones {
				p := map[string]interface{}{"number": phone.Number}
				if phone.Additional != nil {
					p["additional"] = *phone.Additional
				}
				phones[i] = p
			}
			sender["phones"] = phones
		}
		// ИНН для отправителя
		if req.Sender.TIN != nil {
			sender["tin"] = *req.Sender.TIN
		}
		// Паспортные данные для отправителя
		if req.Sender.PassportSeries != nil {
			sender["passport_series"] = *req.Sender.PassportSeries
		}
		if req.Sender.PassportNumber != nil {
			sender["passport_number"] = *req.Sender.PassportNumber
		}
		if req.Sender.PassportDateOfIssue != nil {
			sender["passport_date_of_issue"] = *req.Sender.PassportDateOfIssue
		}
		if req.Sender.PassportOrganization != nil {
			sender["passport_organization"] = *req.Sender.PassportOrganization
		}
		if req.Sender.PassportDateOfBirth != nil {
			sender["passport_date_of_birth"] = *req.Sender.PassportDateOfBirth
		}
		update["sender"] = sender
	}

	// Истинный продавец (третье лицо) - если указан
	if req.Seller != nil {
		seller := make(map[string]interface{})
		if req.Seller.Name != nil {
			seller["name"] = *req.Seller.Name
		}
		if req.Seller.INN != nil {
			seller["inn"] = *req.Seller.INN
		}
		if req.Seller.Phone != nil {
			seller["phone"] = *req.Seller.Phone
		}
		if req.Seller.OwnershipForm != nil {
			seller["ownership_form"] = *req.Seller.OwnershipForm
		}
		if req.Seller.Address != nil {
			seller["address"] = *req.Seller.Address
		}
		update["seller"] = seller
	}

	// Адрес получателя (если указан)
	if req.ToLocation != nil {
		toLoc := make(map[string]interface{})
		if req.ToLocation.Code != nil {
			toLoc["code"] = *req.ToLocation.Code
		}
		if req.ToLocation.Address != nil {
			toLoc["address"] = *req.ToLocation.Address
		}
		if req.ToLocation.City != nil {
			toLoc["city"] = *req.ToLocation.City
		}
		update["to_location"] = toLoc
	}

	// Адрес отправителя (если указан)
	if req.FromLocation != nil {
		fromLoc := make(map[string]interface{})
		if req.FromLocation.Code != nil {
			fromLoc["code"] = *req.FromLocation.Code
		}
		if req.FromLocation.Address != nil {
			fromLoc["address"] = *req.FromLocation.Address
		}
		if req.FromLocation.City != nil {
			fromLoc["city"] = *req.FromLocation.City
		}
		update["from_location"] = fromLoc
	}

	// Комментарий (если указан)
	if req.Comment != nil {
		update["comment"] = *req.Comment
	}

	// Упаковки (если указаны)
	if len(req.Packages) > 0 {
		packages := make([]map[string]interface{}, len(req.Packages))
		for i, pkg := range req.Packages {
			p := map[string]interface{}{
				"number": pkg.Number,
				"weight": pkg.Weight,
			}
			if pkg.Length != nil {
				p["length"] = *pkg.Length
			}
			if pkg.Width != nil {
				p["width"] = *pkg.Width
			}
			if pkg.Height != nil {
				p["height"] = *pkg.Height
			}

			// Товары в упаковке
			items := make([]map[string]interface{}, len(pkg.Items))
			for j, item := range pkg.Items {
				items[j] = map[string]interface{}{
					"name":     item.Name,
					"ware_key": item.WareKey,
					"payment":  map[string]interface{}{"value": item.Payment},
					"cost":     item.Cost,
					"weight":   item.Weight,
					"amount":   item.Amount,
				}
			}
			p["items"] = items
			packages[i] = p
		}
		update["packages"] = packages
	}

	return update, nil
}

// ========================
// Location Reference (Cities/Regions)
// ========================

// fromCDEKCities преобразует Cities ответ → []City
func (m *dtoMapper) fromCDEKCities(data []byte) ([]City, error) {
	var rawResp []map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal cities: %w", err)
	}

	cities := make([]City, 0, len(rawResp))
	for _, c := range rawResp {
		city := City{}
		if code, ok := c["code"].(float64); ok {
			city.Code = int(code)
		}
		if cityName, ok := c["city"].(string); ok {
			city.City = cityName
		}
		if fiasGUID, ok := c["fias_guid"].(string); ok {
			city.FiasGUID = &fiasGUID
		}
		if region, ok := c["region"].(string); ok {
			city.Region = region
		}
		if regionCode, ok := c["region_code"].(float64); ok {
			city.RegionCode = int(regionCode)
		}
		if country, ok := c["country"].(string); ok {
			city.Country = country
		}
		if countryCode, ok := c["country_code"].(string); ok {
			city.CountryCode = countryCode
		}
		if latitude, ok := c["latitude"].(float64); ok {
			city.Latitude = latitude
		}
		if longitude, ok := c["longitude"].(float64); ok {
			city.Longitude = longitude
		}
		if timeZone, ok := c["time_zone"].(string); ok {
			city.TimeZone = timeZone
		}
		if paymentLimit, ok := c["payment_limit"].(float64); ok {
			city.PaymentLimit = paymentLimit
		}
		cities = append(cities, city)
	}
	return cities, nil
}

// fromCDEKRegions преобразует Regions ответ → []Region
func (m *dtoMapper) fromCDEKRegions(data []byte) ([]Region, error) {
	var rawResp []map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal regions: %w", err)
	}

	regions := make([]Region, 0, len(rawResp))
	for _, r := range rawResp {
		region := Region{}
		if code, ok := r["region_code"].(float64); ok {
			region.Code = int(code)
		}
		if regionName, ok := r["region"].(string); ok {
			region.Region = regionName
		}
		if country, ok := r["country"].(string); ok {
			region.Country = country
		}
		if countryCode, ok := r["country_code"].(string); ok {
			region.CountryCode = countryCode
		}
		if fiasGUID, ok := r["fias_region_guid"].(string); ok {
			region.FiasGUID = &fiasGUID
		}
		regions = append(regions, region)
	}
	return regions, nil
}

// ========================
// Intakes
// ========================

// toCDEKIntakeRequest преобразует IntakeRequest → map для Intake API
func (m *dtoMapper) toCDEKIntakeRequest(req *IntakeRequest) (map[string]interface{}, error) {
	intake := map[string]interface{}{
		"intake_date":      req.IntakeDate,
		"intake_time_from": req.IntakeTimeFrom,
		"intake_time_to":   req.IntakeTimeTo,
		"sender":           map[string]interface{}{"name": req.Sender.Name},
	}
	if req.Comment != nil {
		intake["comment"] = *req.Comment
	}
	return intake, nil
}

// fromCDEKIntakeResponse преобразует Intake ответ → IntakeResponse
func (m *dtoMapper) fromCDEKIntakeResponse(data []byte) (*IntakeResponse, error) {
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal intake response: %w", err)
	}
	intakeResp := &IntakeResponse{}
	if entity, ok := rawResp["entity"].(map[string]interface{}); ok {
		if uuid, ok := entity["uuid"].(string); ok {
			intakeResp.UUID = uuid
		}
	}
	// Статус и дата создания - из информации о запросе над сущностью
	if requests, ok := rawResp["requests"].([]interface{}); ok && len(requests) > 0 {
		if firstReq, ok := requests[0].(map[string]interface{}); ok {
			if state, ok := firstReq["state"].(string); ok {
				intakeResp.Status = state
			}
			if dt, ok := firstReq["date_time"].(string); ok {
				intakeResp.CreatedAt = dt
			}
		}
	}
	return intakeResp, nil
}

// fromCDEKIntakeInfo преобразует GetIntake ответ → IntakeInfo
func (m *dtoMapper) fromCDEKIntakeInfo(data []byte) (*IntakeInfo, error) {
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal intake info: %w", err)
	}
	info := &IntakeInfo{}
	if entity, ok := rawResp["entity"].(map[string]interface{}); ok {
		if uuid, ok := entity["uuid"].(string); ok {
			info.UUID = uuid
		}
		if number, ok := entity["intake_number"].(string); ok {
			info.Number = number
		}
		if intakeDate, ok := entity["intake_date"].(string); ok {
			info.IntakeDate = intakeDate
		}
		if timeFrom, ok := entity["intake_time_from"].(string); ok {
			info.IntakeTimeFrom = timeFrom
		}
		if timeTo, ok := entity["intake_time_to"].(string); ok {
			info.IntakeTimeTo = timeTo
		}

		// Отправитель
		if sender, ok := entity["sender"].(map[string]interface{}); ok {
			if name, ok := sender["name"].(string); ok {
				info.Sender.Name = name
			}
			if company, ok := sender["company"].(string); ok {
				info.Sender.Company = &company
			}
			if email, ok := sender["email"].(string); ok {
				info.Sender.Email = &email
			}
			if phones, ok := sender["phones"].([]interface{}); ok {
				info.Sender.Phones = make([]Phone, 0, len(phones))
				for _, ph := range phones {
					if phone, ok := ph.(map[string]interface{}); ok {
						phoneObj := Phone{}
						if number, ok := phone["number"].(string); ok {
							phoneObj.Number = number
						}
						if add, ok := phone["additional"].(string); ok {
							phoneObj.Additional = &add
						}
						info.Sender.Phones = append(info.Sender.Phones, phoneObj)
					}
				}
			}
		}

		// Адрес забора
		if fromLoc, ok := entity["from_location"].(map[string]interface{}); ok {
			info.FromLocation = parseCDEKLocation(fromLoc)
		}

		// Заказ, к которому привязана заявка
		if orderUUID, ok := entity["order_uuid"].(string); ok {
			info.Orders = []IntakeOrder{{OrderUUID: orderUUID}}
		}

		// Статус - последний в списке
		if statuses, ok := entity["statuses"].([]interface{}); ok && len(statuses) > 0 {
			if lastStatus, ok := statuses[len(statuses)-1].(map[string]interface{}); ok {
				if name, ok := lastStatus["name"].(string); ok {
					info.Status = name
				}
			}
		}
	}

	// Дата создания из requests
	if requests, ok := rawResp["requests"].([]interface{}); ok && len(requests) > 0 {
		if firstReq, ok := requests[0].(map[string]interface{}); ok {
			if dt, ok := firstReq["date_time"].(string); ok {
				info.CreatedAt = dt
			}
		}
	}

	return info, nil
}

// ========================
// Webhooks
// ========================

// fromCDEKWebhookResponse преобразует CreateWebhook ответ → WebhookResponse
func (m *dtoMapper) fromCDEKWebhookResponse(data []byte) (*WebhookResponse, error) {
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal webhook response: %w", err)
	}
	webhookResp := &WebhookResponse{}
	if entity, ok := rawResp["entity"].(map[string]interface{}); ok {
		if uuid, ok := entity["uuid"].(string); ok {
			webhookResp.UUID = uuid
		}
	}
	return webhookResp, nil
}

// fromCDEKWebhook преобразует GetWebhook ответ → Webhook
func (m *dtoMapper) fromCDEKWebhook(data []byte) (*Webhook, error) {
	var rawResp map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal webhook: %w", err)
	}

	webhook := &Webhook{}
	if entity, ok := rawResp["entity"].(map[string]interface{}); ok {
		if uuid, ok := entity["uuid"].(string); ok {
			webhook.UUID = uuid
		}
		if url, ok := entity["url"].(string); ok {
			webhook.URL = url
		}
		if webhookType, ok := entity["type"].(string); ok {
			webhook.Type = webhookType
		}
	}
	return webhook, nil
}

// fromCDEKWebhooks преобразует ListWebhooks ответ → []Webhook
func (m *dtoMapper) fromCDEKWebhooks(data []byte) ([]Webhook, error) {
	var rawResp []map[string]interface{}
	if err := json.Unmarshal(data, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal webhooks: %w", err)
	}
	webhooks := make([]Webhook, 0, len(rawResp))
	for _, w := range rawResp {
		webhook := Webhook{}
		if uuid, ok := w["uuid"].(string); ok {
			webhook.UUID = uuid
		}
		if url, ok := w["url"].(string); ok {
			webhook.URL = url
		}
		if webhookType, ok := w["type"].(string); ok {
			webhook.Type = webhookType
		}
		webhooks = append(webhooks, webhook)
	}
	return webhooks, nil
}
