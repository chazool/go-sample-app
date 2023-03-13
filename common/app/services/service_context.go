package services

import "context"

// ServiceContext used to define service interface
type ServiceContext struct {
	_           struct{}
	RequestID   string
	ParentCtx   context.Context
	RequestBody string
}

//CreateServiceContext use to create service context
func CreateServiceContext(parentCtx context.Context, requestID string, requestBody string) ServiceContext {
	return ServiceContext{RequestID: requestID, RequestBody: requestBody, ParentCtx: parentCtx}
}
