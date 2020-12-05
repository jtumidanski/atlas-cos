package com.atlas.cos;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import com.atlas.cos.model.CharacterTemporalData;

public class CharacterTemporalRegistry {
   private static final Object lock = new Object();

   private static volatile CharacterTemporalRegistry instance;

   private final Map<Integer, CharacterTemporalData> temporalDataMap;

   private final Map<Integer, Integer> locks;

   public static CharacterTemporalRegistry getInstance() {
      CharacterTemporalRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new CharacterTemporalRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   private CharacterTemporalRegistry() {
      locks = new ConcurrentHashMap<>();
      temporalDataMap = new HashMap<>();
   }

   protected Object getLock(final Integer characterId) {
      locks.putIfAbsent(characterId, characterId);
      return locks.get(characterId);
   }

   public CharacterTemporalData getTemporalData(int characterId) {
      CharacterTemporalData data = temporalDataMap.get(characterId);
      if (data == null) {
         synchronized (getLock(characterId)) {
            data = getDefaultData();
            temporalDataMap.put(characterId, data);
         }
      }
      return data;
   }

   public void updatePosition(int characterId, int x, int y) {
      synchronized (getLock(characterId)) {
         CharacterTemporalData data = temporalDataMap
               .getOrDefault(characterId, getDefaultData())
               .updatePosition(x, y);
         temporalDataMap.put(characterId, data);
      }
   }

   public void update(int characterId, int x, int y, byte stance) {
      synchronized (getLock(characterId)) {
         CharacterTemporalData data = temporalDataMap
               .getOrDefault(characterId, getDefaultData())
               .update(x, y, stance);
         temporalDataMap.put(characterId, data);
      }
   }

   public void updateStance(int characterId, byte stance) {
      synchronized (getLock(characterId)) {
         CharacterTemporalData data = temporalDataMap
               .getOrDefault(characterId, getDefaultData())
               .updateStance(stance);
         temporalDataMap.put(characterId, data);
      }
   }

   protected CharacterTemporalData getDefaultData() {
      return new CharacterTemporalData(0, 0, (byte) 0);
   }
}
