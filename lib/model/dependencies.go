package model

type Dependencies struct {
	//from db
	DeploymentId string             `json:"deployment_id" bson:"deployment_id"`
	Owner        string             `json:"owner" bson:"owner"`
	Devices      []DeviceDependency `json:"devices" bson:"devices"`
	Events       []EventDependency  `json:"events" bson:"events"`
}

type DeviceDependency struct {
	DeviceId      string         `json:"device_id" bson:"device_id"`
	Name          string         `json:"name" bson:"name"`
	BpmnResources []BpmnResource `json:"bpmn_resources" bson:"bpmn_resources"`
}

type EventDependency struct {
	EventId       string         `json:"event_id" bson:"event_id"`
	BpmnResources []BpmnResource `json:"bpmn_resources" bson:"bpmn_resources"`
}

type BpmnResource struct {
	Id    string `json:"id" bson:"id"`
	label string `json:"label" bson:"label"`
}
