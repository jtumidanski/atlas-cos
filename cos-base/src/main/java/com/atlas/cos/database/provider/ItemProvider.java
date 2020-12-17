package com.atlas.cos.database.provider;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import javax.persistence.EntityManager;
import javax.persistence.TypedQuery;

import com.app.database.util.QueryProviderUtil;
import com.atlas.cos.database.transformer.ItemDataTransformer;
import com.atlas.cos.entity.Item;
import com.atlas.cos.model.ItemData;

public final class ItemProvider {
   private ItemProvider() {
   }

   public static List<ItemData> getItemsForCharacter(EntityManager entityManager, int characterId, byte type, int itemId) {
      TypedQuery<Item> query = entityManager.createQuery("SELECT i FROM Item i WHERE i.characterId = :characterId AND i"
            + ".inventoryType = :type AND i.itemId = :itemId", Item.class);
      query.setParameter("characterId", characterId);
      query.setParameter("type", type);
      query.setParameter("itemId", itemId);
      return QueryProviderUtil.list(query, new ItemDataTransformer());
   }

   public static Optional<Short> getNextFreeEquipmentSlot(EntityManager entityManager, int characterId, byte type) {
      Set<Short> usedSlots = getForCharacterByInventory(entityManager, characterId, type).stream()
            .map(ItemData::slot)
            .filter(slot -> slot >= 0)
            .collect(Collectors.toSet());
      if (usedSlots.isEmpty()) {
         return Optional.of((short) 0);
      }

      Set<Short> slots = IntStream.rangeClosed(Collections.min(usedSlots), Collections.max(usedSlots))
            .boxed()
            .map(Integer::shortValue)
            .collect(Collectors.toSet());

      slots.removeAll(usedSlots);

      return Optional.of(slots.stream()
            .findFirst()
            .orElse((short) usedSlots.size()));
   }

   private static List<ItemData> getForCharacterByInventory(EntityManager entityManager, int characterId, byte type) {
      TypedQuery<Item> query = entityManager.createQuery("SELECT i FROM Item i WHERE i.characterId = :characterId AND i"
            + ".inventoryType = :type", Item.class);
      query.setParameter("characterId", characterId);
      query.setParameter("type", type);
      return QueryProviderUtil.list(query, new ItemDataTransformer());
   }
}
