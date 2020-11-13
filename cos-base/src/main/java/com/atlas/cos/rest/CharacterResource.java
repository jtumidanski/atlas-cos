package com.atlas.cos.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.cos.processor.CharacterResultProcessor;

import builder.ResultBuilder;

@Path("")
public class CharacterResource {
   @GET
   @Path("characters")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getCharacters(@QueryParam("accountId") Integer accountId, @QueryParam("worldId") Integer worldId) {
      if (accountId != null && worldId != null) {
         return CharacterResultProcessor.getInstance().getForAccountAndWorld(accountId, worldId).build();
      }
      return new ResultBuilder().build();
   }
}
