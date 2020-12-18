package com.atlas.cos.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.attribute.CharacterSeedAttributes;
import com.atlas.cos.attribute.LocationAttributes;
import com.atlas.cos.rest.processor.CharacterRequestProcessor;
import com.atlas.cos.rest.processor.CharacterSeedRequestProcessor;
import com.atlas.cos.rest.processor.DamageProcessor;
import com.atlas.cos.rest.processor.InventoryRequestProcessor;
import com.atlas.cos.rest.processor.SavedLocationProcessor;

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
                                 @QueryParam("mapId") Integer mapId,
                                 @QueryParam("name") String name) {
      if (accountId != null && worldId != null) {
         return CharacterRequestProcessor.getForAccountAndWorld(accountId, worldId).build();
      } else if (name != null) {
         return CharacterRequestProcessor.getByName(name).build();
      } else if (worldId != null && mapId != null) {
         return CharacterRequestProcessor.getForWorldInMap(worldId, mapId).build();
      }
      return ResultBuilder.ok().build();
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

   @GET
   @Path("/{characterId}/inventories")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getInventoryForCharacter(@PathParam("characterId") Integer characterId,
                                            @QueryParam("type") String type,
                                            @QueryParam("include") String include) {
      if (type != null) {
         return InventoryRequestProcessor.getInventoryByType(characterId, type, include).build();
      }
      return InventoryRequestProcessor.getAllInventories(characterId, include).build();
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

   @GET
   @Path("/{characterId}/locations")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getSavedLocations(@PathParam("characterId") Integer characterId,
                                     @QueryParam("type") String type) {
      if (type != null) {
         return SavedLocationProcessor.getSavedLocationsByType(characterId, type).build();
      }
      return SavedLocationProcessor.getSavedLocations(characterId).build();
   }

   @POST
   @Path("/{characterId}/locations")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response addSavedLocation(@PathParam("characterId") Integer characterId,
                                    InputBody<LocationAttributes> inputBody) {
      return SavedLocationProcessor.addSavedLocation(characterId, inputBody.attributes()).build();
   }

   @GET
   @Path("/{characterId}/damage/weapon")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getCharacterDamage(@PathParam("characterId") Integer characterId) {
      return DamageProcessor.getWeaponDamage(characterId).build();
   }
}
