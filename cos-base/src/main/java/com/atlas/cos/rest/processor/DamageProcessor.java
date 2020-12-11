package com.atlas.cos.rest.processor;

import com.atlas.cos.rest.attribute.DamageAttributes;
import com.atlas.cos.rest.attribute.DamageType;
import com.atlas.cos.rest.builder.DamageAttributesBuilder;

import builder.ResultBuilder;
import builder.ResultObjectBuilder;

public final class DamageProcessor {
   private DamageProcessor() {
   }

   public static ResultBuilder getWeaponDamage(int characterId) {
      return new ResultBuilder()
            .addData(new ResultObjectBuilder(DamageAttributes.class, 0)
                  .setAttribute(new DamageAttributesBuilder()
                        .setType(DamageType.WEAPON)
                        .setMaximum(com.atlas.cos.processor.DamageProcessor.getMaxBaseDamage(characterId))
                  )
            );
   }
}
