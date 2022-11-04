package repository

import "github.com/google/uuid"

func GetIdOperativeSystem(operativeSystem string) (id uuid.UUID) {
	var idResp uuid.UUID
	for i := 0; i < len(OperativeSystems); i++ {
		if OperativeSystems[i].OperativeSystem == operativeSystem {
			idResp = OperativeSystems[i].Id
		}
	}
	return idResp
}

func GetOperativeSystem(id uuid.UUID) (code string) {
	var operativeSystemResp string
	for i := 0; i < len(OperativeSystems); i++ {
		if OperativeSystems[i].Id == id {
			operativeSystemResp = OperativeSystems[i].OperativeSystem
		}
	}
	return operativeSystemResp
}
