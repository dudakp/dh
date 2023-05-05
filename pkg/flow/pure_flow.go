package flow

//
//func ExecutePureFlow(flow *flow) (error, Data) {
//	return executePure(flow.firstHandler, flow.InitialData, flow)
//}
//
//func executePure(handler *Handler, data Data, flow *flow) (error, Data) {
//	err, resData := handler.Action(data)
//	if err != nil {
//		return handleError(handler, err, flow), nil
//	} else {
//		return executePure(handler.next, resData, flow)
//	}
//}
