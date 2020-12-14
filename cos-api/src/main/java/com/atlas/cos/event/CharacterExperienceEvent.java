package com.atlas.cos.event;

public record CharacterExperienceEvent(int worldId, int channelId, int mapId, int characterId, int personalGain,
                                       int partyGain, boolean show, boolean chat, boolean white) {
}
