package cdek

// service_types.go - Удобные типы для высокоуровневого Service API (PascalCase по регламенту)
//
// Эти типы являются обертками над автогенерированными CDEK типами,
// упрощающими использование библиотеки в приложениях.

// ========================
// Cost Calculation
// ========================

// CostRequest - запрос на расчет стоимости доставки
type CostRequest struct {
	FromCityCode int32     `json:"from_city_code,omitempty"` // Код города отправителя (по КЛАДР)
	ToCityCode   int32     `json:"to_city_code,omitempty"`   // Код города получателя (по КЛАДР)
	Packages     []Package `json:"packages,omitempty"`       // Список мест (упаковок)
}

// Package - информация об одной упаковке (месте)
type Package struct {
	Weight int32 `json:"weight,omitempty"` // Вес в граммах
	Length int32 `json:"length,omitempty"` // Длина в см
	Width  int32 `json:"width,omitempty"`  // Ширина в см
	Height int32 `json:"height,omitempty"` // Высота в см
}

// CostResponse - ответ с вариантами доставки и стоимостью
type CostResponse struct {
	Tariffs []TariffOption `json:"tariffs,omitempty"` // Доступные тарифы
}

// TariffOption - один вариант тарифа доставки
type TariffOption struct {
	TariffCode   int     `json:"tariff_code,omitempty"`   // Код тарифа
	TariffName   string  `json:"tariff_name,omitempty"`   // Название тарифа
	DeliveryMode int     `json:"delivery_mode,omitempty"` // Режим доставки: 1=дверь-дверь, 2=дверь-склад, 3=склад-дверь, 4=склад-склад
	DeliverySum  float64 `json:"delivery_sum,omitempty"`  // Стоимость доставки (рубли)
	PeriodMin    int     `json:"period_min,omitempty"`    // Минимальный срок доставки (дней)
	PeriodMax    int     `json:"period_max,omitempty"`    // Максимальный срок доставки (дней)
}

// AvailableTariffsRequest - запрос на получение списка доступных тарифов по договору
type AvailableTariffsRequest struct {
	Lang *string // Язык вывода информации о тарифах: rus, eng, zho (по умолчанию - rus)
}

// AvailableTariff - тариф, доступный и актуальный по договору
type AvailableTariff struct {
	TariffName              string                               `json:"tariff_name"`                            // Название тарифа
	DeliveryModes           []AvailableDeliveryMode              `json:"delivery_modes"`                         // Доступные режимы доставки для тарифа
	OrderTypes              []int                                `json:"order_types"`                            // Доступные типы заказов (пустой список - доступны все)
	PayerContragentType     []string                             `json:"payer_contragent_type"`                  // Доступные типы контрагентов-плательщиков
	SenderContragentType    []string                             `json:"sender_contragent_type"`                 // Доступные типы контрагентов-отправителей
	RecipientContragentType []string                             `json:"recipient_contragent_type"`              // Доступные типы контрагентов-получателей
	WeightMin               float64                              `json:"weight_min"`                             // Минимальный вес отправления
	WeightMax               float64                              `json:"weight_max"`                             // Максимальный вес отправления
	WeightCalcMax           float64                              `json:"weight_calc_max"`                        // Максимальный расчетный вес
	LengthMin               float64                              `json:"length_min"`                             // Минимальная длина упаковки
	LengthMax               float64                              `json:"length_max"`                             // Максимальная длина упаковки
	WidthMin                float64                              `json:"width_min"`                              // Минимальная ширина упаковки
	WidthMax                float64                              `json:"width_max"`                              // Максимальная ширина упаковки
	HeightMin               float64                              `json:"height_min"`                             // Минимальная высота упаковки
	HeightMax               float64                              `json:"height_max"`                             // Максимальная высота упаковки
	AdditionalOrderTypes    *AvailableTariffAdditionalOrderTypes `json:"additional_order_types_param,omitempty"` // Доп. типы заказа, применимые к тарифу
}

