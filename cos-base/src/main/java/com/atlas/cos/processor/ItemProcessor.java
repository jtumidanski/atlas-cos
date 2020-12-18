package com.atlas.cos.processor;

import java.util.List;
import java.util.Optional;
import java.util.stream.Stream;

import com.atlas.cos.database.administrator.ItemAdministrator;
import com.atlas.cos.database.provider.ItemProvider;
import com.atlas.cos.model.InventoryType;
import com.atlas.cos.model.ItemData;

import database.Connection;

public final class ItemProcessor {
   private ItemProcessor() {
   }

   public static Stream<ItemData> getItemsForCharacter(int characterId, InventoryType inventoryType) {
      return Connection.instance()
            .list(entityManager -> ItemProvider.getForCharacterByInventory(entityManager, characterId, inventoryType.getType()))
            .stream();
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
