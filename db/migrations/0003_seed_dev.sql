-- +goose Up
-- Evento y creador de ejemplo
INSERT INTO
    app.events (
        slug,
        title,
        game,
        league,
        starts_at
    )
VALUES (
        'vct-emea-final',
        'VCT EMEA Final',
        'val',
        'VCT EMEA',
        now() + interval '1 day'
    ) ON CONFLICT (slug) DO NOTHING;

INSERT INTO
    app.creators (
        platform,
        handle,
        url,
        lang,
        country,
        verified
    )
VALUES (
        'twitch',
        'caster1',
        'https://twitch.tv/caster1',
        'es',
        'ES',
        true
    ) ON CONFLICT DO NOTHING;

-- Co-stream “en vivo”
INSERT INTO
    app.co_streams (
        event_id,
        creator_id,
        platform,
        url,
        lang,
        country,
        viewers,
        verified,
        is_live,
        last_seen_at
    )
SELECT e.id, c.id, 'twitch', 'https://twitch.tv/caster1', 'es', 'ES', 8200, true, true, now()
FROM app.events e, app.creators c
WHERE
    e.slug = 'vct-emea-final'
    AND c.handle = 'caster1' ON CONFLICT DO NOTHING;

-- Patch de ejemplo y un cambio
INSERT INTO
    app.patches (game, version, released_at)
VALUES (
        'val',
        '9.15',
        now() - interval '8 days'
    ) ON CONFLICT DO NOTHING;

INSERT INTO
    app.patch_changes (
        patch_id,
        entity_type,
        entity_id,
        field,
        old,
        new,
        impact_score
    )
SELECT p.id, 'agent', 'sova', 'recon bolt cd', '40s', '45s', 0.6
FROM app.patches p
WHERE
    p.game = 'val'
    AND p.version = '9.15';

-- +goose Down
-- (sin down para seed)