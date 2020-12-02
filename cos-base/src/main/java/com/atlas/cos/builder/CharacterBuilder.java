package com.atlas.cos.builder;

import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.configuration.Configuration;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.processor.ConfigurationProcessor;

public class CharacterBuilder {
   private int accountId;

   private int worldId;

   private String name;

   private int level;

   private int strength;

   private int dexterity;

   private int luck;

   private int intelligence;

   private int maxHp;

   private int maxMp;

   private int jobId;

   private int skinColor;

   private byte gender;

   private int hair;

   private int face;

   private int ap;

   private int mapId;

   public CharacterBuilder(int accountId, int worldId, String name, int jobId, int skinColor, byte gender, int hair, int face) {
      this.accountId = accountId;
      this.worldId = worldId;
      this.name = name;
      this.jobId = jobId;
      this.level = 1;

      this.skinColor = skinColor;
      this.gender = gender;
      this.hair = hair;
      this.face = face;

      Configuration configuration = ConfigurationProcessor.getInstance().getConfiguration();
      if (!configuration.useStarting4Ap) {
         if (configuration.useAutoAssignStartersAp) {
            this.strength = 12;
            this.dexterity = 5;
            this.intelligence = 4;
            this.luck = 4;
            this.ap = 0;
         } else {
            this.strength = 4;
            this.dexterity = 4;
            this.intelligence = 4;
            this.luck = 4;
            this.ap = 9;
         }
      } else {
         this.strength = 4;
         this.dexterity = 4;
         this.intelligence = 4;
         this.luck = 4;
         this.ap = 0;
      }

      this.maxHp = 50;
      this.maxMp = 5;
   }

   public CharacterBuilder(int accountId, int worldId, String name, int jobId, int skinColor, byte gender, int hair, int face,
                           int level, int mapId) {
      this(accountId, worldId, name, jobId, skinColor, gender, hair, face);
      this.level = level;
      this.mapId = mapId;
   }

   public CharacterBuilder(CharacterAttributes attributes, int level, int mapId) {
      this(attributes.accountId(), attributes.worldId(), attributes.name(), attributes.jobId(), attributes.skinColor(),
            attributes.gender(), attributes.hair(), attributes.face(), level, mapId);
   }

   public CharacterBuilder setMapId(int mapId) {
      this.mapId = mapId;
      return this;
   }

   public CharacterData build() {
      return new CharacterData(-1, accountId, worldId, name, level, 0, 0, strength, dexterity, luck,
            intelligence, 0, 0, maxHp, maxMp, 0, 0, jobId, skinColor, gender, 0, hair, face, ap,
            "", mapId, 0, 0);
   }
}
