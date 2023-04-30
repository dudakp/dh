package flow

func ExecutePureFlow(flow *Flow) (error, HandlerData) {
	return executePure(flow.firstHandler, flow.InitialData, flow)
}

func executePure(handler *Handler, data HandlerData, flow *Flow) (error, HandlerData) {
	err, resData := handler.Action(data)
	if err != nil {
		return handleError(handler, err, flow), nil
	} else {
		return executePure(handler.next, resData, flow)
	}
}
