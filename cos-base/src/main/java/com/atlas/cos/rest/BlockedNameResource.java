package com.atlas.cos.rest;

import com.atlas.cos.rest.processor.BlockedNameRequestProcessor;

import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

@Path("blockedNames")
public class BlockedNameResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getBlockedNames(@QueryParam("name") String name) {
      if (name != null) {
         return BlockedNameRequestProcessor.getName(name).build();
      } else {
         return BlockedNameRequestProcessor.getNames().build();
      }
   }
}
