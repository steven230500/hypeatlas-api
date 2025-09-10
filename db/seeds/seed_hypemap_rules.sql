-- Evento ejemplo (si no existe)
INSERT INTO
    app.events (
        id,
        slug,
        title,
        game,
        league,
        starts_at,
        ends_at
    )
VALUES (
        gen_random_uuid (),
        'vct-emea-final',
        'VCT EMEA Final',
        'val',
        'VCT EMEA',
        NOW() - INTERVAL '1 day',
        NOW() + INTERVAL '7 days'
    ) ON CONFLICT (slug) DO NOTHING;

-- Reglas de stream -> evento
INSERT INTO
    app.event_stream_rules (platform, handle, event_slug)
VALUES (
        'twitch',
        'koi',
        'vct-emea-final'
    ),
    (
        'twitch',
        'sergiofferra',
        'vct-emea-final'
    ) ON CONFLICT (platform, handle) DO
UPDATE
SET
    event_slug = EXCLUDED.event_slug;

-- Ventana activa (ejemplo)
INSERT INTO
    app.event_windows (
        event_slug,
        starts_at,
        ends_at,
        region,
        lang
    )
VALUES (
        'vct-emea-final',
        NOW() - INTERVAL '1 hour',
        NOW() + INTERVAL '6 hours',
        'EMEA',
        'es'
    );