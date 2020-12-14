package com.atlas.cos.processor;

import com.atlas.cos.model.Monster;
import com.atlas.cos.model.Portal;
import com.atlas.mis.attribute.MonsterDataAttributes;
import com.atlas.mis.attribute.PortalAttributes;

import rest.DataBody;

public final class ModelFactory {
   private ModelFactory() {
   }

   public static Portal createPortal(DataBody<PortalAttributes> body) {
      return new Portal(Integer.parseInt(body.getId()),
            body.getAttributes().name(),
            body.getAttributes().target(),
            body.getAttributes().type(),
            body.getAttributes().x(),
            body.getAttributes().y(),
            body.getAttributes().targetMap(),
            body.getAttributes().scriptName()
      );
   }

   public static Monster createMonster(DataBody<MonsterDataAttributes> body) {
      return new Monster(body.getAttributes().experience(), body.getAttributes().hp());
   }
}
