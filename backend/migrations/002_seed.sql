--🌱🌱🌱
/*Le SEED (ensemencement): peupler la base avec des données initiales.
Consiste à insérer des données de départ dans la base, souvent utiles pour :
Tester l'application
Préparer les catégories, rôles, ou utilisateurs par défaut
Simuler des contenus (topics, messages…)*/
--🌱🌱🌱

INSERT INTO categories (name) VALUES 
	("ERROR"), 
	("Correctifs"), 
	("Ça marche, mais je ne sais pas pourquoi"), 
	("Comment j'ai corrigé ce bug ?");

INSERT OR IGNORE INTO users (username, password, email) VALUES 
	("TheCatdu76", "$2a$10$5p7hI9IcwriOuFoqy8ZRPeVt2/UHU5aeSPBGTn233xy1Pzsu30Ica", "miaou@croquette.com"),
	("Coincoin", "$2a$10$5p7hI9IcwriOuFoqy8ZRPeVt2/UHU5aeSPBGTn233xy1Pzsu30Ica", "agrou@piouipou.cuicui");

INSERT OR IGNORE INTO topics (title, content, user_id) VALUES 
	("Docker c'est trop bien", "J'adore !", 1),
	("Le Front-End, c'est la vie", "En vrai, face au Back... ", 2);

INSERT OR IGNORE INTO posts (content, topic_id, user_id) VALUES 
	("Docker For Ever", 1, 1);

INSERT OR IGNORE INTO topic_categories (topic_id, category_id) VALUES
	(1,3),
	(2,4);

INSERT OR IGNORE INTO users (username, password, email) VALUES 
	("nathox96", "$2a$10$weR/WktH5MhF8EU.6XvBFOG1uZHS4/64oL7wof3vrtEmU5t7SlTMK", "nathan240304@gmail.com");
