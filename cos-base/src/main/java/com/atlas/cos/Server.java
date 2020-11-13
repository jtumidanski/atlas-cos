package com.atlas.cos;

import java.net.URI;

import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;
import org.glassfish.grizzly.http.server.HttpServer;

import database.PersistenceManager;

public class Server {
   public static void main(String[] args) {
      PersistenceManager.construct("atlas-cos");
      URI uri = UriBuilder.host(RestService.CHANNEL).uri();
      final HttpServer server = RestServerFactory.create(uri, "com.atlas.cos.rest");
   }
}
