package com.atlas.cos.model;

public record HpMpSummary(int hp, int mp) {
   public HpMpSummary increaseHp(int hp) {
      return new HpMpSummary(this.hp + hp, mp);
   }

   public HpMpSummary increaseMp(int mp) {
      return new HpMpSummary(hp, this.mp + mp);
   }
}
