package com.atlas.cos.builder;

import com.atlas.cos.ConfigurationRegistry;
import com.atlas.cos.rest.attribute.CharacterAttributes;
import com.atlas.cos.configuration.Configuration;
import com.atlas.cos.model.CharacterData;

public class CharacterBuilder {
   private final int accountId;

   private final int worldId;

   private final String name;

   private int level;

   private final int strength;

   private final int dexterity;

   private final int luck;

   private final int intelligence;

   private final int maxHp;

   private final int maxMp;

   private int jobId;

   private final int skinColor;

   private final byte gender;

   private final int hair;

   private final int face;

   private final int ap;

   private int mapId;

   public CharacterBuilder(int accountId, int worldId, String name, int skinColor, byte gender, int hair, int face) {
      this.accountId = accountId;
      this.worldId = worldId;
      this.name = name;
      this.level = 1;
      this.jobId = 0;

      this.skinColor = skinColor;
      this.gender = gender;
      this.hair = hair;
      this.face = face;

      Configuration configuration = ConfigurationRegistry.getInstance().getConfiguration();
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

   public CharacterBuilder(int accountId, int worldId, String name, int skinColor, byte gender, int hair, int face,
                           int level, int mapId) {
      this(accountId, worldId, name, skinColor, gender, hair, face);
      this.level = level;
      this.mapId = mapId;
   }

   public CharacterBuilder(CharacterAttributes attributes, int level, int mapId) {
      this(attributes.accountId(), attributes.worldId(), attributes.name(), attributes.skinColor(),
            attributes.gender(), attributes.hair(), attributes.face(), level, mapId);
      setJobId(attributes.jobId());
   }

   public CharacterBuilder setJobId(int jobId) {
      this.jobId = jobId;
      return this;
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
