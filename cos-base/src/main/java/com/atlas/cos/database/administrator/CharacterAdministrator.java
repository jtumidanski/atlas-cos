package com.atlas.cos.database.administrator;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.transformer.CharacterDataTransformer;
import com.atlas.cos.entity.Character;
import com.atlas.cos.model.CharacterData;

import javax.persistence.EntityManager;
import java.util.Optional;

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

   public static void update(EntityManager entityManager, int characterId, int hp, int maxHp, int mp, int maxMp, int strength,
                             int dexterity, int intelligence, int luck, int ap, int level) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId, character -> {
         character.setHp(character.getHp() + hp);
         character.setMaxHp(character.getMaxHp() + maxHp);
         character.setMp(character.getMp() + mp);
         character.setMaxMp(character.getMaxMp() + maxMp);
         character.setStrength(character.getStrength() + strength);
         character.setDexterity(character.getDexterity() + dexterity);
         character.setIntelligence(character.getIntelligence() + intelligence);
         character.setLuck(character.getLuck() + luck);
         character.setAp(character.getAp() + ap);
         character.setLevel(character.getLevel() + level);
      });
   }

   public static void increaseMeso(EntityManager entityManager, int characterId, int meso) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId,
            character -> character.setMeso(character.getMeso() + meso));
   }

   public static void updateSp(EntityManager entityManager, int characterId, int newValue, int skillBookId) {
      QueryAdministratorUtil.update(entityManager, Character.class, characterId,
            character -> {
               String[] sps = character.getSp().split(",");
               sps[skillBookId] = String.valueOf(newValue);
               character.setSp(String.join(",", sps));
            });
   }
}
