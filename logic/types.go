package logic

// TraceIDParam
// the traceIDParam use receive traceID in path
type TraceIDParam struct {
	TraceID string `path:"traceid"`
}

func (tid *TraceIDParam) Empty() bool {
	return tid.TraceID == ""
}

type ServiceNameParam struct {
	ServiceName string
}

func (svcN *ServiceNameParam) Empty() bool {
	return svcN.ServiceName == ""
}
