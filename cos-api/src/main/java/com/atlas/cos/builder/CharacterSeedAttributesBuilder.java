package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.CharacterSeedAttributes;

import builder.AttributeResultBuilder;

public class CharacterSeedAttributesBuilder extends RecordBuilder<CharacterSeedAttributes, CharacterSeedAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer accountId;

   private Integer worldId;

   private String name;

   private Integer jobIndex;

   private Integer face;

   private Integer hair;

   private Integer hairColor;

   private Integer skin;

   private Byte gender;

   private Integer top;

   private Integer bottom;

   private Integer shoes;

   private Integer weapon;

   @Override
   public CharacterSeedAttributes construct() {
      return new CharacterSeedAttributes(accountId, worldId, name, jobIndex, face, hair, hairColor, skin, gender, top, bottom,
            shoes, weapon);
   }

   @Override
   public CharacterSeedAttributesBuilder getThis() {
      return this;
   }

   public CharacterSeedAttributesBuilder setAccountId(Integer accountId) {
      this.accountId = accountId;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setWorldId(Integer worldId) {
      this.worldId = worldId;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setName(String name) {
      this.name = name;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setJobIndex(Integer jobIndex) {
      this.jobIndex = jobIndex;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setFace(Integer face) {
      this.face = face;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setHair(Integer hair) {
      this.hair = hair;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setHairColor(Integer hairColor) {
      this.hairColor = hairColor;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setSkin(Integer skin) {
      this.skin = skin;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setGender(Byte gender) {
      this.gender = gender;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setTop(Integer top) {
      this.top = top;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setBottom(Integer bottom) {
      this.bottom = bottom;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setShoes(Integer shoes) {
      this.shoes = shoes;
      return getThis();
   }

   public CharacterSeedAttributesBuilder setWeapon(Integer weapon) {
      this.weapon = weapon;
      return getThis();
   }
}
