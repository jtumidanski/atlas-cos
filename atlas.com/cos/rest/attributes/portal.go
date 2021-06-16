package attributes

import "atlas-cos/rest/response"

type PortalDataContainer struct {
	data     response.DataSegment
	included response.DataSegment
}

type PortalData struct {
	Id         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes PortalAttributes `json:"attributes"`
}

type PortalAttributes struct {
	Name        string `json:"name"`
	Target      string `json:"target"`
	Type        uint8  `json:"type"`
	X           int16  `json:"x"`
	Y           int16  `json:"y"`
	TargetMapId uint32 `json:"target_map_id"`
	ScriptName  string `json:"script_name"`
}

func (a *PortalDataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyPortalData))
	if err != nil {
		return err
	}

	a.data = d
	a.included = i
	return nil
}

func (a *PortalDataContainer) Data() *PortalData {
	if len(a.data) >= 1 {
		return a.data[0].(*PortalData)
	}
	return nil
}

func (a *PortalDataContainer) DataList() []PortalData {
	var r = make([]PortalData, 0)
	for _, x := range a.data {
		r = append(r, *x.(*PortalData))
	}
	return r
}

func EmptyPortalData() interface{} {
	return &PortalData{}
}
