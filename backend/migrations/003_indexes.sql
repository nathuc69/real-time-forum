/*Un index est une structure qui permet à la base de retrouver 
plus rapidement les lignes correspondant à une requête. 
C’est comme un sommaire dans un livre
-Accélère les requêtes SELECT ... WHERE
-Optimise les jointures
-Rend les recherches plus performantes sur les colonnes souvent filtrées
*/

-- Index pour accélérer la recherche utilisateur
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Index sur les sujets
CREATE INDEX idx_topics_user ON topics(user_id);

-- Index sur les messages
CREATE INDEX idx_posts_topic ON posts(topic_id);
CREATE INDEX idx_posts_user ON posts(user_id);

-- Index sur les réactions
CREATE INDEX idx_reactions_target ON reactions(target_type, target_id);
CREATE INDEX idx_reactions_user ON reactions(user_id);

-- pour trier et filtrer par dates
CREATE INDEX idx_topics_created ON topics(created_at);
CREATE INDEX idx_messages_created ON posts(created_at);