// AvailableTariffAdditionalOrderTypes - доп. типы заказа, применимые к тарифу
type AvailableTariffAdditionalOrderTypes struct {
	WithoutAdditionalOrderType bool  `json:"without_additional_order_type"` // Доступность тарифа для заказа без доп.типа
	AdditionalOrderTypes       []int `json:"additional_order_types"`        // Список доступных доп. типов заказа для тарифа
}

// AvailableDeliveryMode - режим доставки, доступный для тарифа
type AvailableDeliveryMode struct {
	DeliveryMode     int    `json:"delivery_mode"`      // Код режима доставки: 1=дверь-дверь, 2=дверь-склад, 3=склад-дверь, 4=склад-склад
	DeliveryModeName string `json:"delivery_mode_name"` // Название режима доставки
	TariffCode       int    `json:"tariff_code"`        // Код тарифа
}

// ========================
// Orders
// ========================

// Seller - истинный продавец (третье лицо)
// Используется для интернет-магазинов, когда фактический продавец отличается от отправителя
type Seller struct {
	Name          *string `json:"name,omitempty"`           // Наименование истинного продавца
	INN           *string `json:"inn,omitempty"`            // ИНН истинного продавца (10 или 12 символов)
	Phone         *string `json:"phone,omitempty"`          // Телефон истинного продавца
	OwnershipForm *int    `json:"ownership_form,omitempty"` // Код формы собственности
	Address       *string `json:"address,omitempty"`        // Адрес истинного продавца
}

// OrderRequest - запрос на создание заказа
type OrderRequest struct {
	Type         string         // Тип заказа: "delivery" (доставка), "pickup" (самовывоз)
	TariffCode   int            // Код тарифа
	Comment      *string        // Комментарий к заказу
	Sender       Recipient      // Отправитель (может быть компания с ИНН или физлицо с паспортом)
	Recipient    Recipient      // Получатель (может быть компания с ИНН или физлицо с паспортом)
	Seller       *Seller        // Истинный продавец (третье лицо) - только для интернет-магазинов
	FromLocation Location       // Адрес отправителя
	ToLocation   Location       // Адрес получателя
	Packages     []OrderPackage // Список мест
}

// Contact - контактная информация
type Contact struct {
	Company *string `json:"company,omitempty"` // Название компании
	Name    string  `json:"name,omitempty"`    // ФИО контактного лица
	Email   *string `json:"email,omitempty"`   // Email
	Phones  []Phone `json:"phones,omitempty"`  // Список телефонов
}

// Recipient - информация о получателе (расширяет Contact)
// Может быть как физическое лицо (Name + паспорт), так и компания (Company + ИНН)
type Recipient struct {
	Contact              `json:"contact"` // Базовая контактная информация (Name обязательно для физлица, Company для юрлица)
	TIN                  *string          `json:"tin,omitempty"`                    // ИНН (Tax Identification Number) - для юридических лиц и ИП (10 или 12 символов)
	PassportSeries       *string          `json:"passport_series,omitempty"`        // Серия паспорта (для физических лиц)
	PassportNumber       *string          `json:"passport_number,omitempty"`        // Номер паспорта (для физических лиц)
	PassportDateOfIssue  *string          `json:"passport_date_of_issue,omitempty"` // Дата выдачи паспорта (для физических лиц)
	PassportOrganization *string          `json:"passport_organization,omitempty"`  // Кем выдан паспорт (для физических лиц)
	PassportDateOfBirth  *string          `json:"passport_date_of_birth,omitempty"` // Дата рождения (yyyy-MM-dd) (для физических лиц)
}

// Phone - телефонный номер
type Phone struct {
	Number     string  `json:"number,omitempty"`     // Номер телефона (обязательно)
	Additional *string `json:"additional,omitempty"` // Добавочный номер
}

