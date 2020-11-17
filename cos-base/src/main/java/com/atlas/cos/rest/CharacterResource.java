package com.atlas.cos.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.processor.CharacterResultProcessor;

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
}
