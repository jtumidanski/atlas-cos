package com.atlas.cos.database.administrator;

import java.util.Optional;
import javax.persistence.EntityManager;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.transformer.ItemDataTransformer;
import com.atlas.cos.entity.Item;
import com.atlas.cos.model.ItemData;

public final class ItemAdministrator {
   private ItemAdministrator() {
   }

   public static void updateQuantity(EntityManager entityManager, int uniqueId, int quantity) {
      QueryAdministratorUtil.update(entityManager, Item.class, uniqueId, item -> item.setQuantity(quantity));
   }

   public static Optional<ItemData> create(EntityManager entityManager, int characterId, byte type, int itemId, int quantity,
                                           short slot) {
      Item item = new Item();
      item.setCharacterId(characterId);
      item.setInventoryType(type);
      item.setItemId(itemId);
      item.setQuantity(quantity);
      item.setSlot(slot);
      return Optional.of(QueryAdministratorUtil.insertAndReturn(entityManager, item, new ItemDataTransformer()));
   }
}
