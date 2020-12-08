package com.atlas.cos.database.provider;

import java.util.List;
import javax.persistence.EntityManager;
import javax.persistence.TypedQuery;

import com.app.database.util.QueryProviderUtil;
import com.atlas.cos.database.transformer.SavedLocationTransformer;
import com.atlas.cos.entity.SavedLocation;
import com.atlas.cos.model.SavedLocationData;

public final class SavedLocationProvider {
   private SavedLocationProvider() {
   }

   public static List<SavedLocationData> getForCharacter(EntityManager entityManager, int characterId) {
      TypedQuery<SavedLocation> query =
            entityManager.createQuery("FROM SavedLocation s WHERE s.characterId = :characterId", SavedLocation.class);
      query.setParameter("characterId", characterId);
      return QueryProviderUtil.list(query, new SavedLocationTransformer());
   }
}