// Location - местоположение (адрес)
type Location struct {
	Code        *int32  `json:"code,omitempty"`         // Код населенного пункта СДЭК
	FiasGUID    *string `json:"fias_guid,omitempty"`    // Уникальный идентификатор ФИАС
	PostalCode  *string `json:"postal_code,omitempty"`  // Почтовый индекс
	CountryCode *string `json:"country_code,omitempty"` // Код страны (ISO 3166-1 alpha-2)
	Region      *string `json:"region,omitempty"`       // Регион
	City        *string `json:"city,omitempty"`         // Город
	Address     *string `json:"address,omitempty"`      // Адрес (улица, дом, квартира)
}

// OrderPackage - информация об упаковке в заказе
type OrderPackage struct {
	Number  string  `json:"number,omitempty"`  // Номер упаковки (артикул)
	Weight  int32   `json:"weight,omitempty"`  // Общий вес (граммы)
	Length  *int32  `json:"length,omitempty"`  // Длина (см)
	Width   *int32  `json:"width,omitempty"`   // Ширина (см)
	Height  *int32  `json:"height,omitempty"`  // Высота (см)
	Comment *string `json:"comment,omitempty"` // Комментарий
	Items   []Item  `json:"items,omitempty"`   // Список вложений
}

// Item - вложение в упаковку
type Item struct {
	Name    string  `json:"name,omitempty"`     // Наименование товара
	WareKey string  `json:"ware_key,omitempty"` // Артикул товара
	Payment float64 `json:"payment,omitempty"`  // Оплата (за единицу товара, в т.ч. частичная предоплата)
	Cost    float64 `json:"cost,omitempty"`     // Объявленная стоимость товара (за единицу)
	Weight  int32   `json:"weight,omitempty"`   // Вес (граммы, за единицу)
	Amount  int32   `json:"amount,omitempty"`   // Количество единиц товара
}

// OrderResponse - ответ при создании заказа
type OrderResponse struct {
	UUID       string        `json:"uuid,omitempty"`        // Идентификатор заказа в CDEK
	Number     *string       `json:"number,omitempty"`      // Номер заказа CDEK (может быть null до обработки)
	TariffCode int           `json:"tariff_code,omitempty"` // Код тарифа
	Statuses   []StatusEvent `json:"statuses,omitempty"`    // История статусов
	CreatedAt  string        `json:"created_at,omitempty"`  // Дата и время создания (ISO 8601)
}

// StatusEvent - событие изменения статуса
type StatusEvent struct {
	Code     string  `json:"code,omitempty"`      // Код статуса
	Name     string  `json:"name,omitempty"`      // Название статуса
	DateTime string  `json:"date_time,omitempty"` // Дата и время статуса (ISO 8601)
	City     *string `json:"city,omitempty"`      // Город, в котором произошло событие
}

// UpdateOrderRequest - запрос на обновление заказа
type UpdateOrderRequest struct {
	OrderUUID    string         // UUID заказа для обновления
	Recipient    *Recipient     // Новые данные получателя
	Sender       *Recipient     // Новые данные отправителя (может быть компания с ИНН или физлицо)
	Seller       *Seller        // Истинный продавец (третье лицо)
	ToLocation   *Location      // Новый адрес доставки
	FromLocation *Location      // Новый адрес отправления
	Comment      *string        // Новый комментарий
	Packages     []OrderPackage // Обновленный список мест (если нужно)
}

// OrderInfo - полная информация о заказе
type OrderInfo struct {
	UUID              string         `json:"uuid,omitempty"`               // Идентификатор заказа в CDEK
	Number            *string        `json:"number,omitempty"`             // Номер заказа CDEK
	Type              string         `json:"type,omitempty"`               // Тип заказа
	TariffCode        int            `json:"tariff_code,omitempty"`        // Код тарифа
	Sender            Recipient      `json:"sender"`                       // Отправитель (может быть компания с ИНН)
	Recipient         Recipient      `json:"recipient"`                    // Получатель (может быть компания с ИНН)
	Seller            *Seller        `json:"seller,omitempty"`             // Продавец (третье лицо, для интернет-магазинов)
	FromLocation      Location       `json:"from_location"`                // Адрес отправления
	ToLocation        Location       `json:"to_location"`                  // Адрес доставки
	Packages          []OrderPackage `json:"packages,omitempty"`           // Список мест
	Statuses          []StatusEvent  `json:"statuses,omitempty"`           // История статусов
	CreatedAt         string         `json:"created_at,omitempty"`         // Дата создания
	DeliveryCost      *float64       `json:"delivery_cost,omitempty"`      // Стоимость доставки
	EstimatedDelivery *string        `json:"estimated_delivery,omitempty"` // Планируемая дата доставки
	ActualDelivery    *string        `json:"actual_delivery,omitempty"`    // Фактическая дата доставки
}

