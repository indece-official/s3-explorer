package web

import (
	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

func (c *Controller) mapObjectV1ToAPIObjectV1(object *model.ObjectV1) webapi.ObjectV1 {
	apiObject := webapi.ObjectV1{}

	apiObject.Key = object.Key
	apiObject.LastModified = object.LastModified
	apiObject.OwnerName = object.OwnerName
	apiObject.OwnerId = object.OwnerID
	apiObject.Size = object.Size

	return apiObject
}
