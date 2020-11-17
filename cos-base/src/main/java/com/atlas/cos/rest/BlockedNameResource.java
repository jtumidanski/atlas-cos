package com.atlas.cos.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.cos.processor.BlockedNameProcessor;

@Path("blockedNames")
public class BlockedNameResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getBlockedNames(@QueryParam("name") String name) {
      if (name != null) {
         return BlockedNameProcessor.getInstance().getName(name).build();
      } else {
         return BlockedNameProcessor.getInstance().getNames().build();
      }
   }
}
