package com.atlas.cos.command;

public record ChangeMapCommand(int worldId, int channelId, int characterId, int mapId, int portalId) {
}
