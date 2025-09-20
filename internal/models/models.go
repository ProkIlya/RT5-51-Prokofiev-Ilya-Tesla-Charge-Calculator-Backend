package models

type DrivingScenarioType string

const (
	RoadType    DrivingScenarioType = "road"
	ComfortType DrivingScenarioType = "comfort"
)

type DrivingScenario struct {
	ID               int                 `json:"id"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	SystemConsuption float64             `json:"systemConsuption"`
	Speed            float64             `json:"speed"`
	AeroCoeff        float64             `json:"aeroCoeff"`
	RollingCoeff     float64             `json:"rollingCoeff"`
	ImageURL         string              `json:"imageURL"`
	Type             DrivingScenarioType `json:"type"`
}

type TripScenario struct {
	Scenario *DrivingScenario `json:"scenario"`
	Value    float64          `json:"value"`
}

type TripCalculationCharge struct {
	ID              int            `json:"id"`
	StartCharge     float64        `json:"startCharge"`
	RemainingCharge float64        `json:"remainingCharge"`
	Scenarios       []TripScenario `json:"scenarios"`
}

var DrivingScenarios = []*DrivingScenario{
	{
		ID:               1,
		Name:             "Экономичная городская поездка",
		Description:      "Поездка по городу без каких-либо систем комфорта",
		SystemConsuption: 0,
		Speed:            50,
		AeroCoeff:        0.208,
		RollingCoeff:     0.012,
		ImageURL:         "http://localhost:9000/tesla-images/1.jpg",
		Type:             RoadType,
	},
	{
		ID:               2,
		Name:             "Кондиционер",
		Description:      "Работа системы кондиционирования увеличивается потребление энергии",
		SystemConsuption: 1.5,
		Speed:            0,
		AeroCoeff:        0,
		RollingCoeff:     0,
		ImageURL:         "http://localhost:9000/tesla-images/2.jpg",
		Type:             ComfortType,
	},
	{
		ID:               3,
		Name:             "Зимний город",
		Description:      "Заснежененная трасса и зимняя резина увеличивает трение и ухудшает аэродинамику",
		SystemConsuption: 3.0,
		Speed:            50,
		AeroCoeff:        0.228,
		RollingCoeff:     0.018,
		ImageURL:         "http://localhost:9000/tesla-images/3.jpg",
		Type:             RoadType,
	},
	{
		ID:               4,
		Name:             "Загородная трасса",
		Description:      "Вождение по шоссе на высоких скоростях",
		SystemConsuption: 0,
		Speed:            110,
		AeroCoeff:        0.208,
		RollingCoeff:     0.012,
		ImageURL:         "http://localhost:9000/tesla-images/4.jpg",
		Type:             RoadType,
	},
	{
		ID:               5,
		Name:             "Зимняя загородная трасса",
		Description:      "Высокая скорость, заснеженная трасса и потеря аэродинамики",
		SystemConsuption: 2.0,
		Speed:            120,
		AeroCoeff:        0.228,
		RollingCoeff:     0.018,
		ImageURL:         "http://localhost:9000/tesla-images/5.jpg",
		Type:             RoadType,
	},
}

var Trips = map[int]*TripCalculationCharge{
	1: {
		ID:              1,
		StartCharge:     100,
		RemainingCharge: 92.82, // Предварительно рассчитанное значение
		Scenarios: []TripScenario{
			{Scenario: DrivingScenarios[0], Value: 50},
			{Scenario: DrivingScenarios[1], Value: 2},
		},
	},
}
