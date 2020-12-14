package com.atlas.cos.database.administrator;

import java.util.Optional;
import javax.persistence.EntityManager;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.transformer.CharacterDataTransformer;
import com.atlas.cos.entity.Character;
import com.atlas.cos.model.CharacterData;

public class CharacterAdministrator {
   private CharacterAdministrator() {
   }

   public static Optional<CharacterData> create(EntityManager entityManager, int accountId, int worldId, String name, int level,
                                                int strength,
                                                int dexterity, int luck, int intelligence, int maxHp, int maxMp, int jobId,
                                                byte gender, int hair, int face, int mapId) {
      Character character = new Character();
      character.setAccountId(accountId);
      character.setWorld(worldId);
      character.setName(name);
      character.setLevel(level);
      character.setStrength(strength);
      character.setDexterity(dexterity);
      character.setLuck(luck);
      character.setIntelligence(intelligence);
      character.setMaxHp(maxHp);
      character.setMaxMp(maxMp);
      character.setJob(jobId);
      character.setGender(gender);
      character.setHair(hair);
      character.setFace(face);
      character.setMap(mapId);
      return Optional.of(QueryAdministratorUtil.insertAndReturn(entityManager, character, new CharacterDataTransformer()));
   }

   public static void setExperience(EntityManager entityManager, int characterId, int experience) {
      QueryAdministratorUtil
            .update(entityManager, Character.class, characterId, character -> character.setExp(experience));
   }

   public static void setLevel(EntityManager entityManager, int characterId, int level) {
      QueryAdministratorUtil
            .update(entityManager, Character.class, characterId, character -> character.setLevel(level));
   }

   public static void increaseExperience(EntityManager entityManager, int characterId, int experience) {
      QueryAdministratorUtil
            .update(entityManager, Character.class, characterId, character -> character.setExp(character.getExp() + experience));
   }

   public static void updateMap(EntityManager entityManager, int characterId, int mapId) {
      QueryAdministratorUtil
            .update(entityManager, Character.class, characterId, character -> character.setMap(mapId));
   }

   public static void updateSpawnPoint(EntityManager entityManager, int characterId, int newSpawnPoint) {
      QueryAdministratorUtil
            .update(entityManager, Character.class, characterId, character -> character.setSpawnPoint(newSpawnPoint));
   }
}
