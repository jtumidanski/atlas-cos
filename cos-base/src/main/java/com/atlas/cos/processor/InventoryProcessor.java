package com.atlas.cos.processor;

import java.util.Arrays;
import java.util.List;
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
      List<InventoryItem> items;
      if (type.equals(InventoryType.EQUIP)) {
         items = getEquipInventoryItems(characterId);
      } else {
         items = getInventoryItems(characterId, type);
      }
      return Optional.of(new Inventory(type.getType(), type.name(), 4, items));
   }

   protected static List<InventoryItem> getEquipInventoryItems(int characterId) {
      return ItemProcessor.getEquipmentForCharacter(characterId).stream()
            .map(equipment -> new InventoryItem(equipment.id(), InventoryItemType.EQUIPMENT, equipment.slot()))
            .collect(Collectors.toList());
   }

   protected static List<InventoryItem> getInventoryItems(int characterId, InventoryType type) {
      return ItemProcessor.getItemsForCharacter(characterId, type).stream()
            .map(item -> new InventoryItem(item.id(), InventoryItemType.ITEM, item.slot()))
            .collect(Collectors.toList());
   }
}
