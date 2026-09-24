package cdek

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDtoMapper_FromCDEKDeliveryPoints(t *testing.T) {
	raw := []byte(`[{
		"code":"NSK1","name":"NSK1, Новосибирск, ул. Кривощековская","uuid":"33258e52-0000-0000-0000-000000000000",
		"nearest_station":"Автовокзал","work_time":"Пн-Пт 08:00-21:00, Сб-Вс 10:00-20:00",
		"phones":[{"number":"+73832022250"}],"email":"nsk@cdek.ru","type":"PVZ","owner_code":"CDEK",
		"take_only":false,"is_handout":true,"is_reception":true,"is_dressing_room":true,"is_ltl":false,
		"have_cashless":true,"have_cash":true,"have_fast_payment_system":true,"allowed_cod":true,
		"site":"https://www.cdek.ru/nsk1",
		"office_image_list":[{"url":"https://img.cdek.ru/1.jpg"},{"url":"https://img.cdek.ru/2.jpg"}],
		"work_time_list":[{"day":1,"time":"08:00/21:00"},{"day":6,"time":"10:00/20:00"}],
		"work_time_exception_list":[{"date_start":"2026-01-01","date_end":"2026-01-01","time_start":"00:00","time_end":"00:00","is_working":false}],
		"weight_min":0.0,"weight_max":100000.0,"status":"ACTIVE",
		"location":{"country_code":"RU","region_code":23,"region":"Новосибирская область","city_code":270,
			"city":"Новосибирск","fias_guid":"fias-guid-000","postal_code":"630007","longitude":82.929124,
			"latitude":55.01591,"address":"ул. Кривощековская, 15, корп.1",
			"address_full":"630007, Россия, Новосибирск, ул. Кривощековская, 15, корп.1","city_uuid":"city-uuid-000"},
		"ltl_acceptance_partners":false,"ltl_issuance_partners":false,"fulfillment":false
	}]`)

	m := newDtoMapper()
	points, err := m.fromCDEKDeliveryPoints(raw)
	require.NoError(t, err)
	require.Len(t, points, 1)

	p := points[0]
	assert.Equal(t, "NSK1", p.Code)
	require.NotNil(t, p.UUID)
	assert.Equal(t, "33258e52-0000-0000-0000-000000000000", *p.UUID)
	require.NotNil(t, p.NearestStation)
	assert.Equal(t, "Автовокзал", *p.NearestStation)
	require.NotNil(t, p.OwnerCode)
	assert.Equal(t, "CDEK", *p.OwnerCode)
	assert.False(t, p.TakeOnly)
	assert.True(t, p.IsHandout)
	assert.True(t, p.IsReception)
	assert.True(t, p.IsDressingRoom)
	assert.False(t, p.IsLtl)
	assert.True(t, p.HaveCashless)
	assert.True(t, p.HaveCash)
	assert.True(t, p.HaveFastPaymentSystem)
	assert.True(t, p.AllowedCod)
	require.NotNil(t, p.Site)
	assert.Equal(t, "https://www.cdek.ru/nsk1", *p.Site)
	require.NotNil(t, p.WeightMin)
	assert.Equal(t, 0.0, *p.WeightMin)
	require.NotNil(t, p.WeightMax)
	assert.Equal(t, 100000.0, *p.WeightMax)
	assert.Equal(t, "ACTIVE", p.Status)

	require.NotNil(t, p.OfficeImage)
	assert.Equal(t, "https://img.cdek.ru/1.jpg", *p.OfficeImage)
	require.Len(t, p.OfficeImageList, 2)
	assert.Equal(t, []string{"https://img.cdek.ru/1.jpg", "https://img.cdek.ru/2.jpg"}, p.OfficeImageList)

	require.Len(t, p.WorkTimeList, 2)
	assert.Equal(t, WorkTimeEntry{Day: 1, Time: "08:00/21:00"}, p.WorkTimeList[0])
	assert.Equal(t, WorkTimeEntry{Day: 6, Time: "10:00/20:00"}, p.WorkTimeList[1])

	require.Len(t, p.WorkTimeExceptionList, 1)
	assert.Equal(t, WorkTimeException{
		DateStart: "2026-01-01",
		DateEnd:   "2026-01-01",
		TimeStart: "00:00",
		TimeEnd:   "00:00",
		IsWorking: false,
	}, p.WorkTimeExceptionList[0])

	loc := p.Location
	require.NotNil(t, loc.CountryCode)
	assert.Equal(t, "RU", *loc.CountryCode)
	require.NotNil(t, loc.RegionCode)
	assert.EqualValues(t, 23, *loc.RegionCode)
	assert.Equal(t, "Новосибирская область", loc.Region)
	require.NotNil(t, loc.CityCode)
	assert.EqualValues(t, 270, *loc.CityCode)
	assert.Equal(t, "Новосибирск", loc.City)
	require.NotNil(t, loc.FiasGUID)
	assert.Equal(t, "fias-guid-000", *loc.FiasGUID)
	require.NotNil(t, loc.CityUUID)
	assert.Equal(t, "city-uuid-000", *loc.CityUUID)
	require.NotNil(t, loc.AddressFull)
	assert.Equal(t, "630007, Россия, Новосибирск, ул. Кривощековская, 15, корп.1", *loc.AddressFull)
	assert.Equal(t, "ул. Кривощековская, 15, корп.1", loc.Address)
	assert.Equal(t, "630007", loc.PostalCode)
	assert.InDelta(t, 55.01591, loc.Latitude, 0.00001)
	assert.InDelta(t, 82.929124, loc.Longitude, 0.00001)
}
