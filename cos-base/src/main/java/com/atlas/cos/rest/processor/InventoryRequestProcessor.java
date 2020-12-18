package com.atlas.cos.rest.processor;

import java.util.Arrays;
import java.util.Collection;
import java.util.List;
import java.util.Optional;

import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.model.Inventory;
import com.atlas.cos.processor.EquipmentProcessor;
import com.atlas.cos.processor.EquipmentStatisticsProcessor;
import com.atlas.cos.processor.InventoryProcessor;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;

public final class InventoryRequestProcessor {
   protected static List<String> getIncludedResources(String include) {
      return Arrays.asList(include.split(",").clone());
   }

   protected static boolean containsIgnoreCase(List<String> list, String search) {
      return list.stream().anyMatch(item -> item.equalsIgnoreCase(search));
   }

   public static ResultBuilder getAllInventories(int characterId, String include) {
      List<String> includedResources = getIncludedResources(include);
      ResultBuilder resultBuilder = getAllInventories(characterId);

      if (containsIgnoreCase(includedResources, "inventoryItems")) {
         InventoryProcessor.getAllInventories(characterId)
               .map(Inventory::items)
               .flatMap(Collection::stream)
               .map(ResultObjectFactory::create)
               .forEach(resultBuilder::addInclude);
      }
      if (containsIgnoreCase(includedResources, "equipmentStatistics")) {
         EquipmentProcessor.getEquipmentForCharacter(characterId)
               .map(EquipmentData::equipmentId)
               .map(EquipmentStatisticsProcessor::getEquipData)
               .flatMap(Optional::stream)
               .map(ResultObjectFactory::create)
               .forEach(resultBuilder::addInclude);
      }

      return resultBuilder;
   }

   protected static ResultBuilder getAllInventories(int characterId) {
      return InventoryProcessor.getAllInventories(characterId)
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public static ResultBuilder getInventoryByType(int characterId, String type, String include) {
      List<String> includedResources = getIncludedResources(include);
      ResultBuilder resultBuilder = getInventoryByType(characterId, type);

      if (containsIgnoreCase(includedResources, "inventoryItems")) {
         InventoryProcessor.getInventoryByType(characterId, type)
               .map(Inventory::items)
               .stream()
               .flatMap(Collection::stream)
               .map(ResultObjectFactory::create)
               .forEach(resultBuilder::addInclude);
      }
      if (containsIgnoreCase(includedResources, "equipmentStatistics")) {
         EquipmentProcessor.getEquipmentForCharacter(characterId)
               .map(EquipmentData::equipmentId)
               .map(EquipmentStatisticsProcessor::getEquipData)
               .flatMap(Optional::stream)
               .map(ResultObjectFactory::create)
               .forEach(resultBuilder::addInclude);
      }

      return resultBuilder;
   }

   protected static ResultBuilder getInventoryByType(int characterId, String type) {
      return InventoryProcessor.getInventoryByType(characterId, type)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleOkResult)
            .orElseGet(ResultBuilder::notFound);
   }
}
