package model

type Deployment struct {
	Id          string        `json:"id"`
	Executable  bool          `json:"executable"`
	XmlRaw      string        `json:"xml_raw"`
	Xml         string        `json:"xml"`
	Svg         string        `json:"svg"`
	Name        string        `json:"name"`
	Elements    []Element     `json:"elements"`
	Lanes       []LaneElement `json:"lanes"`
	Description string        `json:"description,omitempty"`
}

type Element struct {
	Order            int64      `json:"order"`
	Task             *Task      `json:"task,omitempty"`
	MultiTask        *MultiTask `json:"multi_task,omitempty"`
	ReceiveTaskEvent *MsgEvent  `json:"receive_task_event,omitempty"`
	MsgEvent         *MsgEvent  `json:"msg_event,omitempty"`
	TimeEvent        *TimeEvent `json:"time_event,omitempty"`
}

type LaneElement struct {
	Order     int64      `json:"order"`
	MultiLane *MultiLane `json:"multi_lane,omitempty"`
	Lane      *Lane      `json:"lane,omitempty"`
}
