package com.atlas.cos.rest.processor;

import builder.ResultBuilder;
import com.app.rest.util.stream.Mappers;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.MapleJob;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.ItemProcessor;
import com.atlas.cos.processor.JobProcessor;
import com.atlas.cos.rest.ResultObjectFactory;

import javax.ws.rs.core.Response;
import java.util.Optional;
import java.util.function.Function;

public final class CharacterSeedRequestProcessor {
   private CharacterSeedRequestProcessor() {
   }

   public static ResultBuilder create(int accountId, int worldId, String name, int jobIndex, int face, int hair, int hairColor, int skin,
                                      byte gender, int top, int bottom, int shoes, int weapon) {
      MapleJob job = JobProcessor.getJobFromIndex(jobIndex).orElse(MapleJob.BEGINNER);

      Function<CharacterBuilder, Optional<CharacterData>> creator;
      if (MapleJob.BEGINNER.equals(job)) {
         creator = CharacterProcessor::createBeginner;
      } else if (MapleJob.NOBLESSE.equals(job)) {
         creator = CharacterProcessor::createNoblesse;
      } else if (MapleJob.LEGEND.equals(job)) {
         creator = CharacterProcessor::createLegend;
      } else {
         return new ResultBuilder(Response.Status.NOT_IMPLEMENTED);
      }

      Optional<CharacterData> result = creator.apply(new CharacterBuilder(accountId, worldId, name, job.getId(), skin, gender,
            hair + hairColor, face));
      if (result.isEmpty()) {
         return new ResultBuilder(Response.Status.FORBIDDEN);
      }

      ItemProcessor.createEquipmentForCharacter(result.get().id(), top, true)
            .ifPresent(equipment -> ItemProcessor.equipItemForCharacter(result.get().id(), equipment.id()));
      ItemProcessor.createEquipmentForCharacter(result.get().id(), bottom, true)
            .ifPresent(equipment -> ItemProcessor.equipItemForCharacter(result.get().id(), equipment.id()));
      ItemProcessor.createEquipmentForCharacter(result.get().id(), shoes, true)
            .ifPresent(equipment -> ItemProcessor.equipItemForCharacter(result.get().id(), equipment.id()));
      ItemProcessor.createEquipmentForCharacter(result.get().id(), weapon, true)
            .ifPresent(equipment -> ItemProcessor.equipItemForCharacter(result.get().id(), equipment.id()));

      return result
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.FORBIDDEN));
   }
}
