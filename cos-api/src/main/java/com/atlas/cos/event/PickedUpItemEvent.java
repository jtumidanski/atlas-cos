package com.atlas.cos.event;

public record PickedUpItemEvent(int characterId, int itemId, int quantity) {
}
