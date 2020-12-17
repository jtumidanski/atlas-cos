package com.atlas.cos.processor;

import java.util.Optional;

import com.atlas.cos.model.EquipData;
import com.atlas.eso.attribute.EquipmentAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import rest.DataContainer;

public final class EquipProcessor {
   private EquipProcessor() {
   }

   public static Optional<EquipData> getEquipData(int equipmentId) {
      return UriBuilder.service(RestService.EQUIPMENT_STORAGE)
            .pathParam("equipment", equipmentId)
            .getRestClient(EquipmentAttributes.class)
            .getWithResponse()
            .result()
            .flatMap(DataContainer::data)
            .map(ModelFactory::createEquip);
   }
}
