package com.atlas.cos.processor;

import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Collectors;
import com.atlas.cos.attribute.LocationAttributes;
import com.atlas.cos.database.administrator.SavedLocationAdministrator;
import com.atlas.cos.database.provider.SavedLocationProvider;
import com.atlas.cos.entity.LocationType;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;
import database.Connection;

public final class SavedLocationProcessor {
   private SavedLocationProcessor() {
   }

   public static ResultBuilder getSavedLocations(int characterId) {
      return Connection.instance()
            .list(entityManager -> SavedLocationProvider.getForCharacter(entityManager, characterId))
            .stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public static ResultBuilder getSavedLocationsByType(int characterId, String type) {
      return Connection.instance()
            .list(entityManager -> SavedLocationProvider.getForCharacter(entityManager, characterId))
            .stream()
            .filter(location -> location.type().equalsIgnoreCase(type))
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public static ResultBuilder addSavedLocation(int characterId, LocationAttributes attributes) {
      Connection.instance().with(entityManager -> {
         LocationType type = LocationType.valueOf(attributes.type().name());
         SavedLocationAdministrator.deleteForCharacterByType(entityManager, characterId, type);
         SavedLocationAdministrator.create(entityManager, characterId, type, attributes.mapId(), attributes.portalId());
      });
      return new ResultBuilder(Response.Status.NO_CONTENT);
   }
}
