package com.atlas.cos.event;

public record MapChangedEvent(int worldId, int channelId, int mapId, int portalId, int characterId) {
}
