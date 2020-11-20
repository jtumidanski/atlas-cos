package com.atlas.cos.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.DefaultValue;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.app.rest.RelationshipInputBody;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.attribute.EquipmentAttributes;
import com.atlas.cos.processor.CharacterResultProcessor;
import com.atlas.cos.processor.EquippedItemResultProcessor;
import com.atlas.cos.processor.ItemResultProcessor;

import builder.ResultBuilder;
import rest.InputBody;

@Path("characters")
public class CharacterResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getCharacters(@QueryParam("accountId") Integer accountId,
                                 @QueryParam("worldId") Integer worldId,
                                 @QueryParam("name") String name) {
      if (accountId != null && worldId != null) {
         return CharacterResultProcessor.getInstance().getForAccountAndWorld(accountId, worldId).build();
      } else if (name != null) {
         return CharacterResultProcessor.getInstance().getByName(name).build();
      }
      return new ResultBuilder().build();
   }

   @POST
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response createCharacter(InputBody<CharacterAttributes> inputBody) {
      return CharacterResultProcessor.getInstance().createCharacter(inputBody.attributes()).build();
   }

   @POST
   @Path("/{characterId}/equipmentSlots/relationships/equipment")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response equipEquipment(@PathParam("characterId") Integer characterId,
                                  RelationshipInputBody relationshipInputBody) {
      int itemId = Integer.parseInt(relationshipInputBody.getData().getId());
      return EquippedItemResultProcessor.getInstance().equipForCharacter(characterId, itemId).build();
   }

   @GET
   @Path("/{characterId}/inventories/equipment")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getEquipment(@PathParam("characterId") Integer characterId,
                                @DefaultValue("false") @QueryParam("filter[equipped]") Boolean equipped) {
      if (equipped != null) {
         return ItemResultProcessor.getInstance().getEquippedItemsForCharacter(characterId)
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
      return ItemResultProcessor.getInstance().createEquipmentForCharacter(characterId, inputBody.attributes(), characterCreation)
            .build();
   }

   @GET
   @Path("/{characterId}/inventories/equipment/{equipmentId}")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response createEquipment(@PathParam("characterId") Integer characterId, @PathParam("equipmentId") Integer equipmentId) {
      return ItemResultProcessor.getInstance().getEquipmentForCharacter(characterId, equipmentId)
            .build();
   }
}
