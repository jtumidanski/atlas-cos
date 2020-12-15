package com.atlas.cos.command;

public record PickupCommand(int worldId, int channelId, int characterId, int dropId) {
}