// ========================
// Tracking
// ========================

// TrackingInfo - информация об отслеживании заказа
type TrackingInfo struct {
	UUID              string        `json:"uuid,omitempty"`               // Идентификатор заказа в CDEK
	Number            *string       `json:"number,omitempty"`             // Номер заказа CDEK
	CurrentStatus     StatusEvent   `json:"current_status"`               // Текущий статус
	StatusHistory     []StatusEvent `json:"status_history,omitempty"`     // История статусов
	EstimatedDelivery *string       `json:"estimated_delivery,omitempty"` // Планируемая дата доставки (ISO 8601)
	ActualDelivery    *string       `json:"actual_delivery,omitempty"`    // Фактическая дата доставки (ISO 8601)
}

// ========================
// Delivery Points
// ========================

// DeliveryPointsRequest - запрос на получение списка ПВЗ
type DeliveryPointsRequest struct {
	CityCode              string   // Код города (по КЛАДР)
	Type                  string   // Тип пункта: "PVZ" (пункт выдачи), "POSTAMAT" (постамат)
	Code                  string   // Код ПВЗ
	PostalCode            *string  // Почтовый индекс города, для которого необходим список офисов
	CountryCode           *string  // Код страны в формате ISO_3166-1_alpha-2
	RegionCode            *int     // Код региона СДЭК
	HaveCashless          *bool    // Наличие терминала оплаты
	HaveCash              *bool    // Есть прием наличных
	AllowedCod            *bool    // Разрешен наложенный платеж
	IsDressingRoom        *bool    // Наличие примерочной
	WeightMax             *float64 // Максимальный вес в кг, который может принять офис
	WeightMin             *float64 // Минимальный вес в кг, который принимает офис
	Lang                  *string  // Локализация офиса
	TakeOnly              *bool    // Является ли офис только пунктом выдачи
	IsHandout             *bool    // Является пунктом выдачи
	IsReception           *bool    // Есть ли в офисе приём заказов
	IsMarketplace         *bool    // Офис для доставки "До маркетплейса"
	IsLtl                 *bool    // Работает ли офис с LTL (сборный груз)
	LtlAcceptancePartners *bool    // Принимает заказы LTL, которые будут доставляться партнерами
	LtlIssuancePartners   *bool    // Выдает заказы LTL, которые были доставлены партнерами
	Fulfillment           *bool    // Офис с зоной фулфилмента
	FiasGuid              *string  // Код города ФИАС
	Size                  *int     // Ограничение выборки результата (размер страницы)
	Page                  *int     // Номер страницы выборки результата
}

