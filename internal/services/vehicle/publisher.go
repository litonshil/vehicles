package vehicle

import (
	"log"

	"vehicles/infra/conn"
)

type VehicleApprovalMessage struct {
	VehicleID string `json:"vehicle_id"`
	DriverID  string `json:"driver_id"`
}

// PublishVehicleApproval Publish a vehicle approval message
func PublishVehicleApproval(vehicleID, driverID string) {
	message := VehicleApprovalMessage{
		VehicleID: vehicleID,
		DriverID:  driverID,
	}

	err := conn.PublishMessage("vehicle_exchange", "vehicle.approved", message)
	if err != nil {
		log.Printf("Error publishing vehicle approval message: %v", err)
	} else {
		log.Println("Vehicle approval message published successfully")
	}
}
