package com.atlas.cos.processor;

import java.util.Optional;

import com.atlas.cos.event.producer.PickedUpItemProducer;
import com.atlas.cos.event.producer.PickedUpMesoProducer;
import com.atlas.cos.event.producer.PickedUpNxProducer;
import com.atlas.cos.event.producer.PickupDropCommandProducer;
import com.atlas.cos.model.Drop;
import com.atlas.drg.rest.attribute.DropAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import rest.DataContainer;

public final class DropProcessor {
   private DropProcessor() {
   }

   public static void attemptPickup(int characterId, int dropId) {
      CharacterProcessor.getById(characterId)
            .ifPresent(character -> getDrop(dropId).ifPresent(drop -> {
               // TODO some smart verification.
               if (drop.itemId() == 4031865 || drop.itemId() == 4031866) {
                  PickedUpNxProducer.emit(characterId, 10);
               } else if (drop.meso() > 0) {
                  PickedUpMesoProducer.emit(characterId, drop.meso());
               } else {
                  PickedUpItemProducer.emit(characterId, drop.itemId(), drop.quantity());
               }
               PickupDropCommandProducer.emit(characterId, dropId);
            }));
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
}
