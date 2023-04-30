package flow

func ExecuteEffectFlow(flow *Flow) error {
	return executeEffect(flow.firstHandler, flow)
}

func executeEffect(handler *Handler, flow *Flow) error {
	err, _ := handler.Action(flow.InitialData)
	if err != nil {
		return handleError(handler, err, flow)
	} else {
		return executeEffect(handler.next, flow)
	}
}
