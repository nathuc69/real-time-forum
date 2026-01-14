--🌱🌱🌱
/*Le SEED (ensemencement): peupler la base avec des données initiales.
Consiste à insérer des données de départ dans la base, souvent utiles pour :
Tester l'application
Préparer les catégories, rôles, ou utilisateurs par défaut
Simuler des contenus (topics, messages…)*/
--🌱🌱🌱

-- ============= CATEGORIES =============
INSERT OR IGNORE INTO categories (name) VALUES 
	("Général"),
	("Astuces"),
	("Musculation"),
	("Sciences"),
	("Tech");

-- ============= USERS =============
INSERT OR IGNORE INTO users (username, password, email, age, gender, firstName, lastName, is_admin) VALUES 
	("nathox96", "$2a$10$weR/WktH5MhF8EU.6XvBFOG1uZHS4/64oL7wof3vrtEmU5t7SlTMK", "nathan240304@gmail.com", 20, "M", "Nathan", "Allain", 1), /* mot de passe : test */
	("dev_master", "$2a$10$weR/WktH5MhF8EU.6XvBFOG1uZHS4/64oL7wof3vrtEmU5t7SlTMK", "admin@test.com", 25, "M", "Admin", "Master", 1),
	("gym_pro", "$2a$10$weR/WktH5MhF8EU.6XvBFOG1uZHS4/64oL7wof3vrtEmU5t7SlTMK", "gym@example.com", 28, "M", "Jean", "Dupont", 0),
	("tech_girl", "$2a$10$weR/WktH5MhF8EU.6XvBFOG1uZHS4/64oL7wof3vrtEmU5t7SlTMK", "tech@example.com", 24, "F", "Marie", "Martin", 0),
	("science_nerd", "$2a$10$weR/WktH5MhF8EU.6XvBFOG1uZHS4/64oL7wof3vrtEmU5t7SlTMK", "science@example.com", 30, "M", "Pierre", "Leclerc", 0),
	("casual_user", "$2a$10$weR/WktH5MhF8EU.6XvBFOG1uZHS4/64oL7wof3vrtEmU5t7SlTMK", "casual@example.com", 22, "F", "Sophie", "Bernard", 0);

-- ============= POSTS =============
INSERT OR IGNORE INTO posts (title, content, user_id) VALUES 
	("Les bases de la musculation", "Salut à tous ! Je souhaite partager avec vous les fondamentaux de la musculation. D'abord, il faut comprendre qu'une bonne nutrition est essentielle. Ensuite, il faut s'entraîner régulièrement avec des poids progressifs.", 3),
	("Comment débuter en programmation ?", "Bonjour, je commence juste en programmation. Quelqu'un peut me recommander un langage pour débuter ? JavaScript me semble intéressant mais je ne sais pas par où commencer.", 4),
	("Les dernières découvertes en physique quantique", "Je vous propose de discuter des dernières avancées en physique quantique. L'intrication quantique est devenue de plus en plus accessible grâce aux nouvelles technologies.", 5),
	("Tips pour une meilleure récupération", "Après l'entraînement, la récupération est cruciale. Voici mes conseils : bien dormir, manger des protéines, faire des étirements et du stretching régulièrement.", 3),
	("React vs Vue.js en 2024", "Quelle est votre préférence ? React domine le marché mais Vue.js a gagné en popularité. Discutons des avantages et inconvénients de chacun.", 4),
	("Le rôle des mitochondries dans la cellule", "Les mitochondries sont les usines énergétiques de nos cellules. Elles produisent l'ATP qui alimente nos activités quotidiennes. Fascinant, non ?", 5),
	("Bienvenue sur le forum !", "Bonjour à tous et bienvenue ! Nous sommes heureux de vous avoir ici. N'hésitez pas à poser vos questions et à partager vos expériences.", 1),
	("Top 5 des erreurs en musculation", "Voici les 5 erreurs les plus courantes que je vois à la salle : 1) Trop de volume, 2) Mauvaise form, 3) Pas assez de repos, 4) Nutrition mauvaise, 5) Pas de progression.", 3);

-- ============= COMMENTS =============
INSERT OR IGNORE INTO comments (content, post_id, user_id) VALUES 
	("Excellent article ! Je vais suivre ces conseils.", 1, 6),
	("Merci pour le partage, très utile !", 1, 4),
	("Je recommande Python pour débuter, c'est plus facile que JavaScript.", 2, 5),
	("Faux, JavaScript est parfait pour débuter avec le web !", 2, 4),
	("La physique quantique me fascine aussi ! Avez-vous lu le dernier article de Science ?", 3, 5),
	("Super article sur la récupération ! Je vais appliquer ça.", 4, 6),
	("React, c'est une évidence pour moi. La communauté est immense.", 5, 4),
	("Mais Vue est tellement plus simple à apprendre !", 5, 6),
	("Les mitochondries c'est ouf ! Merci pour cette explication claire.", 6, 3);

-- ============= POSTS_CATEGORIES =============
INSERT OR IGNORE INTO posts_categories (post_id, category_id) VALUES
	(1, 3), -- Les bases de la musculation -> Musculation
	(2, 5), -- Comment débuter en programmation -> Tech
	(3, 4), -- Découvertes en physique -> Sciences
	(4, 3), -- Tips récupération -> Musculation
	(5, 5), -- React vs Vue -> Tech
	(6, 4), -- Mitochondries -> Sciences
	(7, 1), -- Bienvenue -> Général
	(8, 3); -- Top 5 erreurs musculation -> Musculation

-- ============= REACTIONS =============
INSERT OR IGNORE INTO reactions (value, target_type, target_id, user_id) VALUES
	(1, 'posts', 1, 4),
	(1, 'posts', 1, 5),
	(1, 'posts', 1, 6),
	(1, 'posts', 2, 3),
	(1, 'posts', 2, 5),
	(-1, 'posts', 2, 3),
	(1, 'posts', 3, 4),
	(1, 'posts', 5, 3),
	(1, 'posts', 8, 4),
	(1, 'posts', 8, 5);
