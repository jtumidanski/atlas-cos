package com.atlas.cos.processor;

import java.util.Optional;

import com.atlas.cos.event.producer.CancelDropReservationProducer;
import com.atlas.cos.event.producer.PickedUpItemProducer;
import com.atlas.cos.event.producer.PickedUpMesoProducer;
import com.atlas.cos.event.producer.PickedUpNxProducer;
import com.atlas.cos.event.producer.PickupDropCommandProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.Drop;
import com.atlas.cos.model.Party;
import com.atlas.drg.rest.attribute.DropAttributes;
import com.atlas.shared.rest.RestService;
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

         // TODO actually add item to inventory
         // TODO update ariant score if 4031868

         PickedUpItemProducer.emit(character.id(), drop.itemId(), drop.quantity());
      }
      PickupDropCommandProducer.emit(character.id(), drop.id());
   }

   protected static void pickupNx(CharacterData character, Drop drop) {
      int gain = drop.itemId() == 4031865 ? 100 : 250;
      PickedUpNxProducer.emit(character.id(), gain);
   }

   protected static void pickupMeso(CharacterData character, Drop drop) {
      CharacterProcessor.gainMeso(character.id(), drop.meso());
      PickedUpMesoProducer.emit(character.id(), drop.meso());
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
      return UriBuilder.service(RestService.DROP_REGISTRY)
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
