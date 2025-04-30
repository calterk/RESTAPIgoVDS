package service

import (
	"app/engine"
	"app/entity"
	"app/storage"
	"encoding/json"
)

func (e *Service) CreateVds(ctx *engine.Context) error { //POST
	decoder := json.NewDecoder(ctx.Request.Body)
	var vds entity.Vds
	err := decoder.Decode(&vds)
	if err != nil {
		return err
	}
	vds = storage.CreateVds(vds)
	coder, _ := json.Marshal(vds)
	ctx.Response.Write(coder)
	return nil
}

func (e *Service) UpdateVds(ctx *engine.Context) error { //PUT

	return nil
}

func (e *Service) DeleteVds(ctx *engine.Context) error { //DELETE

	return nil
}

func (e *Service) GetAll(ctx *engine.Context) error { //GET
	//storage.GetAllVds(ctx)
	return nil
}
func (e *Service) GetById(ctx *engine.Context) error { //GET

	return nil
}
