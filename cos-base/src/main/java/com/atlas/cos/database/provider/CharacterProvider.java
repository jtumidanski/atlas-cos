package com.atlas.cos.database.provider;

import java.util.Arrays;
import java.util.List;
import javax.persistence.EntityManager;
import javax.persistence.TypedQuery;

import com.atlas.cos.database.transformer.CharacterDataTransformer;
import com.atlas.cos.entity.Character;
import com.atlas.cos.model.CharacterData;

import accessor.AbstractQueryExecutor;

public class CharacterProvider extends AbstractQueryExecutor {
   private static final Object lock = new Object();

   private static volatile CharacterProvider instance;

   public static CharacterProvider getInstance() {
      CharacterProvider result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new CharacterProvider();
               instance = result;
            }
         }
      }
      return result;
   }

   public List<CharacterData> getForAccountAndWorld(EntityManager entityManager, int accountId, int worldId) {
      TypedQuery<Character> query = entityManager.createQuery("SELECT c FROM Character c WHERE c.accountId = :accountId AND c"
            + ".world = :worldId", Character.class);
      query.setParameter("accountId", accountId);
      query.setParameter("worldId", worldId);
      return getResultList(query, new CharacterDataTransformer());
   }

   public List<CharacterData> getForName(EntityManager entityManager, String name) {
      TypedQuery<Character> query = entityManager.createQuery("SELECT c FROM Character c WHERE c.name = :name", Character.class);
      query.setParameter("name", name);
      return getResultList(query, new CharacterDataTransformer());
   }
}