// DeliveryPoint - пункт выдачи заказов
type DeliveryPoint struct {
	Code                  string              `json:"code,omitempty"`                     // Код ПВЗ
	UUID                  *string             `json:"uuid,omitempty"`                     // Идентификатор офиса в ИС СДЭК
	Name                  string              `json:"name,omitempty"`                     // Название ПВЗ
	Type                  string              `json:"type,omitempty"`                     // Тип: "PVZ", "POSTAMAT"
	Location              PointLocation       `json:"location"`                           // Адрес расположения
	NearestStation        *string             `json:"nearest_station,omitempty"`          // Ближайшая станция/остановка транспорта
	WorkTime              string              `json:"work_time,omitempty"`                // Режим работы
	WorkTimeList          []WorkTimeEntry     `json:"work_time_list,omitempty"`           // График работы по дням недели
	WorkTimeExceptionList []WorkTimeException `json:"work_time_exception_list,omitempty"` // Исключения в графике работы офиса
	Phones                []Phone             `json:"phones,omitempty"`                   // Телефоны
	Email                 *string             `json:"email,omitempty"`                    // Email
	Note                  *string             `json:"note,omitempty"`                     // Примечание
	OwnerCode             *string             `json:"owner_code,omitempty"`               // Принадлежность офиса компании (CDEK, PickPoint, ...)
	TakeOnly              bool                `json:"take_only,omitempty"`                // Является ли офис только пунктом выдачи
	IsHandout             bool                `json:"is_handout,omitempty"`               // Является пунктом выдачи
	IsReception           bool                `json:"is_reception,omitempty"`             // Является пунктом приёма
	IsDressingRoom        bool                `json:"is_dressing_room,omitempty"`         // Есть ли примерочная
	IsLtl                 bool                `json:"is_ltl,omitempty"`                   // Работает ли офис с LTL (сборный груз)
	HaveCashless          bool                `json:"have_cashless,omitempty"`            // Есть безналичный расчет
	HaveCash              bool                `json:"have_cash,omitempty"`                // Есть приём наличных
	HaveFastPaymentSystem bool                `json:"have_fast_payment_system,omitempty"` // Есть безналичный расчёт по СБП
	AllowedCod            bool                `json:"allowed_cod,omitempty"`              // Разрешен наложенный платеж в ПВЗ
	Site                  *string             `json:"site,omitempty"`                     // Ссылка на офис на сайте СДЭК
	OfficeImage           *string             `json:"office_image,omitempty"`             // URL первого изображения офиса (для обратной совместимости, см. OfficeImageList)
	OfficeImageList       []string            `json:"office_image_list,omitempty"`        // Все фото офиса (кроме фото проезда)
	WeightMin             *float64            `json:"weight_min,omitempty"`               // Минимальный вес (кг), принимаемый в ПВЗ
	WeightMax             *float64            `json:"weight_max,omitempty"`               // Максимальный вес (кг), принимаемый в ПВЗ
	Status                string              `json:"status,omitempty"`                   // Статус офиса: "ACTIVE", "CLOSED"
}

// WorkTimeEntry - график работы офиса на конкретный день недели
type WorkTimeEntry struct {
	Day  int    `json:"day,omitempty"`  // Порядковый номер дня недели (1 = понедельник, 7 = воскресенье)
	Time string `json:"time,omitempty"` // Период работы в этот день
}

// WorkTimeException - исключение в графике работы офиса на определенный период
type WorkTimeException struct {
	DateStart string `json:"date_start,omitempty"` // Дата начала исключения
	DateEnd   string `json:"date_end,omitempty"`   // Дата окончания исключения
	TimeStart string `json:"time_start,omitempty"` // Время начала работы в указанную дату
	TimeEnd   string `json:"time_end,omitempty"`   // Время окончания работы в указанную дату
	IsWorking bool   `json:"is_working,omitempty"` // Признак рабочего/нерабочего дня
}

// PointLocation - местоположение пункта выдачи
type PointLocation struct {
	CountryCode *string `json:"country_code,omitempty"` // Код страны в формате ISO_3166-1_alpha-2
	Region      string  `json:"region,omitempty"`       // Регион
	RegionCode  *int32  `json:"region_code,omitempty"`  // Код региона СДЭК
	City        string  `json:"city,omitempty"`         // Город
	CityCode    *int32  `json:"city_code,omitempty"`    // Код населенного пункта СДЭК
	CityUUID    *string `json:"city_uuid,omitempty"`    // Идентификатор города в ИС СДЭК
	FiasGUID    *string `json:"fias_guid,omitempty"`    // Идентификатор ФИАС населенного пункта
	Address     string  `json:"address,omitempty"`      // Адрес
	AddressFull *string `json:"address_full,omitempty"` // Полный адрес с указанием страны, региона, города и т.д.
	PostalCode  string  `json:"postal_code,omitempty"`  // Почтовый индекс
	Latitude    float64 `json:"latitude,omitempty"`     // Широта
	Longitude   float64 `json:"longitude,omitempty"`    // Долгота
}

