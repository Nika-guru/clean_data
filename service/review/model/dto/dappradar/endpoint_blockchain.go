package dto_dappradar

import (
	"review-service/service/review/model/dao"
)

type EndpointBlockchainRepo struct {
	EndpointBlockchains []*EndpointBlockchain
}

type EndpointBlockchain struct {
	Id                 *uint64
	Endpoint           string
	BlockchainId       string
	BlockchainName     string
	BlockchainImageSvg string
}

func (dtoRepo *EndpointBlockchainRepo) ConvertTo(daoRepo *dao.BlockchainRepo) {
	for _, dto := range dtoRepo.EndpointBlockchains {
		dao := &dao.Blockchain{}
		dto.ConvertTo(dao)
		daoRepo.Blockchains = append(daoRepo.Blockchains, dao)
	}
}

func (dto *EndpointBlockchain) ConvertTo(dao *dao.Blockchain) {
	dto.Id = &dao.Id
	dao.BlockchainName = dto.BlockchainName
	dao.BlockchainId = dto.BlockchainId
	if dao.Info == nil {
		dao.Info = make(map[string]any)
	}
	dao.Info[`svgHtml`] = dto.BlockchainImageSvg
}
