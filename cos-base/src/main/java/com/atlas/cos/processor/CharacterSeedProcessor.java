package com.atlas.cos.processor;

import java.util.Optional;
import java.util.function.Function;
import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Mappers;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.MapleJob;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;

public class CharacterSeedProcessor {
   private static final Object lock = new Object();

   private static volatile CharacterSeedProcessor instance;

   public static CharacterSeedProcessor getInstance() {
      CharacterSeedProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new CharacterSeedProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public ResultBuilder create(int accountId, int worldId, String name, int jobIndex, int face, int hair, int hairColor, int skin,
                               byte gender, int top, int bottom, int shoes, int weapon) {
      MapleJob job = JobProcessor.getInstance().getJobFromIndex(jobIndex).orElse(MapleJob.BEGINNER);

      Function<CharacterBuilder, Optional<CharacterData>> creator;
      if (MapleJob.BEGINNER.equals(job)) {
         creator = CharacterProcessor.getInstance()::createBeginner;
      } else if (MapleJob.NOBLESSE.equals(job)) {
         creator = CharacterProcessor.getInstance()::createNoblesse;
      } else if (MapleJob.LEGEND.equals(job)) {
         creator = CharacterProcessor.getInstance()::createLegend;
      } else {
         return new ResultBuilder(Response.Status.NOT_IMPLEMENTED);
      }

      Optional<CharacterData> result = creator.apply(new CharacterBuilder(accountId, worldId, name, job.getId(), skin, gender,
            hair + hairColor, face));
      if (result.isEmpty()) {
         return new ResultBuilder(Response.Status.FORBIDDEN);
      }

      ItemProcessor.getInstance()
            .createEquipmentForCharacter(result.get().id(), top, true)
            .ifPresent(equipment -> ItemProcessor.getInstance().equipItemForCharacter(result.get().id(), equipment.id()));
      ItemProcessor.getInstance().createEquipmentForCharacter(result.get().id(), bottom, true)
            .ifPresent(equipment -> ItemProcessor.getInstance().equipItemForCharacter(result.get().id(), equipment.id()));
      ItemProcessor.getInstance().createEquipmentForCharacter(result.get().id(), shoes, true)
            .ifPresent(equipment -> ItemProcessor.getInstance().equipItemForCharacter(result.get().id(), equipment.id()));
      ItemProcessor.getInstance().createEquipmentForCharacter(result.get().id(), weapon, true)
            .ifPresent(equipment -> ItemProcessor.getInstance().equipItemForCharacter(result.get().id(), equipment.id()));

      return result
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.FORBIDDEN));
   }
}
