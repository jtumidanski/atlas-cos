package attributes

import (
	"atlas-cos/rest/response"
	"encoding/json"
)

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

func (c *PortalDataContainer) MarshalJSON() ([]byte, error) {
	t := struct {
		Data     interface{} `json:"data"`
		Included interface{} `json:"included"`
	}{}
	if len(c.data) == 1 {
		t.Data = c.data[0]
	} else {
		t.Data = c.data
	}
	return json.Marshal(t)
}

func (c *PortalDataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyPortalData))
	if err != nil {
		return err
	}

	c.data = d
	c.included = i
	return nil
}

func (c *PortalDataContainer) Data() *PortalData {
	if len(c.data) >= 1 {
		return c.data[0].(*PortalData)
	}
	return nil
}

func (c *PortalDataContainer) DataList() []PortalData {
	var r = make([]PortalData, 0)
	for _, x := range c.data {
		r = append(r, *x.(*PortalData))
	}
	return r
}

func EmptyPortalData() interface{} {
	return &PortalData{}
}
