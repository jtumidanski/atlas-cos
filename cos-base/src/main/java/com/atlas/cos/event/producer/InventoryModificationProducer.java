package com.atlas.cos.event.producer;

import java.util.Collections;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterInventoryModifyEvent;
import com.atlas.cos.event.InventoryModification;

public final class InventoryModificationProducer {
   private InventoryModificationProducer() {
   }

   public static void emit(int characterId, int mode, int itemId, int quantity, int inventoryType, int position) {
      EventProducerRegistry.getInstance()
            .send(CharacterInventoryModifyEvent.class, EventConstants.TOPIC_INVENTORY_MODIFICATION, characterId,
                  new CharacterInventoryModifyEvent(characterId, true,
                        Collections.singletonList(new InventoryModification(mode, itemId, quantity, inventoryType, position))));
   }
}
