package com.atlas.cos.model;

public record CharacterTemporalData(int x, int y, int stance) {
   public CharacterTemporalData updatePosition(int x, int y) {
      return new CharacterTemporalData(x, y, stance);
   }
}
