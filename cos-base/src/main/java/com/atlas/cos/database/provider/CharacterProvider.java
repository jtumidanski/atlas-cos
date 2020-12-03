package com.atlas.cos.database.provider;

import java.util.List;
import java.util.Optional;
import javax.persistence.EntityManager;

import com.app.database.provider.NamedQueryClient;
import com.atlas.cos.database.transformer.CharacterDataTransformer;
import com.atlas.cos.entity.Character;
import com.atlas.cos.model.CharacterData;

public class CharacterProvider {
   private CharacterProvider() {
   }

   public static List<CharacterData> getForAccountAndWorld(EntityManager entityManager, int accountId, int worldId) {
      return new NamedQueryClient<>(entityManager, Character.GET_BY_ACCOUNT_AND_WORLD, Character.class)
            .setParameter(Character.ACCOUNT_ID, accountId)
            .setParameter(Character.WORLD_ID, worldId)
            .list(new CharacterDataTransformer());
   }

   public static List<CharacterData> getForName(EntityManager entityManager, String name) {
      return new NamedQueryClient<>(entityManager, Character.GET_BY_NAME, Character.class)
            .setParameter(Character.NAME, name)
            .list(new CharacterDataTransformer());
   }

   public static Optional<CharacterData> getById(EntityManager entityManager, int characterId) {
      return new NamedQueryClient<>(entityManager, Character.GET_BY_ID, Character.class)
            .setParameter(Character.ID, characterId)
            .element(new CharacterDataTransformer());
   }

   public static List<CharacterData> getForWorldInMap(EntityManager entityManager, int worldId, int mapId) {
      return new NamedQueryClient<>(entityManager, Character.GET_BY_WORLD_AND_MAP, Character.class)
            .setParameter(Character.WORLD_ID, worldId)
            .setParameter(Character.MAP_ID, mapId)
            .list(new CharacterDataTransformer());
   }
}
