package com.atlas.cos.command;

public record AssignApCommand(int worldId, int channelId, int mapId, int characterId, AssignApType type) {
}
