package com.atlas.cos.rest;

import builder.ResultBuilder;
import com.app.rest.RelationshipInputBody;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.attribute.CharacterSeedAttributes;
import com.atlas.cos.attribute.EquipmentAttributes;
import com.atlas.cos.rest.processor.CharacterSeedRequestProcessor;
import com.atlas.cos.rest.processor.EquippedItemRequestProcessor;
import com.atlas.cos.rest.processor.ItemRequestProcessor;
import com.atlas.cos.rest.processor.CharacterRequestProcessor;
import rest.InputBody;

import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

@Path("characters")
public class CharacterResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getCharacters(@QueryParam("accountId") Integer accountId,
                                 @QueryParam("worldId") Integer worldId,
                                 @QueryParam("mapId") Integer mapId,
                                 @QueryParam("name") String name) {
      if (accountId != null && worldId != null) {
         return CharacterRequestProcessor.getForAccountAndWorld(accountId, worldId).build();
      } else if (name != null) {
         return CharacterRequestProcessor.getByName(name).build();
      } else if (worldId != null && mapId != null) {
         return CharacterRequestProcessor.getForWorldInMap(worldId, mapId).build();
      }
      return new ResultBuilder().build();
   }

   @POST
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response createCharacter(InputBody<CharacterAttributes> inputBody) {
      return CharacterRequestProcessor.createCharacter(inputBody.attributes()).build();
   }

   @GET
   @Path("/{characterId}")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getCharacter(@PathParam("characterId") Integer characterId) {
      return CharacterRequestProcessor.getById(characterId).build();
   }

   @POST
   @Path("/{characterId}/equipmentSlots/relationships/equipment")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response equipEquipment(@PathParam("characterId") Integer characterId,
                                  RelationshipInputBody relationshipInputBody) {
      int itemId = Integer.parseInt(relationshipInputBody.getData().getId());
      return EquippedItemRequestProcessor.equipForCharacter(characterId, itemId).build();
   }

   @GET
   @Path("/{characterId}/inventories/equipment")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getEquipment(@PathParam("characterId") Integer characterId,
                                @DefaultValue("false") @QueryParam("filter[equipped]") Boolean equipped) {
      if (equipped != null) {
         return ItemRequestProcessor.getEquippedItemsForCharacter(characterId)
               .build();
      } else {
         return new ResultBuilder(Response.Status.NOT_IMPLEMENTED).build();
      }
   }

   @POST
   @Path("/{characterId}/inventories/equipment")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response createEquipment(@PathParam("characterId") Integer characterId,
                                   @DefaultValue("false") @QueryParam("characterCreation") Boolean characterCreation,
                                   InputBody<EquipmentAttributes> inputBody) {
      return ItemRequestProcessor.createEquipmentForCharacter(characterId, inputBody.attributes(), characterCreation)
            .build();
   }

   @GET
   @Path("/{characterId}/inventories/equipment/{equipmentId}")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response createEquipment(@PathParam("characterId") Integer characterId, @PathParam("equipmentId") Integer equipmentId) {
      return ItemRequestProcessor.getEquipmentForCharacter(characterId, equipmentId)
            .build();
   }

   @POST
   @Path("/seeds")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response createCharacterFromSeed(InputBody<CharacterSeedAttributes> inputBody) {
      return CharacterSeedRequestProcessor.create(
            inputBody.attribute(CharacterSeedAttributes::accountId),
            inputBody.attribute(CharacterSeedAttributes::worldId),
            inputBody.attribute(CharacterSeedAttributes::name),
            inputBody.attribute(CharacterSeedAttributes::jobIndex),
            inputBody.attribute(CharacterSeedAttributes::face),
            inputBody.attribute(CharacterSeedAttributes::hair),
            inputBody.attribute(CharacterSeedAttributes::hairColor),
            inputBody.attribute(CharacterSeedAttributes::skin),
            inputBody.attribute(CharacterSeedAttributes::gender),
            inputBody.attribute(CharacterSeedAttributes::top),
            inputBody.attribute(CharacterSeedAttributes::bottom),
            inputBody.attribute(CharacterSeedAttributes::shoes),
            inputBody.attribute(CharacterSeedAttributes::weapon)
      ).build();
   }
}
