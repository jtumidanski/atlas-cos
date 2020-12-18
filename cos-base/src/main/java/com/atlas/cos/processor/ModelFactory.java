package com.atlas.cos.processor;

import com.atlas.cos.model.Drop;
import com.atlas.cos.model.EquipData;
import com.atlas.cos.model.Monster;
import com.atlas.cos.model.Portal;
import com.atlas.drg.rest.attribute.DropAttributes;
import com.atlas.eso.attribute.EquipmentAttributes;
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

   public static Drop createDrop(DataBody<DropAttributes> body) {
      return new Drop(Integer.parseInt(body.getId()),
            body.getAttributes().itemId(),
            body.getAttributes().quantity(),
            body.getAttributes().meso(),
            body.getAttributes().dropTime(),
            body.getAttributes().dropType(),
            body.getAttributes().ownerId(),
            body.getAttributes().playerDrop());
   }

   public static EquipData createEquip(DataBody<EquipmentAttributes> body) {
      return new EquipData(
            Integer.parseInt(body.getId()),
            body.getAttributes().itemId(),
            body.getAttributes().strength(),
            body.getAttributes().dexterity(),
            body.getAttributes().intelligence(),
            body.getAttributes().luck(),
            body.getAttributes().hp(),
            body.getAttributes().mp(),
            body.getAttributes().weaponAttack(),
            body.getAttributes().magicAttack(),
            body.getAttributes().weaponDefense(),
            body.getAttributes().magicDefense(),
            body.getAttributes().accuracy(),
            body.getAttributes().avoidability(),
            body.getAttributes().hands(),
            body.getAttributes().speed(),
            body.getAttributes().jump(),
            body.getAttributes().slots()
      );
   }
}
