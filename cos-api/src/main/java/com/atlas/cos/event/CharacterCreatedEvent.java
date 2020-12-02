package com.atlas.cos.event;

public record CharacterCreatedEvent(int worldId, int characterId, String name) {
}
