package com.atlas.cos.processor;

import java.util.Arrays;
import java.util.Optional;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import com.atlas.cos.model.Inventory;
import com.atlas.cos.model.InventoryItem;
import com.atlas.cos.model.InventoryItemType;
import com.atlas.cos.model.InventoryType;

public final class InventoryProcessor {
   private InventoryProcessor() {
   }

   public static Stream<Inventory> getAllInventories(int characterId) {
      return Arrays.stream(InventoryType.values())
            .map(type -> InventoryProcessor.getInventoryByType(characterId, type))
            .flatMap(Optional::stream);
   }

   public static Optional<Inventory> getInventoryByType(int characterId, String type) {
      return InventoryType.getByName(type)
            .flatMap(inventoryType -> InventoryProcessor.getInventoryByType(characterId, inventoryType));
   }

   protected static Optional<Inventory> getInventoryByType(int characterId, InventoryType type) {
      Stream<InventoryItem> items;
      if (type.equals(InventoryType.EQUIP)) {
         items = getEquipInventoryItems(characterId);
      } else {
         items = getInventoryItems(characterId, type);
      }
      return Optional.of(new Inventory(type.getType(), type.name(), 4, items.collect(Collectors.toList())));
   }

   protected static Stream<InventoryItem> getEquipInventoryItems(int characterId) {
      return ItemProcessor.getEquipmentForCharacter(characterId)
            .map(equipment -> new InventoryItem(equipment.id(), InventoryItemType.EQUIPMENT, equipment.slot()));
   }

   protected static Stream<InventoryItem> getInventoryItems(int characterId, InventoryType type) {
      return ItemProcessor.getItemsForCharacter(characterId, type)
            .map(item -> new InventoryItem(item.id(), InventoryItemType.ITEM, item.slot()));
   }
}
