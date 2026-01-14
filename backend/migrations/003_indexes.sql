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

-- Index sur les posts
CREATE INDEX idx_posts_user ON posts(user_id);
CREATE INDEX idx_posts_created ON posts(created_at);

-- Index sur les commentaires
CREATE INDEX idx_comments_post ON comments(post_id);
CREATE INDEX idx_comments_user ON comments(user_id);
CREATE INDEX idx_comments_created ON comments(created_at);

-- Index sur les réactions
CREATE INDEX idx_reactions_target ON reactions(target_type, target_id);
CREATE INDEX idx_reactions_user ON reactions(user_id);