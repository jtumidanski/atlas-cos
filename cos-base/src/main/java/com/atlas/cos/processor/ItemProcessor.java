package com.atlas.cos.processor;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.stream.Stream;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.administrator.EquipmentAdministrator;
import com.atlas.cos.database.administrator.ItemAdministrator;
import com.atlas.cos.database.provider.EquipmentProvider;
import com.atlas.cos.database.provider.ItemProvider;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.model.InventoryType;
import com.atlas.cos.model.ItemData;
import com.atlas.eso.attribute.EquipmentAttributes;
import com.atlas.eso.builder.EquipmentAttributesBuilder;
import com.atlas.iis.attribute.EquipmentSlotAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import builder.ResultObjectBuilder;
import database.Connection;
import rest.DataBody;
import rest.DataContainer;

public final class ItemProcessor {
   private ItemProcessor() {
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

      return UriBuilder.service(RestService.EQUIPMENT_STORAGE)
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

   protected static Stream<Short> getEquipmentSlotDestination(int itemId) {
      return UriBuilder.service(RestService.ITEM_INFORMATION)
            .pathParam("equipment", itemId)
            .path("slots")
            .getRestClient(EquipmentSlotAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::dataList)
            .orElse(Collections.emptyList()).stream()
            .map(DataBody::getAttributes)
            .map(EquipmentSlotAttributes::slot);
   }

   public static void equipItemForCharacter(int characterId, int equipmentId) {
      Connection.instance()
            .element(entityManager -> EquipmentProvider.getByEquipmentId(entityManager, equipmentId))
            .map(EquipmentData::equipmentId)
            .flatMap(ItemProcessor::getItemIdForEquipment)
            .flatMap(itemId -> getEquipmentSlotDestination(itemId).findFirst())
            .ifPresent(destinationSlot -> equipItemForCharacter(characterId, equipmentId, destinationSlot));
   }

   protected static Optional<Integer> getItemIdForEquipment(int equipmentId) {
      return UriBuilder.service(RestService.EQUIPMENT_STORAGE)
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

   public static List<EquipmentData> getEquipmentForCharacter(int characterId) {
      return Connection.instance().list(entityManager -> EquipmentProvider.getForCharacter(entityManager, characterId));
   }

   public static Optional<EquipmentData> getEquippedItemForCharacterBySlot(int characterId, short slotId) {
      return Connection.instance().list(entityManager -> EquipmentProvider.getForCharacter(entityManager, characterId)).stream()
            .filter(equipmentData -> equipmentData.slot() == slotId)
            .findFirst();
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

   public static List<ItemData> getItemsForCharacter(int characterId, InventoryType inventoryType) {
      return Connection.instance()
            .list(entityManager -> ItemProvider.getForCharacterByInventory(entityManager, characterId, inventoryType.getType()));
   }

   public static List<ItemData> getItemsForCharacter(int characterId, InventoryType inventoryType, int itemId) {
      return Connection.instance()
            .list(entityManager -> ItemProvider.getItemsForCharacter(entityManager, characterId, inventoryType.getType(), itemId));
   }

   public static void updateItemQuantity(int uniqueId, int quantity) {
      Connection.instance().with(entityManager -> ItemAdministrator.updateQuantity(entityManager, uniqueId, quantity));
   }

   public static void createItemForCharacter(int characterId, InventoryType inventoryType, int itemId, int quantity) {
      short nextOpenSlot = Connection.instance()
            .element(entityManager -> ItemProvider.getNextFreeEquipmentSlot(entityManager, characterId, inventoryType.getType()))
            .orElse((short) 0);

      Connection.instance().element(entityManager ->
            ItemAdministrator.create(entityManager, characterId, inventoryType.getType(), itemId, quantity, nextOpenSlot));
   }

   public static Optional<ItemData> getItemById(int id) {
      return Connection.instance().element(entityManager -> ItemProvider.getById(entityManager, id));
   }
}
