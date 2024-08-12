CREATE TABLE game_user (
    id uuid,
    name VARCHAR(255) NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE game_state (
    id uuid,
    user_id uuid NOT NULL,
    games_played integer NOT NULL DEFAULT 0,
    score integer NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES game_user(id) ON DELETE CASCADE
);

CREATE TABLE game_friends (
    user_id uuid,
    friend_id uuid,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES game_user(id) ON DELETE CASCADE,
    CONSTRAINT fk_friend_id FOREIGN KEY (friend_id) REFERENCES game_user(id) ON DELETE CASCADE
);
