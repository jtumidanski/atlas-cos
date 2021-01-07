package com.atlas.cos.processor;

import java.util.Iterator;
import java.util.List;
import java.util.Optional;

import com.atlas.cos.event.producer.CancelDropReservationProducer;
import com.atlas.cos.event.producer.GainMesoProducer;
import com.atlas.cos.event.producer.InventoryModificationProducer;
import com.atlas.cos.event.producer.PickedUpItemProducer;
import com.atlas.cos.event.producer.PickedUpNxProducer;
import com.atlas.cos.event.producer.PickupDropCommandProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.Drop;
import com.atlas.cos.model.InventoryType;
import com.atlas.cos.model.ItemData;
import com.atlas.cos.model.Party;
import com.atlas.drg.constant.RestConstants;
import com.atlas.drg.rest.attribute.DropAttributes;
import com.atlas.shared.rest.UriBuilder;

import rest.DataContainer;

public final class DropProcessor {
   private DropProcessor() {
   }

   public static void attemptPickup(int characterId, int dropId) {
      CharacterProcessor.getById(characterId)
            .ifPresent(character -> getDrop(dropId)
                  .ifPresent(drop -> attemptPickup(character, drop)));
   }

   protected static void attemptPickup(CharacterData character, Drop drop) {
      if (System.currentTimeMillis() - drop.dropTime() < 400) {
         CancelDropReservationProducer.emit(drop.id(), character.id());
         return;
      }
      if (!canBePickedBy(character, drop)) {
         CancelDropReservationProducer.emit(drop.id(), character.id());
         return;
      }

      if (isOwnerLockedMap(character.mapId()) && drop.playerDrop() && drop.ownerId() != character.id()) {
         CancelDropReservationProducer.emit(drop.id(), character.id());
         // emit item unavailable.
         return;
      }

      if (drop.itemId() == 4031865 || drop.itemId() == 4031866) {
         pickupNx(character, drop);
      } else if (drop.meso() > 0) {
         pickupMeso(character, drop);
      } else if (consumeOnPickup(drop.itemId())) {

      } else {
         if (!needsQuestItem(character, drop)) {
            CancelDropReservationProducer.emit(drop.id(), character.id());
            // emit item unavailable.
            return;
         }

         if (!hasInventorySpace(character, drop)) {
            CancelDropReservationProducer.emit(drop.id(), character.id());
            // emit inventory full.
            // emit show inventory full.
            return;
         }

         if (scriptedItem(drop.itemId())) {
            // TODO handle scripted item
         }

         InventoryType inventoryType = getInventoryType(drop.itemId()).orElseThrow();
         if (inventoryType == InventoryType.EQUIP) {
            pickupEquip(character, drop);
         } else {
            pickupItem(character, inventoryType, drop);
         }

         // TODO update ariant score if 4031868

         PickedUpItemProducer.emit(character.id(), drop.itemId(), drop.quantity());
      }
      PickupDropCommandProducer.emit(character.id(), drop.id());
   }

   protected static void pickupEquip(CharacterData character, Drop drop) {
      EquipmentProcessor.createEquipmentForCharacter(character.id(), drop.itemId(), false)
            .ifPresent(data -> InventoryModificationProducer
                  .emit(character.id(), 0, drop.itemId(), drop.quantity(), InventoryType.EQUIP.getType(), data.slot()));
   }

   protected static void pickupItem(CharacterData character, InventoryType inventoryType, Drop drop) {
      int slotMax = maxInSlot(character, drop);

      int quantity = drop.quantity();

      List<ItemData> existingItems = ItemProcessor.getItemsForCharacter(character.id(), inventoryType, drop.itemId());
      // Breaks for a rechargeable item.
      if (existingItems.size() > 0) {
         Iterator<ItemData> itemIterator = existingItems.iterator();
         while (quantity > 0) {
            if (itemIterator.hasNext()) {
               ItemData itemData = itemIterator.next();
               int oldQuantity = itemData.quantity();
               if (oldQuantity < slotMax) {
                  int newQuantity = Math.min(oldQuantity + quantity, slotMax);
                  quantity -= (newQuantity - oldQuantity);
                  ItemProcessor.updateItemQuantity(itemData.id(), newQuantity);
                  InventoryModificationProducer.emit(character.id(), 1, drop.itemId(), newQuantity, inventoryType.getType(),
                        itemData.slot());
               }
            } else {
               break;
            }
         }
      }
      while (quantity > 0) {
         int newQuantity = Math.min(quantity, slotMax);
         quantity -= newQuantity;
         ItemProcessor
               .createItemForCharacter(character.id(), inventoryType, drop.itemId(), newQuantity)
               .ifPresent(data -> InventoryModificationProducer
                     .emit(character.id(), 0, drop.itemId(), newQuantity, inventoryType.getType(), data.slot()));
      }
   }

   protected static int maxInSlot(CharacterData character, Drop drop) {
      return 200;
   }

   protected static Optional<InventoryType> getInventoryType(int itemId) {
      byte type = (byte) (itemId / 1000000);
      if (type >= 1 && type <= 5) {
         return InventoryType.getByType(type);
      }
      return Optional.empty();
   }

   protected static void pickupNx(CharacterData character, Drop drop) {
      int gain = drop.itemId() == 4031865 ? 100 : 250;
      PickedUpNxProducer.emit(character.id(), gain);
   }

   protected static void pickupMeso(CharacterData character, Drop drop) {
      GainMesoProducer.command(character.id(), drop.meso());
   }

   protected static boolean isOwnerLockedMap(int mapId) {
      return (mapId > 209000000 && mapId < 209000016)
            || (mapId >= 990000500 && mapId <= 990000502);
   }

   protected static boolean scriptedItem(int itemId) {
      return itemId / 10000 == 243;
   }

   protected static boolean consumeOnPickup(int itemId) {
      return false;
   }

   protected static boolean canBePickedBy(CharacterData character, Drop drop) {
      if (drop.ownerId() <= 0 || drop.isFFADrop()) {
         return true;
      }

      Optional<Party> ownerParty = PartyProcessor.getPartyForCharacter(drop.ownerId());
      if (ownerParty.isEmpty()) {
         if (character.id() == drop.ownerId()) {
            return true;
         }
      } else {
         Optional<Party> party = PartyProcessor.getPartyForCharacter(character.id());
         if (party.isPresent() && party.get() == ownerParty.get()) {
            return true;
         } else if (character.id() == drop.ownerId()) {
            return true;
         }
      }
      return drop.hasExpiredOwnershipTime();
   }

   protected static Optional<Drop> getDrop(int dropId) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("drops", dropId)
            .getRestClient(DropAttributes.class)
            .getWithResponse()
            .result()
            .flatMap(DataContainer::data)
            .map(ModelFactory::createDrop);
   }

   protected static boolean needsQuestItem(CharacterData character, Drop drop) {
      return true;
   }

   protected static boolean hasInventorySpace(CharacterData character, Drop drop) {
      return true;
   }
}
