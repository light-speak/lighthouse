package resolver

import (
	"search/kitex_gen/search/searchservice"

	"gorm.io/gorm"
)

type Resolver struct {
	Db           *gorm.DB
	SearchClient *searchservice.Client
}
