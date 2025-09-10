-- +goose Up
INSERT INTO
    app.leagues (game, region, name, slug)
VALUES (
        'val',
        'EMEA',
        'VCT EMEA',
        'vct-emea'
    ),
    ('lol', 'EMEA', 'LEC', 'lec') ON CONFLICT (slug) DO NOTHING;

-- comp de ejemplo (VALORANT - EMEA - Ascent)
INSERT INTO app.comps (game,region,league,patch,map,side,slots,pick_rate,win_rate,delta_win)
VALUES (
  'val','EMEA','VCT EMEA','9.15','Ascent','attack',
  '{
     "roles": ["smokes","initiator","duelist","sentinel","flex"],
     "members": [
       {"agent":"omen"},{"agent":"sova"},{"agent":"jett"},
       {"agent":"killjoy"},{"agent":"skye"}
     ]
   }'::jsonb,
  24.300, 52.100, 1.600
)
ON CONFLICT DO NOTHING;

-- +goose Down
-- (sin down para seeds)