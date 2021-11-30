package party

import (
	"atlas-cos/rest/response"
	"encoding/json"
)

type dataContainer struct {
	data     response.DataSegment
	included response.DataSegment
}

type dataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
	LeaderId uint32 `json:"leader_id"`
}

func (c *dataContainer) MarshalJSON() ([]byte, error) {
	t := struct {
		Data     interface{} `json:"data"`
		Included interface{} `json:"included"`
	}{}
	if len(c.data) == 1 {
		t.Data = c.data[0]
	} else {
		t.Data = c.data
	}
	if len(c.included) == 1 {
		t.Included = c.included[0]
	} else {
		t.Included = c.included
	}
	return json.Marshal(t)
}

func (c *dataContainer) UnmarshalJSON(data []byte) error {
	d, _, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyData))
	if err != nil {
		return err
	}
	c.data = d
	return nil
}

func (c *dataContainer) Data() *dataBody {
	if len(c.data) >= 1 {
		return c.data[0].(*dataBody)
	}
	return nil
}

func (c *dataContainer) DataList() []dataBody {
	var r = make([]dataBody, 0)
	for _, x := range c.data {
		r = append(r, *x.(*dataBody))
	}
	return r
}

func EmptyData() interface{} {
	return &dataBody{}
}
