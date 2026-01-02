export default {
    fr: {
        lang: {
            fr: '🇫🇷',
            en: '🇬🇧',
            undefined: '🏁', // Easter-egg from development
        },
        a11y: {
            logo: "Logo de TiL"
        },
        menu: {
            home: "Accueil",
            bookmarks: "Favoris",
            extension: "Extension",
            logout: "Déconnexion",
            publish: "Publier",
            profile: "Profil"
        },
        item: {
            by: "Par {{user}}, le {{date}}"
        },
        post: {
            empty_list: "Aucun contenu à afficher... Changez vos filtres ?",
            page_title: "Ajouter un contenu",
            created: "Le contenu a bien été publié !",
            title: "Nom du contenu :",
            title_placeholder: "Ma super news",
            link: "Lien :",
            link_placeholder: "https://google.com/",
            lang: "Langue :",
            tags: "Tags :",
            comment: "Commentaire additionnel :",
            create_button: "Ajouter le contenu",
            deleted: "Article supprimé"
        },
        bookmarks: {
            added: "Article ajouté aux favoris !",
            deleted: "Article supprimé des favoris.",
            empty_list: "Vous n'avez aucun article favori"
        },
        index: {
            filter: {
                label: "Filtrer :",
                remove_all: "Supprimer filtres",
                restore: "Réinitialiser filtres"
            }
        },
        profile: {
            updated: "Les informations ont bien été mises à jour.",
            renewed_feed_key: "La clé de flux a bien été régénérée.",
            automatic_tags: {
                title: "Mes tags automatiques",
                description: "Définir des tags vous permet de personnaliser votre flux RSS et d'afficher directement les articles correspondant à vos préférences sur la page principale.",
                save_button: "Sauvegarder"
            },
            rss_feed: {
                title: "Mon flux RSS",
                description: "Vous pouvez directement récupérer votre flux d'actualité au format RSS (Atom 2.0). Pour cela, copiez/collez le lien suivant dans votre aggrégateur favori :",
                tags_manual: "Vous pouvez spécifier des tags différents de ceux liés à votre profil avec l'argument",
                warning: {
                    content: "Ne partagez jamais votre clé de flux RSS ! Un utilisateur mal-intentionné pourrait s'en servir pour accèder à votre flux de veille technologique. Si votre clé est compromise, veuillez cliquer sur le bouton ci-dessous pour régénérer une clé.",
                    disconnection_warning: "Attention, cette opération déconnectera automatiquement tout vos aggrégateurs et extensions de navigateurs liées !",
                    renew_button: "Réinitialiser ma clé de flux RSS"
                }
            }
        },
        generic_error: "Une erreur est survenue : \"{{err}}\"",
    },
    en: {
        lang: {
            fr: '🇫🇷',
            en: '🇬🇧',
            undefined: '🏁', // Easter-egg from development
        },
        a11y: {
            logo: "TiL's logo"
        },
        menu: {
            home: "Home",
            bookmarks: "Bookmarks",
            extension: "Extension",
            logout: "Logout",
            publish: "Publish",
            profile: "Profile"
        },
        item: {
            by: "By {{user}}, on {{date}}"
        },
        post: {
            empty_list: "There's no content... Maybe change your filters?",
            page_title: "Add content",
            created: "Content is now published!",
            title: "Name:",
            title_placeholder: "My epic news",
            link: "Link:",
            link_placeholder: "https://google.com/",
            lang: "Language:",
            tags: "Tags:",
            comment: "Additional comment:",
            create_button: "Publish content",
            deleted: "Content deleted"
        },
        bookmarks: {
            added: "Post added to bookmarks!",
            deleted: "Post removed from bookmarks.",
            empty_list: "You don't have any bookmarks"
        },
        index: {
            filter: {
                label: "Filter:",
                remove_all: "Remove filters",
                restore: "Restore user filters"
            }
        },
        profile: {
            updated: "Informations updated.",
            renewed_feed_key: "RSS Feed Key regenerated.",
            automatic_tags: {
                title: "My automatic tags",
                description: "The following tags will be automatically loaded on the index page, and will personalize your RSS feed.",
                save_button: "Save"
            },
            rss_feed: {
                title: "My RSS Feed",
                description: "You can retrieve your feed directly from RSS (Atom 2.0). Simply copy/paste this link to your favorite aggregator:",
                tags_manual: "You can specify custom tags with the following argument next to your feed key:",
                warning: {
                    content: "Never share your RSS feed key ! A malicious third-party can use it to access your data. If you think your key is compromised, please click the following button to generate a new one.",
                    disconnection_warning: "Please note: this action will disconnect any aggregator and browser extension linked to this account!",
                    renew_button: "Renew my RSS Feed Key"
                }
            }
        },
        generic_error: "An error occurred : \"{{err}}\"",
    }
};