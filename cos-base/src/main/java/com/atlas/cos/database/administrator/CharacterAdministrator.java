package com.atlas.cos.database.administrator;

import javax.persistence.EntityManager;

import com.atlas.cos.database.transformer.CharacterDataTransformer;
import com.atlas.cos.entity.Character;
import com.atlas.cos.model.CharacterData;

import accessor.AbstractQueryExecutor;

public class CharacterAdministrator extends AbstractQueryExecutor {
   private static final Object lock = new Object();

   private static volatile CharacterAdministrator instance;

   public static CharacterAdministrator getInstance() {
      CharacterAdministrator result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new CharacterAdministrator();
               instance = result;
            }
         }
      }
      return result;
   }

   public CharacterData create(EntityManager entityManager, int accountId, int worldId, String name, int level, int strength,
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
      return insertAndReturn(entityManager, character, new CharacterDataTransformer());
   }
}
