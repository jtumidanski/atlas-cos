package com.atlas.cos.event;

import java.util.Collection;

public record CharacterStatUpdateEvent(int worldId, int channelId, int mapId, int characterId, Collection<StatUpdateType> updates) {
}