// ========================
// Print (Barcode/Waybill)
// ========================

// PrintBarcodeRequest - запрос на создание этикеток
type PrintBarcodeRequest struct {
	Orders []PrintOrder // Список заказов для печати
	Copy   *int         // Количество копий (по умолчанию 1)
	Format *string      // Формат: "A4", "A5", "A6" (по умолчанию A4)
}

// PrintWaybillRequest - запрос на создание накладных
type PrintWaybillRequest struct {
	Orders []PrintOrder // Список заказов для печати
	Copy   *int         // Количество копий (по умолчанию 1)
	Format *string      // Формат: "A4", "A5" (по умолчанию A4)
}

// PrintOrder - заказ для печати
type PrintOrder struct {
	OrderUUID string `json:"order_uuid,omitempty"` // UUID заказа
}

// PrintResponse - ответ на запрос печати
type PrintResponse struct {
	UUID      string `json:"uuid,omitempty"`       // UUID задания на печать
	URL       string `json:"url,omitempty"`        // URL для скачивания PDF (доступен после готовности)
	Status    string `json:"status,omitempty"`     // Статус: "ACCEPTED", "PROCESSING", "READY", "INVALID"
	CreatedAt string `json:"created_at,omitempty"` // Дата создания задания
}

// ========================
// Location Reference (Cities/Regions)
// ========================

// CitiesRequest - запрос на получение списка городов
type CitiesRequest struct {
	CountryCode    *string  // Код страны (ISO 3166-1 alpha-2)
	RegionCode     *int     // Код региона
	FiasRegionGUID *string  // ФИАС код региона. Устаревшее поле - значения могут быть не актуальны
	KladrCode      *string  // Код КЛАДР населенного пункта
	FiasGUID       *string  // ФИАС код населенного пункта
	PostalCode     *string  // Почтовый индекс
	Code           *int     // Код населенного пункта СДЭК
	City           *string  // Название города (поиск)
	PaymentLimit   *float64 // Ограничение на сумму наложенного платежа. -1 - ограничения нет; 0 - наложенный платеж не принимается
	Size           *int     // Количество результатов (по умолчанию 1000)
	Page           *int     // Номер страницы
	Lang           *string  // Язык локализации ответа
}

// City - город из справочника СДЭК
type City struct {
	Code         int      `json:"code,omitempty"`          // Код населенного пункта СДЭК
	City         string   `json:"city,omitempty"`          // Название города
	FiasGUID     *string  `json:"fias_guid,omitempty"`     // Уникальный идентификатор ФИАС
	Region       string   `json:"region,omitempty"`        // Регион
	RegionCode   int      `json:"region_code,omitempty"`   // Код региона
	Country      string   `json:"country,omitempty"`       // Страна
	CountryCode  string   `json:"country_code,omitempty"`  // Код страны
	Latitude     float64  `json:"latitude,omitempty"`      // Широта
	Longitude    float64  `json:"longitude,omitempty"`     // Долгота
	TimeZone     string   `json:"time_zone,omitempty"`     // Часовой пояс
	PaymentLimit float64  `json:"payment_limit,omitempty"` // Ограничение оплаты наличными
	PostalCodes  []string `json:"postal_codes,omitempty"`  // Почтовые индексы
}

// RegionsRequest - запрос на получение списка регионов
type RegionsRequest struct {
	CountryCode *string // Код страны (ISO 3166-1 alpha-2)
	RegionCode  *int    // Код региона
	Region      *string // Название региона (поиск)
	Size        *int    // Количество результатов
	Page        *int    // Номер страницы
}

// Region - регион из справочника СДЭК
type Region struct {
	Code        int     `json:"code,omitempty"`         // Код региона
	Region      string  `json:"region,omitempty"`       // Название региона
	Country     string  `json:"country,omitempty"`      // Страна
	CountryCode string  `json:"country_code,omitempty"` // Код страны
	FiasGUID    *string `json:"fias_guid,omitempty"`    // ФИАС код региона
}

