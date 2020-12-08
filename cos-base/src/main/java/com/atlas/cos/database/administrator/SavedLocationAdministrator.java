package com.atlas.cos.database.administrator;

import javax.persistence.EntityManager;
import javax.persistence.Query;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.entity.LocationType;
import com.atlas.cos.entity.SavedLocation;

public final class SavedLocationAdministrator {
   private SavedLocationAdministrator() {
   }

   public static void deleteForCharacter(EntityManager entityManager, int characterId) {
      Query query = entityManager.createQuery("DELETE FROM SavedLocation WHERE characterId = :characterId");
      query.setParameter("characterId", characterId);
      QueryAdministratorUtil.execute(entityManager, query);
   }

   public static void create(EntityManager entityManager, int characterId, LocationType locationType, int mapId, int portalId) {
      SavedLocation savedLocation = new SavedLocation();
      savedLocation.setCharacterId(characterId);
      savedLocation.setLocationType(locationType);
      savedLocation.setMap(mapId);
      savedLocation.setPortal(portalId);
      QueryAdministratorUtil.insert(entityManager, savedLocation);
   }

   public static void deleteForCharacterByType(EntityManager entityManager, int characterId, LocationType type) {
      Query query = entityManager.createQuery("DELETE FROM SavedLocation WHERE characterId = :characterId and locationType = "
            + ":type");
      query.setParameter("characterId", characterId);
      query.setParameter("type", type);
      QueryAdministratorUtil.execute(entityManager, query);
   }
}