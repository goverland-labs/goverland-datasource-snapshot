create table spaces
(
    id         text not null
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    snapshot   jsonb
);

create index idx_spaces_deleted_at
    on spaces (deleted_at);

create table proposals
(
    id             text not null
        primary key,
    space_id       text,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    snapshot       jsonb,
    vote_processed boolean
);

create index idx_proposals_deleted_at
    on proposals (deleted_at);

create index idx_proposals_space_id
    on proposals (space_id);

create table votes
(
    id             text not null
        primary key,
    ipfs           text,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    voter          text,
    space_id       text,
    proposal_id    text,
    choice         text,
    reason         text,
    app            text,
    vp             numeric,
    vp_by_strategy text,
    vp_state       text,
    published      boolean
);

create index idx_votes_deleted_at
    on votes (deleted_at);

create index votes_published_idx
    on votes (published, deleted_at, created_at);

