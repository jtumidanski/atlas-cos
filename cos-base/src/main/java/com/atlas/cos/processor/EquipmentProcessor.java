package com.atlas.cos.processor;

import java.util.Arrays;
import java.util.Collection;
import java.util.Optional;
import java.util.stream.Stream;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.administrator.EquipmentAdministrator;
import com.atlas.cos.database.provider.EquipmentProvider;
import com.atlas.cos.model.EquipmentData;
import com.atlas.eso.attribute.EquipmentAttributes;
import com.atlas.eso.builder.EquipmentAttributesBuilder;
import com.atlas.eso.constant.RestConstants;
import com.atlas.iis.attribute.EquipmentSlotAttributes;
import com.atlas.shared.rest.UriBuilder;

import builder.ResultObjectBuilder;
import database.Connection;
import rest.DataBody;
import rest.DataContainer;

public final class EquipmentProcessor {
   private EquipmentProcessor() {
   }

   public static Optional<EquipmentData> createEquipmentForCharacter(int characterId, int itemId, boolean characterCreation) {
      if (characterCreation) {
         boolean valid = validCharacterCreationItem(itemId);
         if (!valid) {
            return Optional.empty();
         }
      }

      short nextOpenSlot = Connection.instance()
            .element(entityManager -> EquipmentProvider.getNextFreeEquipmentSlot(entityManager, characterId))
            .orElse((short) 0);

      return UriBuilder.service(RestConstants.SERVICE)
            .path("equipment")
            .getRestClient(EquipmentAttributes.class)
            .createWithResponse(new ResultObjectBuilder(EquipmentAttributes.class, 0)
                  .setAttribute(new EquipmentAttributesBuilder().setItemId(itemId))
                  .inputObject()
            )
            .result()
            .flatMap(DataContainer::data)
            .map(DataBody::getId)
            .map(Integer::parseInt)
            .flatMap(id -> Connection.instance().element(entityManager ->
                  EquipmentAdministrator.create(entityManager, characterId, id, nextOpenSlot)));
   }

   protected static boolean validCharacterCreationItem(int itemId) {
      return Stream.of(
            1302000, 1312004, 1322005, 1442079,// weapons
            1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011, 1042167,// bottom
            1060002, 1060006, 1061002, 1061008, 1062115, // top
            1072001, 1072005, 1072037, 1072038, 1072383,// shoes
            30000, 30010, 30020, 30030, 31000, 31040, 31050,// hair
            20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100  //face
      ).anyMatch(id -> id == itemId);
   }

   protected static Stream<Short> getEquipmentSlotDestination(int itemId) {
      return UriBuilder.service(com.atlas.iis.constant.RestConstants.SERVICE)
            .pathParam("equipment", itemId)
            .path("slots")
            .getRestClient(EquipmentSlotAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::dataList)
            .stream()
            .flatMap(Collection::stream)
            .map(DataBody::getAttributes)
            .map(EquipmentSlotAttributes::slot);
   }

   public static void equipItemForCharacter(int characterId, int equipmentId) {
      Connection.instance()
            .element(entityManager -> EquipmentProvider.getByEquipmentId(entityManager, equipmentId))
            .map(EquipmentData::equipmentId)
            .flatMap(EquipmentProcessor::getItemIdForEquipment)
            .flatMap(itemId -> getEquipmentSlotDestination(itemId).findFirst())
            .ifPresent(destinationSlot -> equipItemForCharacter(characterId, equipmentId, destinationSlot));
   }

   protected static Optional<Integer> getItemIdForEquipment(int equipmentId) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("equipment", equipmentId)
            .getRestClient(EquipmentAttributes.class)
            .getWithResponse()
            .result()
            .flatMap(DataContainer::data)
            .map(DataBody::getAttributes)
            .map(EquipmentAttributes::itemId);
   }

   protected static void equipItemForCharacter(int characterId, int equipmentId, short destinationSlot) {
      Connection.instance().with(entityManager ->
            QueryAdministratorUtil.inTransaction(entityManager, transactionEm -> {
               short temporarySlot = Short.MIN_VALUE;

               EquipmentProvider.findEquipmentInSlot(transactionEm, characterId, destinationSlot)
                     .ifPresent(itemToMove -> EquipmentAdministrator.updateSlot(transactionEm, itemToMove, temporarySlot));

               short currentSlot = EquipmentProvider
                     .getByEquipmentId(transactionEm, equipmentId)
                     .map(EquipmentData::slot)
                     .orElse(EquipmentProvider.getNextFreeEquipmentSlot(transactionEm, characterId).orElse((short) 0));

               EquipmentAdministrator.updateSlot(transactionEm, equipmentId, destinationSlot);

               EquipmentProvider.findEquipmentInSlot(transactionEm, characterId, temporarySlot)
                     .ifPresent(itemToMove -> EquipmentAdministrator.updateSlot(transactionEm, itemToMove, currentSlot));
            }));
   }

   public static Optional<EquipmentData> getEquipmentById(int id) {
      return Connection.instance().element(entityManager -> EquipmentProvider.getById(entityManager, id));
   }

   public static Stream<EquipmentData> getEquipmentForCharacter(int characterId) {
      return Connection.instance()
            .list(entityManager -> EquipmentProvider.getForCharacter(entityManager, characterId))
            .stream();
   }

   public static Optional<EquipmentData> getEquippedItemForCharacterBySlot(int characterId, short slotId) {
      return Connection.instance().list(entityManager -> EquipmentProvider.getForCharacter(entityManager, characterId)).stream()
            .filter(equipmentData -> equipmentData.slot() == slotId)
            .findFirst();
   }

   public static void createAndEquip(int characterId, int... items) {
      Arrays.stream(items).forEach(itemId -> createAndEquip(characterId, itemId));
   }

   public static void createAndEquip(int characterId, int itemId) {
      createEquipmentForCharacter(characterId, itemId, true)
            .ifPresent(equipment -> equipItemForCharacter(characterId, equipment.equipmentId()));
   }
}
