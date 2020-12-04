package com.atlas.cos.model;

public record CharacterTemporalData(int x, int y, byte stance) {
   public CharacterTemporalData update(int x, int y, byte stance) {
      return new CharacterTemporalData(x, y, stance);
   }

   public CharacterTemporalData updatePosition(int x, int y) {
      return new CharacterTemporalData(x, y, stance);
   }

   public CharacterTemporalData updateStance(byte stance) {
      return new CharacterTemporalData(x, y, stance);
   }
}
