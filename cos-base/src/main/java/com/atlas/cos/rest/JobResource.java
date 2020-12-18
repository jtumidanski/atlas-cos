package com.atlas.cos.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.cos.rest.processor.JobRequestProcessor;

import builder.ResultBuilder;

@Path("jobs")
public class JobResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getJobs(@QueryParam("createIndex") Integer createIndex) {
      if (createIndex != null) {
         return JobRequestProcessor.getByCreateIndex(createIndex).build();
      }
      return ResultBuilder.notFound().build();
   }
}
