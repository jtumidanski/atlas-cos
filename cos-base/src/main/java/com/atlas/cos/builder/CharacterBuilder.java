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

   public CharacterBuilder(CharacterAttributes attributes, int level, int mapId) {
      this.accountId = attributes.accountId();
      this.worldId = attributes.worldId();
      this.name = attributes.name();
      this.jobId = attributes.jobId();
      this.level = level;
      this.mapId = mapId;

      this.skinColor = attributes.skinColor();
      this.gender = attributes.gender();
      this.hair = attributes.hair();
      this.face = attributes.face();

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

   public CharacterData build() {
      return new CharacterData(-1, accountId, worldId, name, level, 0, 0, strength, dexterity, luck,
            intelligence, 0, 0, maxHp, maxMp, 0, 0, jobId, skinColor, gender, 0, hair, face, ap,
            "", mapId, 0, 0);
   }
}