// ========================
// Intake (Заявка на забор)
// ========================

// IntakeRequest - запрос на создание заявки на забор груза
type IntakeRequest struct {
	IntakeDate     string        // Дата ожидаемого забора (ISO 8601: YYYY-MM-DD)
	IntakeTimeFrom string        // Время начала ожидания (HH:MM)
	IntakeTimeTo   string        // Время окончания ожидания (HH:MM)
	LunchTimeFrom  *string       // Время начала обеда (HH:MM)
	LunchTimeTo    *string       // Время окончания обеда (HH:MM)
	Comment        *string       // Комментарий
	Sender         Contact       // Отправитель
	FromLocation   Location      // Адрес забора
	NeedCall       *bool         // Нужен ли звонок
	Orders         []IntakeOrder // Список заказов для забора
}

// IntakeOrder - заказ в заявке на забор
type IntakeOrder struct {
	OrderUUID string // UUID заказа
}

// IntakeResponse - ответ при создании заявки на забор
type IntakeResponse struct {
	UUID       string `json:"uuid,omitempty"`        // UUID заявки
	Number     string `json:"number,omitempty"`      // Номер заявки СДЭК
	IntakeDate string `json:"intake_date,omitempty"` // Дата забора
	Status     string `json:"status,omitempty"`      // Статус заявки
	CreatedAt  string `json:"created_at,omitempty"`  // Дата создания
}

// IntakeInfo - информация о заявке на забор
type IntakeInfo struct {
	UUID           string        `json:"uuid,omitempty"`             // UUID заявки
	Number         string        `json:"number,omitempty"`           // Номер заявки
	IntakeDate     string        `json:"intake_date,omitempty"`      // Дата забора
	IntakeTimeFrom string        `json:"intake_time_from,omitempty"` // Время начала
	IntakeTimeTo   string        `json:"intake_time_to,omitempty"`   // Время окончания
	Status         string        `json:"status,omitempty"`           // Статус
	Sender         Contact       `json:"sender"`                     // Отправитель
	FromLocation   Location      `json:"from_location"`              // Адрес забора
	Orders         []IntakeOrder `json:"orders,omitempty"`           // Заказы
	CreatedAt      string        `json:"created_at,omitempty"`       // Дата создания
}

// ========================
// Webhooks
// ========================

// Типы событий webhook
const (
	WebhookTypeOrderStatus         = "ORDER_STATUS"         // Изменение статуса заказа
	WebhookTypeOrderModified       = "ORDER_MODIFIED"       // Изменение заказа
	WebhookTypePrintForm           = "PRINT_FORM"           // Готовность печатной формы
	WebhookTypeReceipt             = "RECEIPT"              // Квитанция
	WebhookTypePrealertClosed      = "PREALERT_CLOSED"      // Закрытие преалерта
	WebhookTypeAccompanyingWaybill = "ACCOMPANYING_WAYBILL" // Информация о транспорте
	WebhookTypeOfficeAvailability  = "OFFICE_AVAILABILITY"  // Доступность офиса
	WebhookTypeDelivAgreement      = "DELIV_AGREEMENT"      // Договоренность о доставке
	WebhookTypeDelivProblem        = "DELIV_PROBLEM"        // Проблемы доставки
	WebhookTypeCourierInfo         = "COURIER_INFO"         // Информация о курьере
)

// WebhookRequest - запрос на создание webhook
type WebhookRequest struct {
	URL  string // URL для получения уведомлений (обязательно)
	Type string // Тип события (обязательно): ORDER_STATUS, PRINT_FORM, etc.
}

// WebhookResponse - ответ при создании webhook
type WebhookResponse struct {
	UUID string // UUID созданного webhook
}

// Webhook - информация о webhook
type Webhook struct {
	UUID string // UUID webhook
	URL  string // URL для уведомлений
	Type string // Тип события
}
