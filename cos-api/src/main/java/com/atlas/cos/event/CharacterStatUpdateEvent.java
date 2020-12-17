package com.atlas.cos.event;

import java.util.Collection;

public record CharacterStatUpdateEvent(int characterId, Collection<StatUpdateType> updates) {
}
