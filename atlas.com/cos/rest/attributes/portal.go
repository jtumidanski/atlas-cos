package attributes

type PortalDataContainer struct {
   data     dataSegment
   included dataSegment
}

type PortalData struct {
   Id         string           `json:"id"`
   Type       string           `json:"type"`
   Attributes PortalAttributes `json:"attributes"`
}

type PortalAttributes struct {
   Name       string `json:"name"`
   Target     string `json:"target"`
   Type       uint32 `json:"type"`
   X          int32 `json:"x"`
   Y          int32 `json:"y"`
   TargetMap  uint32 `json:"targetMap"`
   ScriptName string `json:"scriptName"`
}

func (a *PortalDataContainer) UnmarshalJSON(data []byte) error {
   d, i, err := unmarshalRoot(data, mapperFunc(EmptyPortalData))
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
