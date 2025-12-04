package calculations

import "tesla-app/internal/app/ds"

const (
	MassKg         = 2068
	AirDensityKgM3 = 1.225
	FrontalAreaM2  = 2.34
	GravityMpS2    = 9.8
)

func CalculateRemainingCharge(trip *ds.TripApplication, tripScenarios []ds.TripScenario) float64 {
	if trip.Status != "завершён" {
		return 0 // Расчет только для завершенных заявок
	}

	// Проверяем что StartCharge не nil
	if trip.StartCharge == nil {
		return 0
	}

	totalEnergy := 0.0

	for _, tripScenario := range tripScenarios {
		scenario := tripScenario.DrivingScenario

		// Проверяем что Duration не nil, иначе пропускаем
		if tripScenario.Duration == nil {
			continue
		}

		value := *tripScenario.Duration // Разыменовываем указатель

		if scenario.Type == "дорога" {
			// Расчет энергии на движение для дорожных условий
			F_rolling := float64(MassKg) * GravityMpS2 * scenario.RollingCoeff
			speedMs := scenario.Speed / 3.6
			F_air := 0.5 * AirDensityKgM3 * scenario.AeroCoeff * FrontalAreaM2 * speedMs * speedMs
			E_drive := ((F_rolling + F_air) * value * 1000) / 3600000
			totalEnergy += E_drive
		} else if scenario.Type == "комфорт" {
			// Расчет энергии на системы комфорта
			E_systems := scenario.SystemConsuption * value
			totalEnergy += E_systems
		}
	}

	remainingCharge := *trip.StartCharge - totalEnergy // Разыменовываем указатель

	if remainingCharge < 0 {
		return 0
	}
	return remainingCharge
}
