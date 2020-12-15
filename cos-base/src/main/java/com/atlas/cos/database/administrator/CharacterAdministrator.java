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

   public static void increaseLevel(EntityManager entityManager, int characterId) {
      QueryAdministratorUtil
            .update(entityManager, Character.class, characterId, character -> character.setLevel(character.getLevel() + 1));
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

   public static void setAp(EntityManager entityManager, int characterId, int ap) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> character.setAp(ap));
   }

   public static void setLuck(EntityManager entityManager, int characterId, int luck) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> character.setLuck(luck));
   }

   public static void setIntelligence(EntityManager entityManager, int characterId, int intelligence) {
      QueryAdministratorUtil
            .update(entityManager, Character.class, characterId, character -> character.setIntelligence(intelligence));
   }

   public static void setDexterity(EntityManager entityManager, int characterId, int dexterity) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> character.setDexterity(dexterity));
   }

   public static void setStrength(EntityManager entityManager, int characterId, int strength) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> character.setStrength(strength));
   }

   public static void setMp(EntityManager entityManager, int characterId, int mp) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> character.setMp(mp));
   }

   public static void setMaxMp(EntityManager entityManager, int characterId, int maxMp) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> character.setMaxMp(maxMp));
   }

   public static void setHp(EntityManager entityManager, int characterId, int hp) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> character.setHp(hp));
   }

   public static void setMaxHp(EntityManager entityManager, int characterId, int hp) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> character.setMaxHp(hp));
   }
}
