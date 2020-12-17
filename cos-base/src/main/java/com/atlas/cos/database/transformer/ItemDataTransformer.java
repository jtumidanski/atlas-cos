package com.atlas.cos.database.transformer;

import com.atlas.cos.entity.Item;
import com.atlas.cos.model.ItemData;

import transformer.SqlTransformer;

public class ItemDataTransformer implements SqlTransformer<ItemData, Item> {
   @Override
   public ItemData transform(Item item) {
      return new ItemData(item.getId(), item.getCharacterId(), item.getInventoryType(), item.getItemId(), item.getQuantity(),
            item.getSlot());
   }
}
