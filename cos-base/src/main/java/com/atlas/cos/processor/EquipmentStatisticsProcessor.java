package com.atlas.cos.processor;

import java.util.Optional;

import com.atlas.cos.model.EquipmentStatistics;
import com.atlas.eso.attribute.EquipmentAttributes;
import com.atlas.eso.constant.RestConstants;
import com.atlas.shared.rest.UriBuilder;

import rest.DataContainer;

public final class EquipmentStatisticsProcessor {
   private EquipmentStatisticsProcessor() {
   }

   public static Optional<EquipmentStatistics> getEquipData(int equipmentId) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("equipment", equipmentId)
            .getRestClient(EquipmentAttributes.class)
            .getWithResponse()
            .result()
            .flatMap(DataContainer::data)
            .map(ModelFactory::createEquip);
   }
}
